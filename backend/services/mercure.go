package services

import (
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/utils"
)

func ShowDebug() error {
	err := mercure_client.CLIENT.PublishEvent("/debug", map[string]any{
		"ip_addresses": utils.GetIPs(),
	})
	if err != nil {
		log.LOG.Error("Failed to publish show_debug message", "err", err)
		return err
	}

	return nil
}

func ShowSnackbar(snackType, msg string) error {
	err := mercure_client.CLIENT.PublishEvent("/snackbar", map[string]any{
		"type": snackType,
		"msg":  msg,
	})
	if err != nil {
		log.LOG.Error("Failed to publish show_debug message", "err", err)
		return err
	}

	return nil
}
