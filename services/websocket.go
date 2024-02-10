package services

import (
	"fmt"

	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
)

// Global websocket handlers

type UpdateVolumeHandler struct{}

func (h UpdateVolumeHandler) GetType() string {
	return "SET_VOLUME"
}

func (h UpdateVolumeHandler) Do(s *easyws.Socket, payload interface{}) {
	val, ok := payload.(float64)
	if GET.PulseAudio.MainDevice != nil && ok {
		GET.PulseAudio.MainDevice.SetVolume(int(val))
		s.Send("SET_SOUND_CARD", GET.PulseAudio.MainDevice)
	}
}

type SetMuteHandler struct{}

func (h SetMuteHandler) GetType() string {
	return "SET_MUTE"
}

func (h SetMuteHandler) Do(s *easyws.Socket, payload interface{}) {
	val, ok := payload.(bool)
	if !ok {
		return
	}

	GET.PulseAudio.MainDevice.SetMute(val)
	if err := GET.PulseAudio.Refresh(); err != nil {
		logs.Error(err)
	}
	fmt.Println(GET.PulseAudio.MainDevice.Mute)
	s.Send("SET_SOUND_CARD", GET.PulseAudio.MainDevice)
}

type SetSoundcardHandler struct{}

func (h SetSoundcardHandler) GetType() string {
	return "SET_SOUND_CARD"
}

func (h SetSoundcardHandler) Do(s *easyws.Socket, payload interface{}) {
	val, ok := payload.(float64)
	if !ok {
		return
	}

	valInt := int(val)

	for _, x := range GET.PulseAudio.Devices {
		if x.Index == valInt {
			GET.PulseAudio.SetDefaultOutput(x)
			GET.PulseAudio.Refresh()
			s.Send("SOUND_DEVICE", GET.PulseAudio.MainDevice)
		}
	}
}
