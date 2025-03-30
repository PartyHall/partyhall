package mercure_client

import "github.com/partyhall/partyhall/models"

func (mc Client) SetCurrentEvent(event *models.Event) error {
	return mc.PublishEvent("/event", event)
}
