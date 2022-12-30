package mqtt_handler

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
	"golang.org/x/exp/slices"
)

type AdminSetModeHandler struct{}

func (h AdminSetModeHandler) GetTopic() string {
	return "admin/set_mode"
}

func (h AdminSetModeHandler) Do(client mqtt.Client, msg mqtt.Message) {
	mode := string(msg.Payload())
	if !slices.Contains(config.MODES, mode) {
		logs.Error("given mode is not allowed")
		return
	}

	services.GET.CurrentMode = mode
	socket.SOCKETS.BroadcastState()
}
