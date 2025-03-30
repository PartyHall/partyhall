package mercure_client

import "github.com/partyhall/partyhall/state"

func (mc Client) SendTakePicture(unattended bool) error {
	return mc.PublishEvent("/take-picture", map[string]any{
		"unattended": unattended,
	})
}

func (mc Client) SetFlash(powered bool, brightness int) error {
	return mc.PublishEvent(
		"/flash",
		map[string]any{
			"powered":    powered,
			"brightness": brightness,
		},
	)
}

func (mc Client) SendBackdropState() error {
	return mc.PublishEvent("/backdrop-state", map[string]any{
		"backdrop_album":    state.STATE.BackdropAlbum,
		"selected_backdrop": state.STATE.SelectedBackdrop,
	})
}
