package cron

import (
	"time"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/state"
)

var lastUnattendedPicture time.Time

/**
 * Note:
 * Cron method runs at all time and diff with the previous run time
 * because we want them to update snappily when they change their
 * settings
 **/

func RunCron() {
	// Sending time every seconds
	go func() {
		if !config.GET.SendTime {
			return
		}

		for {
			time.Sleep(1 * time.Second)

			if err := mercure_client.CLIENT.SendTime(); err != nil {
				log.Warn("Failed to publish time to mercure hub", "err", err)
			}
		}
	}()

	// Taking unattended pictures every X minutes
	lastUnattendedPicture = time.Now()
	go func() {
		for {
			time.Sleep(1 * time.Second)

			module := config.GET.UserSettings.Photobooth.Unattended

			if !module.Enabled || state.STATE.CurrentEvent == nil {
				continue
			}

			diff := time.Since(lastUnattendedPicture)
			if diff <= (time.Duration(module.Interval) * time.Second) {
				continue
			}

			lastUnattendedPicture = time.Now()
			if err := mercure_client.CLIENT.SendTakePicture(true); err != nil {
				log.Warn("Failed to publish take unattended to mercure hub", "err", err)
			}
		}
	}()

	// Sync-ing songs, images, and sessions every 5 minutes
	go func() {
		for {
			if !state.STATE.SyncInProgress {
				err := nexus.INSTANCE.Sync(state.STATE.CurrentEvent)
				if err != nil {
					mercure_client.CLIENT.PublishSyncInProgress()
					log.Error("Failed to sync songs", "err", err)
				}
			} else {
				log.Info("CRON Synchronizing has been skipped as it's already in progress")
			}

			time.Sleep(5 * time.Minute)
		}
	}()
}
