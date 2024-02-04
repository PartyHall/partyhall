package module_karaoke

import (
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
	"gopkg.in/yaml.v2"
)

var (
	INSTANCE = &ModuleKaraoke{
		Actions: Actions{},
	}

	CONFIG = Config{}

	ALLOWED_EXT_LIST = []string{
		"mp4",
		"webm",
	}
)

/**
 * State explanation:
 * Started = false: waiting, nothing happening, maybe renamed to pause ?
 * Preplay timer: Delay before the start of the song to display stuff on the screen
 */
type ModuleKaraoke struct {
	Actions Actions

	CurrentSong  *models.Song
	Queue        []*models.Song
	Started      bool
	PreplayTimer int

	// This means the queue is running
	// E.g. you click on play a song and want to play it alone, it won't play the rest of the queue
	// E.g. you want to add multiple song to the queue before starting, it won't start right away
	IsQueuePlaying bool
}

func (m ModuleKaraoke) GetModuleName() string {
	return "Karaoke"
}

func (m ModuleKaraoke) LoadConfig(filename string) error {
	if !utils.FileExists(filename) {
		CONFIG = Config{
			AmtSongsPerPage: 20,
			PrePlayTimer:    5,
		}

		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	CONFIG = cfg

	return nil
}

func (m ModuleKaraoke) PreInitialize() error {
	return nil
}

func (m ModuleKaraoke) Initialize() error {
	go func() {
		for {
			m.ScanSongs()
			time.Sleep(5 * time.Minute)
		}
	}()

	return nil
}

// #region ScanSongs
// ScanSongs checks the song directory to add new cdg files to the database
// It also removes from the DB songs that do no longer have cdg+mp3 files
func (m ModuleKaraoke) ScanSongs() error {
	logs.Info("[Karaoke] Scanning the songs")
	baseDir := filepath.Join(config.GET.RootPath, "karaoke")
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		logs.Error("Failed to read songs directory!", err)
		return err
	}

	dbSongs, err := ormFetchSongFilenames()
	if err != nil {
		logs.Error("Failed to fetch songs from DB: ", err)
		return err
	}
	folderSongs := []string{}

	if err != nil {
		logs.Error("Failed to fetch songs in DB, skipping the scanning process")
		return err
	}

	//#region Purge invalid songs
	// @TODO: In case of video
	// Check whether the DB song has the correct extension to consider it valid
	for _, entry := range entries {
		songName := entry.Name()
		basePath := filepath.Join(baseDir, songName)

		if fi, err := os.Stat(basePath); err == nil && !fi.IsDir() {
			continue
		}

		isSongValid := false

		hasCdg := utils.FileExists(filepath.Join(basePath, "song.cdg"))
		hasMp3 := utils.FileExists(filepath.Join(basePath, "song.mp3"))
		foundExt := utils.FileExistsForAnyExt(filepath.Join(basePath, "song"), ALLOWED_EXT_LIST)

		if len(foundExt) > 0 || (hasMp3 && hasCdg) {
			isSongValid = true
		}

		if _, err := os.Stat(filepath.Join(basePath, "cover.jpg")); isSongValid && os.IsNotExist(err) {
			logs.Warn("No cover found for song " + songName)
		}

		if !isSongValid {
			logs.Error("Invalid song detected: " + songName + ". Removing it.")

			// Removing the song
			os.RemoveAll(basePath)
			ormDeleteSong(songName)
		} else {
			folderSongs = append(folderSongs, songName)
		}
	}
	//#endregion

	//#region Removing songs that are in database but not in folders
	for _, songName := range dbSongs {
		if slices.Contains(folderSongs, songName) {
			continue
		}

		logs.Error("Song " + songName + " no longer available, removing it from the database.")
		ormDeleteSong(songName)
	}
	//#endregion

	//#region Finally, adding songs that are not in the DB yet
	for _, songName := range folderSongs {
		if slices.Contains(dbSongs, songName) {
			continue
		}

		artist := ""
		title := ""
		format := ""

		infoFile := filepath.Join(baseDir, songName, "info.txt")
		data, err := os.ReadFile(infoFile)
		if err != nil {
			logs.Warn("Song " + songName + " doesn't have info.txt. Fill the info manually!")
		} else {
			dataStr := strings.Split(string(data), "\n")
			if len(dataStr) < 3 {
				logs.Warn("Song " + songName + " have invalid info.txt. Fill the info manually!")
			} else {
				artist = dataStr[0]
				title = dataStr[1]
				format = dataStr[2]
			}
		}

		if len(format) == 0 {
			// Wtf ?
			foundCdg := false
			foundVideo := false

			if !foundCdg && !foundVideo {
				logs.Error("Failed to create song " + songName + ", the folder should contain either a cdg or a video")
				continue
			}
		}

		_, err = ormCreateSong(songName, artist, title, format)
		if err != nil {
			logs.Error("Failed to create song "+songName+": ", err)
		} else {
			logs.Info("Song " + songName + " created.")
		}
	}
	//#endregion

	return nil
}

//#endregion

func (m ModuleKaraoke) GetMqttHandlers() map[string]mqtt.MessageHandler {
	return map[string]mqtt.MessageHandler{}
}

func (m ModuleKaraoke) GetWebsocketHandlers() []easyws.MessageHandler {
	return []easyws.MessageHandler{
		PlaySongHandler{},
		PlayingStatusHandler{},
		PlayingEndedHandler{},
		AddToQueueHandler{},
		DelFromQueueHandler{},
		PauseHandler{},
		QueueMoveUp{},
		QueueMoveDown{},
	}
}

func (m ModuleKaraoke) UpdateFrontendSettings() {
	queue := m.Queue
	if queue == nil {
		queue = []*models.Song{}
	}

	services.GET.ModuleSettings["karaoke"] = map[string]interface{}{
		"currentSong":  m.CurrentSong,
		"queue":        queue,
		"started":      m.Started,
		"preplayTimer": m.PreplayTimer,
	}
}

func (m ModuleKaraoke) RegisterApiRoutes(g *echo.Group) {
	g.GET("/song", searchSongs, services.GET.EchoJwtMiddleware)
	g.POST("/song", songPost, services.GET.EchoJwtMiddleware, middlewares.RequireAdmin)
	g.POST("/rescan", rescanSongs, services.GET.EchoJwtMiddleware, middlewares.RequireAdmin)
	g.POST("/spotify-search", spotifySearch, services.GET.EchoJwtMiddleware, middlewares.RequireAdmin)

	g.GET("/fallback-image", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "image/jpeg")
		return c.Blob(http.StatusOK, "", services.KARAOKE_FALLBACK_IMAGE)
	})

	g.Static("/medias", utils.GetPath("karaoke"))
}
