package mercure_client

import "github.com/partyhall/partyhall/models"

func (mc Client) SendAudioDevices(devices *models.PwDevices) error {
	return mc.PublishEvent("/audio-devices", devices)
}
