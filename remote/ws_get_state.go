package remote

import "github.com/partyhall/easyws"

type GetStateHandler struct{}

func (h GetStateHandler) GetType() string {
	return "GET_STATE"
}

func (h GetStateHandler) Do(s *easyws.Socket, payload interface{}) {
	SendState(s)
}
