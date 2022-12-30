package message_handler

import (
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
	"golang.org/x/exp/slices"
)

type SetModeHandler struct{}

func (h SetModeHandler) GetType() string {
	return "SET_MODE"
}

func (h SetModeHandler) Do(s *socket.Socket, payload interface{}) {
	mode, ok := payload.(string)
	if !ok {
		s.Send("ERR_MODAL", "Bad request")
		return
	}

	if !slices.Contains(config.MODES, mode) {
		s.Send("ERR_MODAL", "Unknown mode")
		return
	}

	services.GET.CurrentMode = mode
	socket.SOCKETS.BroadcastState()
}
