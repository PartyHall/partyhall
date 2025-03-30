package mercure_client

import "github.com/partyhall/partyhall/state"

func (mc Client) SendUserSettings() error {
	return mc.PublishEvent("/user-settings", state.STATE.UserSettings)
}
