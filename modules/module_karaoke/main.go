package module_karaoke

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
	"gopkg.in/yaml.v2"
)

var (
	INSTANCE = &ModuleKaraoke{
		Actions: Actions{},
	}

	CONFIG = Config{}
)

/**
 * State explanation:
 * Started = false: waiting, nothing happening, maybe renamed to pause ?
 * Preplay timer: Delay before the start of the song to display stuff on the screen
 */
type ModuleKaraoke struct {
	Actions Actions

	CurrentSong  *SongSession
	Queue        []SongSession
	Started      bool
	PreplayTimer int

	// This means the queue is running
	// E.g. you click on play a song and want to play it alone, it won't play the rest of the queue
	// E.g. you want to add multiple song to the queue before starting, it won't start right away
	IsQueuePlaying bool

	VolumeInstru float64
	VolumeVocals float64
	VolumeFull   float64
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
	remote.RegisterOnJoin("karaoke", func(socketType string, s *easyws.Socket) {
		if socketType == utils.SOCKET_TYPE_BOOTH {
			if CONFIG.UnattendedInterval == 0 {
				logs.Warn("No unattended delay set for karaoke so no per-song timelapse will be made!")

				return
			}

			unattendedInterval := time.Duration(CONFIG.UnattendedInterval) * time.Second
			lastTime := time.Now()

			// This way of doing it should also be applied to
			// the photobooth module
			go func() {
				for s.Open {
					time.Sleep(1 * time.Second)
					if INSTANCE.CurrentSong == nil || !INSTANCE.Started {
						continue
					}

					currentTime := time.Now()

					if currentTime.Sub(lastTime) >= unattendedInterval {
						logs.Info("Unattended karaoke picture")
						s.Send("UNATTENDED_KARAOKE_PICTURE", nil)
						lastTime = currentTime
					}
				}
			}()
		}
	})

	return nil
}

func (m ModuleKaraoke) Initialize() error {
	if err := m.ScanSongs(); err != nil {
		fmt.Println(err)
	}

	return nil
}

// #region ScanSongs
// ScanSongs checks the song directory to add new cdg files to the database
// It also removes from the DB songs that do no longer have cdg+mp3 files
func (m ModuleKaraoke) ScanSongs() error {
	logs.Info("[Karaoke] Scanning the songs")
	baseDir, err := getModuleDir()
	if err != nil {
		return err
	}

	//#region Load folder songs UUIDs
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		logs.Error("Failed to read songs directory!", err)
		return err
	}

	folderSongs := map[string]*PhkSong{}
	folderSongsIds := []string{}

	for _, file := range entries {
		if file.IsDir() {
			continue
		}

		songFile, _ := getModuleFile(file.Name())
		song, err := LoadPhkSong(songFile)
		if err != nil {
			fmt.Printf("Failed to read file %s: %s\n", file.Name(), err)
			continue
		}

		folderSongsIds = append(folderSongsIds, song.Uuid.String())
		folderSongs[file.Name()] = song
	}
	//#endregion

	dbSongs, err := ormFetchSongUUIDs()
	if err != nil {
		logs.Error("Failed to fetch songs from DB: ", err)
		return err
	}

	//#region Removing songs that are in database but not in folders
	for _, songId := range dbSongs {
		if slices.Contains(folderSongsIds, songId) {
			continue
		}

		logs.Error("Song " + songId + " no longer available, removing it from the database.")
		ormDeleteSong(songId)
	}
	//#endregion

	//#region Then adding songs that are not in the DB yet
	for filename, song := range folderSongs {
		if slices.Contains(dbSongs, song.Uuid.String()) {
			continue
		}

		song.Filename = filename
		s, err := ormCreateSong(*song)
		if err != nil {
			logs.Error("Failed to create song "+song.Artist+" - "+song.Title+": ", err)
			continue
		}

		logs.Info(fmt.Sprintf("Song %s - %s created (ID %v).", song.Artist, song.Title, s.Id))
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
		QueueAndPlayHandler{},
		PlayingStatusHandler{},
		PlayingEndedHandler{},
		AddToQueueHandler{},
		DelFromQueueHandler{},
		PauseHandler{},
		QueueMoveUp{},
		QueueMoveDown{},
		SetVolumeInstru{},
		SetVolumeVocals{},
		SetVolumeFull{},
	}
}

func (m ModuleKaraoke) UpdateFrontendSettings() {
	queue := m.Queue
	if queue == nil {
		queue = []SongSession{}
	}

	services.GET.ModuleSettings["karaoke"] = map[string]interface{}{
		"currentSong":  m.CurrentSong,
		"queue":        queue,
		"started":      m.Started,
		"preplayTimer": m.PreplayTimer,
		"volumeInstru": m.VolumeInstru,
		"volumeVocals": m.VolumeVocals,
		"volumeFull":   m.VolumeFull,
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

	g.GET("/song/:uuid/cover", streamFileFromZip("cover.jpg", "image/jpeg"))
	g.GET("/song/:uuid/instrumental-mp3", streamFileFromZip("instrumental.mp3", "audio/mpeg"))
	g.GET("/song/:uuid/full-mp3", streamFileFromZip("full.mp3", "audio/mpeg"))
	g.GET("/song/:uuid/instrumental-webm", streamFileFromZip("instrumental.webm", "video/webm"))
	g.GET("/song/:uuid/vocals-mp3", streamFileFromZip("vocals.mp3", "audio/mpeg"))
	g.GET("/song/:uuid/cdg", streamFileFromZip("lyrics.cdg", "application/octet-stream"))

	g.POST("/picture", takePictureRoute, middlewares.BoothOnlyMiddleware)
}
