package services

import (
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
)

// @TODO: Kick this fcking function
// its a pain to maintain, its badly engineered
// it makes cycle imports everywhere, it should not work like that
// burn it with fire
func BuildFrontendSettings() map[string]interface{} {
	if GET.CurrentState.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*GET.CurrentState.CurrentEvent)
		if err == nil {
			GET.CurrentState.CurrentEventObj = evt
		}
	}

	settings := map[string]interface{}{
		"app_state":    GET.CurrentState,
		"current_mode": GET.CurrentMode,

		"known_modes":    config.MODES,
		"guests_allowed": config.GET.GuestsAllowed,
		"modules":        GET.ModuleSettings,

		"pulseaudio_devices":  GET.PulseAudio.Devices,
		"pulseaudio_selected": GET.PulseAudio.MainDevice,
		"language":            config.GET.Language,
		"ip_addresses":        utils.GetIPs(),
		"partyhall_version":   utils.CURRENT_VERSION,
		"partyhall_commit":    utils.CURRENT_COMMIT,
	}

	//#region Adding known events
	var err error
	settings["known_events"], err = orm.GET.Events.GetEvents()
	if err != nil {
		logs.Error("Failed to get events: ", err)
		return nil
	}
	//#endregion

	return settings
}
