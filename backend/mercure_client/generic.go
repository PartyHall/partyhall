package mercure_client

import (
	"time"

	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/utils"
)

func (mc Client) SendTime() error {
	return mc.PublishEvent("/time", map[string]string{
		"time": time.Now().Format(time.RFC3339),
	})
}

func (mc Client) SetMode(mode string) error {
	return mc.PublishEvent("/mode", map[string]string{
		"mode": mode,
	})
}

func (mc Client) ShowDebug() error {
	return mc.PublishEvent("/debug", map[string]any{
		"ip_addresses": utils.GetIPs(),
	})
}

func (mc Client) ShowSnackbar(snackType, msg string) error {
	return mc.PublishEvent("/snackbar", map[string]any{
		"type": snackType,
		"msg":  msg,
	})
}

func (mc Client) SendLog(log models.Log) error {
	return mc.PublishEvent("/logs", log)
}
