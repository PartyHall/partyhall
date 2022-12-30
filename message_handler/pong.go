package message_handler

import "github.com/partyhall/partyhall/socket"

type PongHandler struct{}

func (h PongHandler) GetType() string {
	return "PONG"
}

func (h PongHandler) Do(s *socket.Socket, payload interface{}) {
	// Nothing to do!
}
