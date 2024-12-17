package cron

import (
	"time"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/state"
)

func RunCron() {
	// Sending time every seconds
	go func() {
		if !config.GET.SendTime {
			return
		}

		for {
			time.Sleep(1 * time.Second)
			err := mercure_client.CLIENT.PublishEvent("/time", map[string]any{
				"time": time.Now().Format(time.RFC3339),
			})

			if err != nil {
				log.Warn("Failed to publish time to mercure hub", "err", err)
			}
		}
	}()

	// Taking unattended pictures every X minutes
	go func() {
		module := config.GET.ModulesSettings.Photobooth.Unattended

		if !module.Enabled {
			return
		}

		for {
			time.Sleep(time.Duration(module.Interval) * time.Second)
			if state.STATE.CurrentEvent == nil {
				continue
			}

			err := mercure_client.CLIENT.PublishEvent("/take-picture", map[string]any{
				"unattended": true,
			})

			if err != nil {
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
