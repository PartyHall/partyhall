package nexus

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

var INSTANCE NexusSync

var ErrUnauthorized = errors.New("the credentials for PartyNexus are invalid (HTTP 401)")
var ErrForbidden = errors.New("you cant access the requested resource (HTTP 403)")

// @TODO: Handle 4xx (especially 401), 5xx and let the user know properly in the logs !

type NexusSync struct {
	BaseURL    string
	HardwareID string
	ApiKey     string

	http      *http.Client
	ignoreSsl bool

	IsSetup bool
}

func (ns NexusSync) setUserAgent(req *http.Request) {
	req.Header.Set(
		"User-Agent",
		fmt.Sprintf("PartyHall Appliance (%s [%s])", utils.CURRENT_VERSION, utils.CURRENT_COMMIT),
	)
}

func (ns *NexusSync) doJsonRequest(httpMethod string, endpoint string, data map[string]any) (map[string]any, error) {
	var reader io.Reader = nil

	if data != nil {
		input, _ := json.Marshal(data)
		reader = bytes.NewReader(input)
	}

	req, err := http.NewRequest(httpMethod, ns.BaseURL+endpoint, reader)
	if err != nil {
		return nil, err
	}

	ns.setUserAgent(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-HARDWARE-ID", ns.HardwareID)
	req.Header.Set("X-API-TOKEN", ns.ApiKey)

	res, err := ns.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		return nil, ErrUnauthorized
	}

	if res.StatusCode == 403 {
		return nil, ErrForbidden
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var respData map[string]any
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (ns *NexusSync) fetchSongPage(songPage int) ([]ApiSong, int, error) {
	data, err := ns.doJsonRequest(http.MethodGet, fmt.Sprintf("/api/songs?page=%v", songPage), nil)
	if err != nil {
		return nil, -1, err
	}

	songs := []ApiSong{}

	itemsRaw, ok := data["member"]
	if !ok {
		return nil, -1, errors.New("bad request from the api: no member in the response")
	}

	itemsSlice, ok := itemsRaw.([]interface{})
	if !ok {
		return nil, -1, errors.New("bad request from the api: member is not an array")
	}

	for _, item := range itemsSlice {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to marshal song data: %w", err)
		}

		var song ApiSong
		if err := json.Unmarshal(itemJSON, &song); err != nil {
			return nil, -1, fmt.Errorf("failed to unmarshal song: %w", err)
		}

		songs = append(songs, song)
	}

	totalItemsRaw, ok := data["totalItems"]
	if !ok {
		return nil, -1, errors.New("total items field is missing from response")
	}

	totalItemsFloat, ok := totalItemsRaw.(float64)
	if !ok {
		return nil, -1, errors.New("total items is not a number")
	}

	totalItems := int(totalItemsFloat)

	return songs, totalItems, nil
}

func (ns *NexusSync) fetchAllSongs() ([]ApiSong, error) {
	total := -1

	songs := []ApiSong{}

	page := 1
	for total < 0 || len(songs) < total {
		pageSongs, amtSongs, err := ns.fetchSongPage(page)
		if err != nil {
			return nil, err
		}

		total = int(amtSongs)
		songs = append(songs, pageSongs...)
		page++
	}

	return songs, nil
}

func (ns *NexusSync) CreateEvent(eventId int64) error {
	event, err := dal.EVENTS.Get(eventId)
	if err != nil {
		return err
	}

	if event.NexusId.Valid {
		return errors.New("this event is already created on PartyNexus")
	}

	resp, err := ns.doJsonRequest(http.MethodPost, "/api/events", map[string]any{
		"name":     event.Name,
		"author":   event.Author,
		"datetime": event.Date.Format(time.RFC3339),
		"location": event.Location,
	})
	if err != nil {
		return err
	}

	id, ok := resp["id"].(string)
	if !ok {
		return errors.New("failed to create event: failed to parse id from PartyNexus")
	}

	state.STATE.CurrentEvent.NexusId = models.JsonnableNullstring{String: id, Valid: true}

	return dal.EVENTS.Update(state.STATE.CurrentEvent)
}

func (ns *NexusSync) downloadSong(id int64) (string, error) {
	tmpFile, err := os.CreateTemp("", "phksong-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%v/api/songs/%v/download", ns.BaseURL, id),
		nil,
	)
	if err != nil {
		return "", err
	}

	ns.setUserAgent(req)
	req.Header.Set("X-HARDWARE-ID", ns.HardwareID)
	req.Header.Set("X-API-TOKEN", ns.ApiKey)

	resp, err := ns.http.Do(req)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpFile.Name())
		return "", errors.New("bad status: " + resp.Status)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	absPath, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return absPath, nil
}

func (ns *NexusSync) syncPicture(event *models.Event, p models.Picture) error {
	log.Info("Syncing picture", "picture_id", p.Id)

	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)

	part, err := mwriter.CreateFormFile("file", p.Filename)
	if err != nil {
		return err
	}

	filename := p.Filename
	if p.Unattended {
		filename = filepath.Join("unattended", filename)
	}

	imagePath := filepath.Join(
		config.GET.EventPath,
		fmt.Sprintf("%v", event.Id),
		"photobooth",
		filename,
	)

	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	file.Close()

	mwriter.WriteField("event", fmt.Sprintf("/api/events/%s", event.NexusId.String))
	mwriter.WriteField("takenAt", p.TakenAt.Format(time.RFC3339))
	mwriter.WriteField("applianceUuid", p.Id)

	if p.Unattended {
		mwriter.WriteField("unattended", "true")
	} else {
		mwriter.WriteField("unattended", "false")
	}

	mwriter.Close()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/api/pictures", // events/%s/
			ns.BaseURL,
			// event.NexusId.String,
		),
		body,
	)
	if err != nil {
		return err
	}

	ns.setUserAgent(req)
	req.Header.Set("Content-Type", mwriter.FormDataContentType())
	req.Header.Set("X-HARDWARE-ID", ns.HardwareID)
	req.Header.Set("X-API-TOKEN", ns.ApiKey)

	resp, err := ns.http.Do(req)
	if err != nil {
		return err
	}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New("failed to upload picture: " + string(bodyData))
	}

	var data map[string]any
	err = json.Unmarshal(bodyData, &data)
	if err != nil {
		return err
	}

	pictureId, ok := data["id"]
	if !ok {
		return errors.New("no id for the uploaded picture")
	}

	nexusId, ok := pictureId.(string)
	if !ok {
		return errors.New("the picture id should be a string")
	}

	p.NexusId = models.JsonnableNullstring{String: nexusId, Valid: true}

	return dal.EVENTS.UpdatePicture(p)
}

func (ns *NexusSync) syncPictures(event *models.Event, shouldSendStatus bool) error {
	if !ns.IsSetup {
		log.Info("No PartyNexus credentials. No sync will be made")

		return nil
	}

	if state.STATE.CurrentEvent == nil {
		log.Info("No current event. No sync will be made ")
		return nil
	}

	if !state.STATE.CurrentEvent.NexusId.Valid {
		log.Info("Current event not created on Nexus. Skipping sync pictures")
		return nil
	}

	log.Info("Fetching pictures to upload")
	if shouldSendStatus {
		state.STATE.SyncInProgress = true
		mercure_client.CLIENT.PublishSyncInProgress()
	}

	pictures, err := dal.EVENTS.GetUnsubmittedPictures(state.STATE.CurrentEvent.Id)
	if err != nil {
		state.STATE.SyncInProgress = false
		mercure_client.CLIENT.PublishSyncInProgress()

		return err
	}

	// We are bulking the errors so that
	// if one fail, the other will still upload
	// and the karaoke songs will still sync
	pictureUploadErrors := []error{}

	for _, picture := range pictures {
		if picture.NexusId.Valid {
			continue
		}

		err := ns.syncPicture(state.STATE.CurrentEvent, picture)
		if err != nil {
			pictureUploadErrors = append(pictureUploadErrors, err)
		}
	}

	if shouldSendStatus {
		state.STATE.SyncInProgress = false
		mercure_client.CLIENT.PublishSyncInProgress()
	}

	return errors.Join(pictureUploadErrors...)
}

func (ns *NexusSync) syncSongs() error {
	log.Info("Fetching songs from Nexus")

	// Fetch all songs from the API
	nexusSongs, err := ns.fetchAllSongs()
	if err != nil {
		return err
	}

	// Fetch all song in the DB
	applianceSongs, err := dal.SONGS.GetAll()
	if err != nil {
		return err
	}

	// Removing those that are no longer on Nexus
	log.Info("Removing song no longer available")
	for _, aSong := range applianceSongs {
		found := false

		for _, nSong := range nexusSongs {
			if nSong.NexusBuildId != nil && aSong.NexusId == *nSong.NexusBuildId {
				found = true
				break
			}
		}

		if !found {
			log.Info("Removing song", "title", aSong.Title, "artist", aSong.Artist)
			err = dal.SONGS.Delete(aSong)
			if err != nil {
				return err
			}
		}
	}

	phkToDownload := []int64{}

	// Downloading those which are not available locally
	log.Info("Downloading new songs")
	for _, nSong := range nexusSongs {
		found := false

		for _, aSong := range applianceSongs {
			if nSong.NexusBuildId != nil && aSong.NexusId == *nSong.NexusBuildId {
				found = true
				break
			}
		}

		if !found {
			phkToDownload = append(phkToDownload, nSong.Id)
		}
	}

	for _, id := range phkToDownload {
		log.Info("Downloading song", "id", id)
		songPhkPath, err := ns.downloadSong(id)
		if err != nil {
			return err
		}

		log.Info("Importing song", "id", id)
		if err := importSong(songPhkPath); err != nil {
			return err
		}
	}

	log.Info("Song sync done")

	return nil
}

func (ns *NexusSync) syncSession(session *models.SongSession) error {
	if session.SessionNexusId.Valid {
		log.Warn("Trying to sync an already synced song session", "sessionId", session.Id)

		return nil
	}

	resp, err := ns.doJsonRequest(
		http.MethodPost,
		"/api/song_sessions",
		map[string]any{
			"title":  session.Title,
			"artist": session.Artist,
			"sungAt": session.StartedAt.Time.Format(time.RFC3339),
			"singer": session.SungBy,
			"event":  "/api/events/" + state.STATE.CurrentEvent.NexusId.String,
		},
	)

	if err != nil {
		return err
	}

	nexusId, ok := resp["id"]
	if !ok {
		return errors.New("no id retreived from creating the session on PartyNexus")
	}

	nexusIdInt, ok := nexusId.(float64)
	if !ok {
		return errors.New("id retreived from creating the session on PartyNexus is not a float64")
	}

	session.SessionNexusId = models.JsonnableNullInt64{
		Int64: int64(nexusIdInt),
		Valid: true,
	}

	return dal.SONGS.UpdateSession(session)
}

func (ns *NexusSync) syncSessions() error {
	if state.STATE.CurrentEvent == nil || !state.STATE.CurrentEvent.NexusId.Valid {
		log.Info("No current event or no partynexus id. No sync session will be made ")

		return nil
	}

	notSyncedSessions, err := dal.SONGS.GetNotSyncedSessions(state.STATE.CurrentEvent.Id)
	if err != nil {
		return err
	}

	for _, session := range notSyncedSessions {
		err := ns.syncSession(&session)
		if err != nil {
			return err
		}
	}

	log.Info("Song session sync done")

	return nil
}

func (ns *NexusSync) Sync(event *models.Event) error {
	if !ns.IsSetup {
		log.Info("No PartyNexus credentials. No sync will be made")

		return nil
	}

	if state.STATE.SyncInProgress {
		return errors.New("trying to sync while already syncronizing")
	}

	state.STATE.SyncInProgress = true
	mercure_client.CLIENT.PublishSyncInProgress()

	err := ns.syncSongs()
	if err != nil {
		mercure_client.CLIENT.PublishSyncInProgress()
		log.Error("Failed to sync songs", "err", err)
	}

	if state.STATE.CurrentEvent == nil {
		state.STATE.SyncInProgress = false
		mercure_client.CLIENT.PublishSyncInProgress()

		log.Info("No current event. No sync will be made ")

		return nil
	}

	err = ns.syncPictures(state.STATE.CurrentEvent, false)
	if err != nil {
		mercure_client.CLIENT.PublishSyncInProgress()
		log.Error("Failed to sync songs", "err", err)
	}

	err = ns.syncSessions()
	if err != nil {
		mercure_client.CLIENT.PublishSyncInProgress()
		log.Error("Failed to sync songs", "err", err)
	}

	state.STATE.SyncInProgress = false
	mercure_client.CLIENT.PublishSyncInProgress()

	return nil
}

func NewClient(baseUrl string, hwid string, apiKey string, ignoreSsl bool) {
	log.Info("Initializing PartyNexus client")
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	isSetup := len(baseUrl) > 0 && len(hwid) > 0 && len(apiKey) > 0
	if !isSetup {
		log.Warn("No credentials were setup for PartyNexus. No sync will be made.")
	}

	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if ignoreSsl {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	httpClient := &http.Client{
		Transport: &transport,
		Timeout:   30 * time.Second,
	}

	INSTANCE = NexusSync{
		BaseURL:    baseUrl,
		HardwareID: hwid,
		ApiKey:     apiKey,

		http:      httpClient,
		ignoreSsl: ignoreSsl,
		IsSetup:   isSetup,
	}
}
