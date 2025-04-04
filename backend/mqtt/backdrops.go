package mqtt

import (
	"time"

	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/state"
)

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
