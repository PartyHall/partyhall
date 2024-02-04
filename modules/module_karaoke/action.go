package module_karaoke

import (
	"time"

	"github.com/partyhall/partyhall/dto"
	"github.com/partyhall/partyhall/remote"
)

type Actions struct{}

func (a Actions) StartSong(song dto.SongDto) {
	INSTANCE.CurrentSong = &song
	INSTANCE.PreplayTimer = CONFIG.PrePlayTimer
	INSTANCE.Started = true
	INSTANCE.UpdateFrontendSettings()
	remote.BroadcastState()

	go func() {
		for INSTANCE.PreplayTimer > 0 {
			time.Sleep(1 * time.Second)
			INSTANCE.PreplayTimer -= 1
			INSTANCE.UpdateFrontendSettings()
			remote.BroadcastState()
		}
	}()
}

func (a Actions) StartNextSong() {
	if INSTANCE.Queue != nil && len(INSTANCE.Queue) > 0 {
		nextSong := INSTANCE.Queue[0]
		INSTANCE.Queue = INSTANCE.Queue[1:]

		a.StartSong(nextSong)
	} else {
		INSTANCE.CurrentSong = nil
		INSTANCE.UpdateFrontendSettings()
		remote.BroadcastState()
	}
}
