package mqtt

import (
	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/services"
)

func OnDisplayDebug(client emqtt.Client, msg emqtt.Message) {
	log.Debug("HW Display debug request")
	services.ShowDebug()

	msg.Ack()
}
