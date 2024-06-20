package module_karaoke

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
)

type PlaySongHandler struct{}

func (h PlaySongHandler) GetType() string {
	return "karaoke/PLAY"
}

// PLAY = Session ID OK
func (h PlaySongHandler) Do(s *easyws.Socket, payload interface{}) {
	sessionIdFloat, ok := payload.(float64)
	if !ok {
		logs.Error("Failed to play song: Session ID could not be parsed")
		return
	}

	sessionId := int(sessionIdFloat)

	var songSession *SongSession = nil
	toRemoveSong := -1

	if INSTANCE.CurrentSong != nil && INSTANCE.CurrentSong.Id == sessionId {
		songSession = INSTANCE.CurrentSong
	} else if INSTANCE.Queue != nil {
		for idx, session := range INSTANCE.Queue {
			if session.Id == sessionId {
				songSession = &session
				toRemoveSong = idx
				break
			}
		}
	}

	if INSTANCE.CurrentSong != nil {
		now := models.Timestamp(time.Now())
		INSTANCE.CurrentSong.CancelledAt = &now
		_, err := ormUpsertSession(0, *INSTANCE.CurrentSong)
		if err != nil {
			logs.Errorf("Failed to cancel previous song (%v): %v", INSTANCE.CurrentSong.Id, err)
		}
	}

	if songSession == nil {
		logs.Errorf("Failed to play song: Session %v is not in the queue!", sessionId)
		return
	}

	if toRemoveSong >= 0 {

		INSTANCE.Queue = append(INSTANCE.Queue[:toRemoveSong], INSTANCE.Queue[toRemoveSong+1:]...)
	}

	INSTANCE.Actions.StartSong(*songSession)
}

type PlayingStatusHandler struct{}

func (h PlayingStatusHandler) GetType() string {
	return "karaoke/PLAYING_STATUS"
}

// PLAYING_STATUS: Session ID OK
func (h PlayingStatusHandler) Do(s *easyws.Socket, payload interface{}) {
	remote.BroadcastAdmin("karaoke/PLAYING_STATUS", payload)
}

type PlayingEndedHandler struct{}

func (h PlayingEndedHandler) GetType() string {
	return "karaoke/PLAYING_ENDED"
}

// PLAYING_ENDED: Session ID OK
func (h PlayingEndedHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.CurrentSong != nil {
		now := models.Timestamp(time.Now())
		INSTANCE.CurrentSong.EndedAt = &now
		_, err := ormUpsertSession(0, *INSTANCE.CurrentSong)
		if err != nil {
			logs.Errorf("Failed to upsert song session %v: %v", INSTANCE.CurrentSong.Id, err)
		}
	}

	INSTANCE.Actions.StartNextSong()
}

type AddToQueueHandler struct{}

func (h AddToQueueHandler) GetType() string {
	return "karaoke/ADD_TO_QUEUE"
}

// ADD_TO_QUEUE: Song UUID OK
func (h AddToQueueHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.Queue == nil {
		INSTANCE.Queue = []SongSession{}
	}

	song, err := ormLoadSongByUuid(payload.(string))
	if err != nil {
		logs.Error("Failed to load song: ", err)
		return
	}

	// No need to have the same song twice
	for _, queueSong := range INSTANCE.Queue {
		if queueSong.Song.Id == song.Id {
			return
		}
	}

	singer := ""
	if s.User != nil {
		user := s.User.(*jwt.Token).Claims.(*models.JwtCustomClaims)
		singer = user.Name
	}

	eventId := 0
	if services.GET.CurrentState.CurrentEvent != nil {
		eventId = *services.GET.CurrentState.CurrentEvent
	}

	session, err := ormUpsertSession(eventId, SongSession{
		Song:   *song,
		SungBy: singer,
	})
	if err != nil {
		logs.Errorf("Failed to upsert session: %v", err)
	}

	INSTANCE.Queue = append(INSTANCE.Queue, *session)
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

type QueueAndPlayHandler struct{}

func (h QueueAndPlayHandler) GetType() string {
	return "karaoke/QUEUE_AND_PLAY"
}

// QUEUE_AND_PLAY: Song UUID
func (h QueueAndPlayHandler) Do(s *easyws.Socket, payload interface{}) {
	song, err := ormLoadSongByUuid(payload.(string))
	if err != nil {
		logs.Error("Failed to load song: ", err)
		return
	}

	singer := ""
	if s.User != nil {
		user := s.User.(*jwt.Token).Claims.(*models.JwtCustomClaims)
		singer = user.Name
	}

	eventId := 0
	if services.GET.CurrentState.CurrentEvent != nil {
		eventId = *services.GET.CurrentState.CurrentEvent
	}

	session, err := ormUpsertSession(eventId, SongSession{
		Song:   *song,
		SungBy: singer,
	})
	if err != nil {
		logs.Errorf("Failed to upsert session: %v", err)
	}

	INSTANCE.Actions.StartSong(*session)
}

type DelFromQueueHandler struct{}

func (h DelFromQueueHandler) GetType() string {
	return "karaoke/DEL_FROM_QUEUE"
}

// DEL_FROM_QUEUE: Session ID OK
func (h DelFromQueueHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.Queue == nil {
		INSTANCE.Queue = []SongSession{}
	}

	sessionIdFloat, ok := payload.(float64)
	if !ok {
		logs.Error("Failed to play song: Session ID could not be parsed")
		return
	}

	sessionId := int(sessionIdFloat)

	var session *SongSession = nil
	toRemoveSong := -1

	if INSTANCE.CurrentSong != nil && sessionId == INSTANCE.CurrentSong.Id {
		session = INSTANCE.CurrentSong
	} else {
		for i, queueSong := range INSTANCE.Queue {
			if queueSong.Id == sessionId {
				toRemoveSong = i
				session = &queueSong
				break
			}
		}
	}

	if session != nil {
		now := models.Timestamp(time.Now())
		session.CancelledAt = &now
		_, err := ormUpsertSession(0, *session)
		if err != nil {
			logs.Errorf("Failed to save the song cancellation for session %v: %v", session.Id, err)
		}
	}

	if INSTANCE.CurrentSong != nil && INSTANCE.CurrentSong.Id == sessionId {
		INSTANCE.Actions.StartNextSong()
	} else if toRemoveSong >= 0 {
		INSTANCE.Queue = append(INSTANCE.Queue[:toRemoveSong], INSTANCE.Queue[toRemoveSong+1:]...)
		INSTANCE.UpdateFrontendSettings()
		remote.BroadcastState()
	}
}

type PauseHandler struct{}

func (h PauseHandler) GetType() string {
	return "karaoke/PAUSE"
}

// PAUSE: Session ID
func (h PauseHandler) Do(s *easyws.Socket, payload interface{}) {
	if INSTANCE.CurrentSong == nil {
		return
	}

	INSTANCE.Started = !INSTANCE.Started
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

func findSongPosition(sessionId int) int {
	if INSTANCE.Queue == nil {
		return -1
	}

	for x, s := range INSTANCE.Queue {
		if s.Id == sessionId {
			return x
		}
	}

	return -1
}

type QueueMoveUp struct{}

func (h QueueMoveUp) GetType() string {
	return "karaoke/QUEUE_MOVE_UP"
}

// QUEUE_MOVE_UP: Session ID OK
func (h QueueMoveUp) Do(s *easyws.Socket, payload interface{}) {
	sessionIdFloat, ok := payload.(float64)
	if !ok {
		logs.Error("Failed to play song: Session ID could not be parsed")
		return
	}

	sessionId := int(sessionIdFloat)

	idx := findSongPosition(sessionId)

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

// QUEUE_MOVE_DOWN: Session ID OK
func (h QueueMoveDown) Do(s *easyws.Socket, payload interface{}) {
	sessionIdFloat, ok := payload.(float64)
	if !ok {
		logs.Error("Failed to play song: Session ID could not be parsed")
		return
	}

	sessionId := int(sessionIdFloat)

	idx := findSongPosition(sessionId)

	if idx < len(INSTANCE.Queue)-1 {
		tmp := INSTANCE.Queue[idx+1]
		INSTANCE.Queue[idx+1] = INSTANCE.Queue[idx]
		INSTANCE.Queue[idx] = tmp

		INSTANCE.UpdateFrontendSettings()
		remote.BroadcastState()
	}
}

type SetVolumeVocals struct{}

func (h SetVolumeVocals) GetType() string {
	return "karaoke/VOLUME_VOCALS"
}

func (h SetVolumeVocals) Do(s *easyws.Socket, payload interface{}) {
	vol, ok := payload.(float64)
	if !ok {
		fmt.Println("bad value:", payload)
		return
	}

	if vol < 0 {
		vol = 0
	} else if vol > 1 {
		vol = 1
	}

	INSTANCE.VolumeVocals = vol
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

type SetVolumeInstru struct{}

func (h SetVolumeInstru) GetType() string {
	return "karaoke/VOLUME_INSTRU"
}

func (h SetVolumeInstru) Do(s *easyws.Socket, payload interface{}) {
	vol, ok := payload.(float64)
	if !ok {
		fmt.Println("bad value:", payload)
		return
	}

	if vol < 0 {
		vol = 0
	} else if vol > 1 {
		vol = 1
	}

	INSTANCE.VolumeInstru = vol
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}

type SetVolumeFull struct{}

func (h SetVolumeFull) GetType() string {
	return "karaoke/VOLUME_FULL"
}

func (h SetVolumeFull) Do(s *easyws.Socket, payload interface{}) {
	vol, ok := payload.(float64)
	if !ok {
		fmt.Println("bad value:", ok)
		return
	}

	if vol < 0 {
		vol = 0
	} else if vol > 1 {
		vol = 1
	}

	INSTANCE.VolumeFull = vol
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()
}
