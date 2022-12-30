package mqtt_handler

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/logs"
)

type SyncRequestedHandler struct{}

func (h SyncRequestedHandler) GetTopic() string {
	return "sync"
}

func (h SyncRequestedHandler) Do(client mqtt.Client, msg mqtt.Message) {
	logs.Error("Sync not implemented !")
	msg.Ack()
}
