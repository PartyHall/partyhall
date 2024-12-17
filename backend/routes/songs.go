package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/nexus"
	routes_requests "github.com/partyhall/partyhall/routes/requests"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/state"
)

func routeGetSongs(c *gin.Context) {
	offset := 0

	search := c.Query("search")
	page := c.Query("page")

	var pageInt int = 1
	var err error
	if len(page) > 0 {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
				"page": "The page should be an integer",
			}))

			return
		}

		offset = (pageInt - 1) * config.AMT_RESULTS_PER_PAGE
	}

	songs, err := dal.SONGS.GetCollection(search, config.AMT_RESULTS_PER_PAGE, offset)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	songs.Page = pageInt

	c.JSON(http.StatusOK, songs)
}

func routeGetSong(c *gin.Context) {
	id := strings.TrimSpace(c.Params.ByName("songId"))
	if len(id) == 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Song ID should not be empty",
		}))

		return
	}

	evt, err := dal.SONGS.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.JSON(http.StatusOK, evt)
}

func routeServeSong(c *gin.Context) {
	songID := c.Param("songId")
	path := c.Param("filepath")

	// Basic traversal protection
	// Probably can be bypassed but meh
	// Appliances are "single use" anyway
	if strings.Contains(path, "..") {
		c.String(http.StatusForbidden, "Invalid path")
		return
	}

	songPath := filepath.Join(config.GET.RootPath, "karaoke", songID)
	if _, err := os.Stat(songPath); os.IsNotExist(err) {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Song does not exists",
		}))

		return
	}

	fs := http.FileServer(http.Dir(songPath))
	c.Request.URL.Path = path
	fs.ServeHTTP(c.Writer, c.Request)
}

// @TODO: Replace Status 409 with a custom error telling "the song is already in the queue"
func routeQueueAdd(c *gin.Context) {
	var req routes_requests.SongEnqueue
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	req.SongId = strings.ToLower(req.SongId)

	currentSong := state.STATE.Karaoke.Current
	if currentSong != nil && strings.ToLower(currentSong.NexusId) == req.SongId {
		c.Status(409)

		return
	}

	for _, session := range state.STATE.KaraokeQueue {
		if strings.ToLower(session.NexusId) == req.SongId {
			c.Status(409)

			return
		}
	}

	song, err := dal.SONGS.Get(req.SongId)
	if err != nil {
		log.Error("Failed to find song to add in queue", "err", err)

		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Song does not exists",
		}))

		return
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)

	log.Info(
		"Song requested",
		"username", claims.Username,
		"displayName", req.DisplayName,
		"direct", req.DirectPlay,
		"song", song.NexusId,
		"title", fmt.Sprintf("%v by %v", song.Title, song.Artist),
	)

	session := models.SongSession{
		EventId:  state.STATE.CurrentEvent.Id,
		Title:    song.Title,
		Artist:   song.Artist,
		SungBy:   req.DisplayName,
		SungById: fmt.Sprintf("%v", claims.Subject),
		Song:     song,
		AddedAt: models.JsonnableNullTime{
			Time:  time.Now(),
			Valid: true,
		},
		StartedAt:   models.JsonnableNullTime{},
		EndedAt:     models.JsonnableNullTime{},
		CancelledAt: models.JsonnableNullTime{},
	}

	directPlay := state.STATE.Karaoke.Current == nil || req.DirectPlay

	err = dal.SONGS.CreateSession(&session)
	if err != nil {
		log.Error("Failed to add in queue", "err", err)

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	if directPlay {
		if err := services.KARAOKE.StartSong(&session); err != nil {
			log.Error("Failed to direct play the song: ", "err", err)
			c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
				"err": err,
			}))

			return
		}
	} else {
		state.STATE.KaraokeQueue = append(state.STATE.KaraokeQueue, &session)

		mercure_client.CLIENT.PublishEvent(
			"/karaoke-queue",
			state.STATE.KaraokeQueue,
		)
	}

	c.JSON(200, session)
}

func routeSessionDirectPlay(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")

	sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid",
		}))

		return
	}

	var idx = -1
	var session *models.SongSession = nil
	for i, s := range state.STATE.KaraokeQueue {
		if s.Id == sessionId {
			session = s
			idx = i
			break
		}
	}

	if session == nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid (not in queue)",
		}))

		return
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Directplay from queue",
		"username", claims.Username,
		"session", session.Id,
		"song", session.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", session.Song.Title, session.Song.Artist),
	)

	state.STATE.KaraokeQueue = append(state.STATE.KaraokeQueue[:idx], state.STATE.KaraokeQueue[idx+1:]...)

	mercure_client.CLIENT.PublishEvent(
		"/karaoke-queue",
		state.STATE.KaraokeQueue,
	)

	if err := services.KARAOKE.StartSong(session); err != nil {
		log.Error("Failed to direct play the song: ", "err", err)
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	c.JSON(200, session)
}

func routeQueueRemove(c *gin.Context) {
	sessionId := c.Param("sessionId")

	session, err := dal.SONGS.GetSession(sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(404)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	if session.EndedAt.Valid || session.CancelledAt.Valid {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session is already over",
		}))

		return
	}

	session.CancelledAt = models.JsonnableNullTime{Valid: true, Time: time.Now()}

	err = dal.SONGS.UpdateSession(session)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err.Error(),
		}))

		return
	}

	title := ""
	nexusId := ""
	if session != nil && session.Song != nil {
		title = fmt.Sprintf("%v by %v", session.Song.Title, session.Song.Artist)
		nexusId = session.Song.NexusId
	} else {
		title = fmt.Sprintf("%v by %v", session.Title, session.Artist)
		nexusId = "Song unloaded"
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Removing from queue",
		"username", claims.Username,
		"session", session.Id,
		"song", nexusId,
		"title", title,
	)

	if state.STATE.Karaoke.Current != nil && state.STATE.Karaoke.Current.Id == session.Id {
		services.KARAOKE.StartNextSong()
	} else {
		newSessions := []*models.SongSession{}
		for _, currSession := range state.STATE.KaraokeQueue {
			if currSession.Id != session.Id {
				newSessions = append(newSessions, currSession)
			}
		}

		state.STATE.KaraokeQueue = newSessions

		mercure_client.CLIENT.PublishEvent(
			"/karaoke-queue",
			state.STATE.KaraokeQueue,
		)
	}

	c.JSON(200, session)
}

func routeSetTimecode(c *gin.Context) {
	var req routes_requests.SetTimecode
	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	state.STATE.Karaoke.Timecode = req.Timecode

	mercure_client.CLIENT.PublishEvent(
		"/karaoke-timecode",
		routes_requests.SetTimecode{Timecode: state.STATE.Karaoke.Timecode},
	)

	c.Status(200)
}

func routeSetEnded(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")
	current := state.STATE.Karaoke.Current

	sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
	if err != nil || current == nil || current.Id != sessionId {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid",
		}))

		return
	}

	current.EndedAt = models.JsonnableNullTime{
		Time:  time.Now(),
		Valid: true,
	}

	log.Info(
		"Song ended",
		"song", current.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", current.Song.Title, current.Song.Artist),
	)

	err = dal.SONGS.UpdateSession(current)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{
			"err": err,
		}))

		return
	}

	err = nexus.INSTANCE.Sync(state.STATE.CurrentEvent)
	if err != nil {
		log.Error("Failed to sync", "err", err)
	}

	state.STATE.Karaoke.Current = nil

	services.KARAOKE.StartNextSong()
}

func routeSetPlayingStatus(c *gin.Context) {
	statusStr := c.Param("status")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Status should be true or false",
		}))

		return
	}

	state.STATE.Karaoke.IsPlaying = status
	tx := ""
	if status {
		tx = "resumed"
	} else {
		tx = "paused"
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Song "+tx,
		"username", claims.Username,
		"session", state.STATE.Karaoke.Current.Id,
		"song", state.STATE.Karaoke.Current.Song.NexusId,
		"title", fmt.Sprintf("%v by %v", state.STATE.Karaoke.Current.Song.Title, state.STATE.Karaoke.Current.Song.Artist),
	)

	mercure_client.CLIENT.PublishEvent(
		"/karaoke",
		state.STATE.Karaoke,
	)
}

func routeMoveInQueue(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")
	direction := strings.ToLower(c.Param("direction"))

	sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Session ID is invalid",
		}))

		return
	}

	if direction != "up" && direction != "down" {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The direction should be either up or down",
		}))

		return
	}

	sessionIdx := -1

	for idx, session := range state.STATE.KaraokeQueue {
		if session.Id == sessionId {
			sessionIdx = idx
			break
		}
	}

	if sessionIdx < 0 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "The session was not found in the queue",
		}))

		return
	}

	if len(state.STATE.KaraokeQueue) < 2 {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Not enough sessions in the queue to perform movement",
		}))

		return
	}

	newPos := sessionIdx - 1
	if direction == "up" {
		if sessionIdx == 0 {
			c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
				"err": "Cannot move the first session up",
			}))

			return
		}

	} else if direction == "down" {
		if sessionIdx == len(state.STATE.KaraokeQueue)-1 {
			c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
				"err": "Cannot move the last session down",
			}))

			return
		}

		newPos = sessionIdx + 1
	}

	ctxClaims, _ := c.Get("TokenClaims")
	claims := ctxClaims.(*models.JwtCustomClaims)
	log.Info(
		"Session moved in queue",
		"username", claims.Username,
		"session", state.STATE.KaraokeQueue[sessionIdx].Id,
		"song", state.STATE.KaraokeQueue[sessionIdx].Song.NexusId,
		"title", fmt.Sprintf("%v by %v", state.STATE.KaraokeQueue[sessionIdx].Song.Title, state.STATE.KaraokeQueue[sessionIdx].Song.Artist),
		"previousPos", sessionIdx,
		"newPos", newPos,
	)

	state.STATE.KaraokeQueue[sessionIdx], state.STATE.KaraokeQueue[newPos] = state.STATE.KaraokeQueue[newPos], state.STATE.KaraokeQueue[sessionIdx]

	mercure_client.CLIENT.PublishEvent(
		"/karaoke-queue",
		state.STATE.KaraokeQueue,
	)

	c.JSON(200, state.STATE.KaraokeQueue)
}

func routeSetVolume(c *gin.Context) {
	volumeType := strings.ToLower(c.Param("type"))

	if volumeType != "instrumental" && volumeType != "vocals" && volumeType != "combined" {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Volume type should be either instrumental, vocals or combined",
		}))

		return
	}

	volStr := c.Param("volume")
	vol, err := strconv.Atoi(volStr)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.BAD_REQUEST.WithExtra(map[string]any{
			"err": "Volume should be an int",
		}))

		return
	}

	if vol < 0 {
		vol = 0
	} else if vol > 100 {
		vol = 100
	}

	if volumeType == "instrumental" {
		state.STATE.Karaoke.Volume = vol
	} else if volumeType == "vocals" {
		state.STATE.Karaoke.VolumeVocals = vol
	}

	mercure_client.CLIENT.PublishEvent(
		"/karaoke",
		state.STATE.Karaoke,
	)

	c.JSON(200, state.STATE.Karaoke)
}
