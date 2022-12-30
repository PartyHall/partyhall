package mqtt_handler

import (
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
)

type ButtonPressHandler struct{}

func (h ButtonPressHandler) GetTopic() string {
	return "button_press"
}

func (h ButtonPressHandler) Do(client mqtt.Client, msg mqtt.Message) {
	// @TODO: Button press should not be in the payload
	// It should be part of the topic so that it can be different handlers

	switch string(msg.Payload()) {
	case "TAKE_PICTURE":
		socket.SOCKETS.BroadcastTakePicture()
	case "DISPLAY_DEBUG":
		logs.Debug("HW Display debug request")
		socket.SOCKETS.BroadcastState()
		socket.SOCKETS.BroadcastBooth("DISPLAY_DEBUG", nil)
	case "SHUTDOWN":
		if config.IsInDev() {
			logs.Info("[IN DEV] Shutdown requested")
			return
		}

		if err := services.GET.Shutdown(); err != nil {
			socket.SOCKETS.BroadcastTo("", "ERR_MODAL", "Failed to shutdown: "+err.Error())
		}

		err := exec.Command("shutdown", "-h", "now").Run()
		if err != nil {
			socket.SOCKETS.BroadcastTo("", "ERR_MODAL", "[MANUAL REBOOT REQUIRED] Failed to shutdown: "+err.Error())
		}
	default:
		logs.Error("Unknown button pressent: ", string(msg.Payload()))
	}

	msg.Ack()
}
