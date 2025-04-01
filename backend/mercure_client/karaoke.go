package mercure_client

import (
	"github.com/partyhall/partyhall/state"
)

func (mc Client) SendKaraokeState() error {
	return mc.PublishEvent("/karaoke", state.STATE.Karaoke)
}

func (mc Client) SendKaraokeQueue() error {
	return mc.PublishEvent("/karaoke-queue", state.STATE.KaraokeQueue)
}

func (mc Client) SendKaraokeTimecode() error {
	return mc.PublishEvent("/karaoke-timecode", map[string]any{
		"timecode": state.STATE.Karaoke.Timecode,
	})
}
