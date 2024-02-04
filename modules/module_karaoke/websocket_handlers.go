package module_karaoke

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/dto"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/remote"
)

type PlaySongHandler struct{}

func (h PlaySongHandler) GetType() string {
	return "karaoke/PLAY"
}

func (h PlaySongHandler) Do(s *easyws.Socket, payload interface{}) {
	song, err := ormLoadSongByFilename(payload.(string))
	if err != nil {
		logs.Error("Failed to load song: ", err)
		return
	}

	if INSTANCE.Queue != nil {
		toRemoveSong := -1
		for i, queueSong := range INSTANCE.Queue {
			if queueSong.Id == song.Id {
				toRemoveSong = i
				break
			}
		}
		if toRemoveSong >= 0 {
			INSTANCE.Queue = append(INSTANCE.Queue[:toRemoveSong], INSTANCE.Queue[toRemoveSong+1:]...)
		}
	}

	singer := ""
	if s.User != nil {
		user := s.User.(*jwt.Token).Claims.(*models.JwtCustomClaims)
		singer = user.Name
	}

	INSTANCE.Actions.StartSong(dto.SongDto{
		Song:   song,
		SungBy: singer,
	})
}

type PlayingStatusHandler struct{}

func (h PlayingStatusHandler) GetType() string {
	return "karaoke/PLAYING_STATUS"
}

func (h PlayingStatusHandler) Do(s *easyws.Socket, payload interface{}) {
	remote.BroadcastAdmin("karaoke/PLAYING_STATUS", payload)
}

type PlayingEndedHandler struct{}

func (h PlayingEndedHandler) GetType() string {
	return "karaoke/PLAYING_ENDED"
}

func (h PlayingEndedHandler) Do(s *easyws.Socket, payload interface{}) {
	ormCountSongPlayed(INSTANCE.CurrentSong.Filename)
	INSTANCE.Actions.StartNextSong()
}

type AddToQueueHandler struct{}

func (h AddToQueueHandler) GetType() string {
	return "karaoke/ADD_TO_QUEUE"
}

func (h AddToQueueHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.Queue == nil {
		INSTANCE.Queue = []dto.SongDto{}
	}

	song, err := ormLoadSongByFilename(payload.(string))
	if err != nil {
		logs.Error("Failed to load song: ", err)
		return
	}

	// No need to have the same song twice
	for _, queueSong := range INSTANCE.Queue {
		if queueSong.Id == song.Id {
			return
		}
	}

	singer := ""
	if s.User != nil {
		user := s.User.(*jwt.Token).Claims.(*models.JwtCustomClaims)
		singer = user.Name
	}

	INSTANCE.Queue = append(INSTANCE.Queue, dto.SongDto{
		Song:   song,
		SungBy: singer,
	})
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

type DelFromQueueHandler struct{}

func (h DelFromQueueHandler) GetType() string {
	return "karaoke/DEL_FROM_QUEUE"
}

func (h DelFromQueueHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.Queue == nil {
		INSTANCE.Queue = []dto.SongDto{}
	}

	filename := payload.(string)

	if INSTANCE.CurrentSong != nil && INSTANCE.CurrentSong.Filename == filename {
		INSTANCE.Actions.StartNextSong()
	} else {
		toRemoveSong := -1
		for i, queueSong := range INSTANCE.Queue {
			if queueSong.Filename == filename {
				toRemoveSong = i
				break
			}
		}

		if toRemoveSong >= 0 {
			INSTANCE.Queue = append(INSTANCE.Queue[:toRemoveSong], INSTANCE.Queue[toRemoveSong+1:]...)
			INSTANCE.UpdateFrontendSettings()
			remote.BroadcastState()
		}
	}
}

type PauseHandler struct{}

func (h PauseHandler) GetType() string {
	return "karaoke/PAUSE"
}

func (h PauseHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.CurrentSong == nil {
		return
	}

	INSTANCE.Started = !INSTANCE.Started
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

func findSongPosition(filename string) int {
	if INSTANCE.Queue == nil {
		return -1
	}

	for x, s := range INSTANCE.Queue {
		if s.Filename == filename {
			return x
		}
	}

	return -1
}

type QueueMoveUp struct{}

func (h QueueMoveUp) GetType() string {
	return "karaoke/QUEUE_MOVE_UP"
}

func (h QueueMoveUp) Do(s *easyws.Socket, payload interface{}) {
	idx := findSongPosition(payload.(string))

	if idx > 0 {
		tmp := INSTANCE.Queue[idx-1]
		INSTANCE.Queue[idx-1] = INSTANCE.Queue[idx]
		INSTANCE.Queue[idx] = tmp

		INSTANCE.UpdateFrontendSettings()
		remote.BroadcastState()
	}
}

type QueueMoveDown struct{}

func (h QueueMoveDown) GetType() string {
	return "karaoke/QUEUE_MOVE_DOWN"
}

func (h QueueMoveDown) Do(s *easyws.Socket, payload interface{}) {
	idx := findSongPosition(payload.(string))

	if idx < len(INSTANCE.Queue)-1 {
		tmp := INSTANCE.Queue[idx+1]
		INSTANCE.Queue[idx+1] = INSTANCE.Queue[idx]
		INSTANCE.Queue[idx] = tmp

		INSTANCE.UpdateFrontendSettings()
		remote.BroadcastState()
	}
}
