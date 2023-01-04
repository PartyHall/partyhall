package remote

import "github.com/partyhall/easyws"

type PongHandler struct{}

func (h PongHandler) GetType() string {
	return "PONG"
}

func (h PongHandler) Do(s *easyws.Socket, payload interface{}) {
	// Nothing to do!
}
