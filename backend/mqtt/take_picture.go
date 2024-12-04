package mqtt

import (
	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/state"
)

func OnTakePicture(client emqtt.Client, msg emqtt.Message) {
	if state.STATE.CurrentEvent == nil {
		log.Error("Tried to take a picture from mqtt but no event selected")
		return
	}

	log.Debug("Taking a picture from MQTT")

	err := mercure_client.CLIENT.PublishEvent("/take-picture", map[string]any{
		"unattended": false,
	})

	if err != nil {
		log.Error("Failed to publish take-picture on mercure")
	}
}
