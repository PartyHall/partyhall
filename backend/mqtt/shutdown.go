package mqtt

import (
	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/services"
)

func OnShutdown(client emqtt.Client, msg emqtt.Message) {
	msg.Ack()

	if config.GET.IsInDev {
		log.Info("[IN DEV] Shutdown requested")
		return
	}

	if err := services.Shutdown(); err != nil {
		mercure_client.CLIENT.ShowSnackbar("error", "Failed to shutdown: "+err.Error())
	}
}
