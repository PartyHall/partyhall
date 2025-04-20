package mqtt

import (
	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/os_mgmt"
)

func OnDisplayDebug(client emqtt.Client, msg emqtt.Message) {
	log.Debug("HW Display debug request")
	mercure_client.CLIENT.ShowDebug(os_mgmt.GetCleanIps())

	msg.Ack()
}
