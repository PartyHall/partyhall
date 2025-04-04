package state

import (
	"time"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
)

const MODE_DISABLED = "disabled"
const MODE_PHOTOBOOTH = "photobooth"
const MODE_BTN_SETUP = "btn_setup"

var MODES = []string{
	MODE_DISABLED,
	MODE_PHOTOBOOTH,
	MODE_BTN_SETUP,
}

var STATE State

type State struct {
	CurrentEvent *models.Event       `json:"current_event"`
	CurrentMode  string              `json:"current_mode"`
	previousMode string              `json:"-"`
	ModeSetAt    time.Time           `json:"-"`
	IpAddresses  map[string][]string `json:"ip_addresses"`

	UserSettings         config.UserSettings `json:"user_settings"`
	HardwareFlashPowered bool                `json:"hardware_flash_powered"`

	BackdropAlbum      *models.BackdropAlbum `json:"backdrop_album"`
	BackdropSelectedAt time.Time             `json:"-"`
	SelectedBackdrop   int                   `json:"selected_backdrop"`

	Karaoke      KaraokeState          `json:"karaoke"`
	KaraokeQueue []*models.SongSession `json:"karaoke_queue"`

	SyncInProgress bool `json:"sync_in_progress"`

	GuestsAllowed bool `json:"guests_allowed"`

	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func (s *State) SetMode(mode string) {
	if s.CurrentMode != mode {
		s.CurrentMode = mode
		s.previousMode = s.CurrentMode
	}

	s.ModeSetAt = time.Now()
}

type KaraokeState struct {
	Current *models.SongSession `json:"current"`

	IsPlaying    bool `json:"is_playing"`
	Countdown    int  `json:"countdown"`
	Timecode     int  `json:"timecode"`
	Volume       int  `json:"volume"`
	VolumeVocals int  `json:"volume_vocals"`
}
