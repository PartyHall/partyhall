package remote

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/logs"
)

func OnDisplayDebug(client mqtt.Client, msg mqtt.Message) {
	logs.Debug("HW Display debug request")
	BroadcastState()
	BroadcastBooth("DISPLAY_DEBUG", nil)

	msg.Ack()
}
