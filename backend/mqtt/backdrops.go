package mqtt

import (
	"time"

	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/state"
)

func OnNextBackdropAlbum(client emqtt.Client, msg emqtt.Message) {
	defer msg.Ack()

	log.Info("Going to next backdrop album")

	albums, err := dal.BACKDROPS.GetAllAlbums()
	if err != nil {
		log.Error("Failed to retreive backdrop albums", "err", err)
		return
	}

	selected := -1
	if state.STATE.BackdropAlbum != nil {
		for x, alb := range albums {
			if alb.Id == state.STATE.BackdropAlbum.Id {
				selected = x
			}
		}

		selected = selected + 1

		if selected > len(albums)-1 {
			selected = -1
		}
	} else if len(albums) > 0 {
		selected = 0
	} else {
		return
	}

	if selected > -1 {
		album, err := dal.BACKDROPS.GetAlbum(albums[selected].Id)
		if err != nil {
			log.Error("Failed to load backdrop album", "err", err)
			return
		}

		state.STATE.BackdropAlbum = &album
		log.Info("New backdrop album selected", "id", album.Id, "name", album.Name)
	} else {
		state.STATE.BackdropAlbum = nil
		log.Info("Looping to no backdrop album selected")
	}

	state.STATE.SelectedBackdrop = 0
	state.STATE.BackdropSelectedAt = time.Now()

	mercure_client.CLIENT.SendBackdropState()
}

func OnNextBackdrop(client emqtt.Client, msg emqtt.Message) {
	defer msg.Ack()

	if state.STATE.BackdropAlbum == nil {
		return
	}

	state.STATE.SelectedBackdrop++
	if state.STATE.SelectedBackdrop > len(state.STATE.BackdropAlbum.Backdrops) {
		state.STATE.SelectedBackdrop = 0
	}

	state.STATE.BackdropSelectedAt = time.Now()

	mercure_client.CLIENT.SendBackdropState()
}
