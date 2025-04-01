package mqtt

import (
	"fmt"

	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
)

func SetFlash(powered bool, brightness int) {
	brightness = utils.ClampInt(brightness, 0, 100)
	if brightness == 0 {
		powered = false
	}

	state.STATE.HardwareFlashPowered = powered

	mercure_client.CLIENT.SetFlash(
		powered,
		config.GET.UserSettings.Photobooth.FlashBrightness,
	)

	err := EasyMqtt.Send("partyhall/flash", fmt.Sprintf("%v", brightness))
	if err != nil {
		api_errors.MQTT_PUBLISH_FAILURE.WithExtra(map[string]any{
			"err": err,
		})

		log.Error("[MQTT] Faied to publish partyhall/flash", "err", err)

		return
	}
}
