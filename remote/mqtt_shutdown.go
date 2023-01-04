package remote

import (
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
)

func OnShutdown(client mqtt.Client, msg mqtt.Message) {
	msg.Ack()

	if config.IsInDev() {
		logs.Info("[IN DEV] Shutdown requested")
		return
	}

	if err := services.GET.Shutdown(); err != nil {
		Broadcast("ERR_MODAL", "Failed to shutdown: "+err.Error())
		return
	}

	err := exec.Command("shutdown", "-h", "now").Run()
	if err != nil {
		Broadcast("ERR_MODAL", "Failed to shutdown: "+err.Error())
	}
}
