package services

import (
	"time"

	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
)

var KARAOKE = Karaoke{}

type Karaoke struct{}

func (k Karaoke) EndCurrentSong(publish bool) {
	if state.STATE.Karaoke.Current == nil {
		return
	}

	if state.STATE.Karaoke.Current.CancelledAt.Valid || state.STATE.Karaoke.Current.EndedAt.Valid {
		if publish {
			mercure_client.CLIENT.PublishEvent(
				"/karaoke",
				state.STATE.Karaoke,
			)
		}

		return
	}

	state.STATE.Karaoke.Current.CancelledAt = models.JsonnableNullTime{Time: time.Now(), Valid: true}
	dal.SONGS.UpdateSession(state.STATE.Karaoke.Current)

	state.STATE.Karaoke.Current = nil

	if publish {
		mercure_client.CLIENT.PublishEvent(
			"/karaoke",
			state.STATE.Karaoke,
		)
	}
}

func (k Karaoke) StartSong(session *models.SongSession) error {
	k.EndCurrentSong(false)

	state.STATE.Karaoke.Current = session
	state.STATE.Karaoke.Countdown = 5
	state.STATE.Karaoke.Timecode = 0
	state.STATE.Karaoke.IsPlaying = false

	go func() {
		for state.STATE.Karaoke.Countdown > 0 {
			time.Sleep(1 * time.Second)
			state.STATE.Karaoke.Countdown--

			if state.STATE.Karaoke.Countdown == 0 {
				state.STATE.Karaoke.IsPlaying = true

				state.STATE.Karaoke.Current.StartedAt = models.JsonnableNullTime{
					Valid: true,
					Time:  time.Now(),
				}
				dal.SONGS.UpdateSession(state.STATE.Karaoke.Current)
			}

			mercure_client.CLIENT.PublishEvent(
				"/karaoke",
				state.STATE.Karaoke,
			)
		}
	}()

	return nil
}

func (k Karaoke) StartNextSong() error {
	if len(state.STATE.KaraokeQueue) > 0 {
		nextSong := state.STATE.KaraokeQueue[0]
		state.STATE.KaraokeQueue = state.STATE.KaraokeQueue[1:]

		err := k.StartSong(nextSong)
		if err != nil {
			log.Error("Failed to start next song", "err", err)
		}
	} else {
		k.EndCurrentSong(false)
	}

	mercure_client.CLIENT.PublishEvent(
		"/karaoke-queue",
		state.STATE.KaraokeQueue,
	)

	mercure_client.CLIENT.PublishEvent(
		"/karaoke",
		state.STATE.Karaoke,
	)

	return nil
}
