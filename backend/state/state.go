package state

import (
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
)

const MODE_DISABLED = "disabled"
const MODE_PHOTOBOOTH = "photobooth"

var MODES = []string{
	MODE_DISABLED,
	MODE_PHOTOBOOTH,
}

var STATE State

type State struct {
	CurrentEvent *models.Event       `json:"current_event"`
	CurrentMode  string              `json:"current_mode"`
	IpAddresses  map[string][]string `json:"ip_addresses"`

	ModulesSettings config.ModulesSettings `json:"modules_settings"`

	Karaoke      KaraokeState          `json:"karaoke"`
	KaraokeQueue []*models.SongSession `json:"karaoke_queue"`

	SyncInProgress bool `json:"sync_in_progress"`

	GuestsAllowed bool `json:"guests_allowed"`

	HardwareId string `json:"hwid"`
	Version    string `json:"version"`
	Commit     string `json:"commit"`
}

type KaraokeState struct {
	Current *models.SongSession `json:"current"`

	IsPlaying    bool `json:"is_playing"`
	Countdown    int  `json:"countdown"`
	Timecode     int  `json:"timecode"`
	Volume       int  `json:"volume"`
	VolumeVocals int  `json:"volume_vocals"`
}
