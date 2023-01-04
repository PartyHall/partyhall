package module_photobooth

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/remote"
)

func OnTakePicture(client mqtt.Client, msg mqtt.Message) {
	for _, s := range remote.EasyWS.Sockets {
		INSTANCE.Actions.TakePicture(s)
	}

	msg.Ack()
}
