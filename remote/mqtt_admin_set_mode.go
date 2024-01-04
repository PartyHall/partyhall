package remote

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"golang.org/x/exp/slices"
)

func OnSetMode(client mqtt.Client, msg mqtt.Message) {
	mode := string(msg.Payload())
	if !slices.Contains(config.MODES, mode) {
		logs.Error("given mode is not allowed")
		return
	}

	services.GET.CurrentMode = mode
	BroadcastState()
}
