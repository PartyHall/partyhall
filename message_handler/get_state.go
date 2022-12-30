package message_handler

import "github.com/partyhall/partyhall/socket"

type GetStateHandler struct{}

func (h GetStateHandler) GetType() string {
	return "GET_STATE"
}

func (h GetStateHandler) Do(s *socket.Socket, payload interface{}) {
	s.SendState()
}
