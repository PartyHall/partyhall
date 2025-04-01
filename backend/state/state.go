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

	UserSettings         config.UserSettings `json:"user_settings"`
	HardwareFlashPowered bool                `json:"hardware_flash_powered"`

	BackdropAlbum    *models.BackdropAlbum `json:"backdrop_album"`
	SelectedBackdrop int                   `json:"selected_backdrop"`

	Karaoke      KaraokeState          `json:"karaoke"`
	KaraokeQueue []*models.SongSession `json:"karaoke_queue"`

	SyncInProgress bool `json:"sync_in_progress"`

	GuestsAllowed bool `json:"guests_allowed"`

	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type KaraokeState struct {
	Current *models.SongSession `json:"current"`

	IsPlaying    bool `json:"is_playing"`
	Countdown    int  `json:"countdown"`
	Timecode     int  `json:"timecode"`
	Volume       int  `json:"volume"`
	VolumeVocals int  `json:"volume_vocals"`
}
