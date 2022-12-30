package message_handler

import (
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
)

type SetEventHandler struct{}

func (h SetEventHandler) GetType() string {
	return "SET_EVENT"
}

func (h SetEventHandler) Do(s *socket.Socket, payload interface{}) {
	evtIdFloat, ok := payload.(float64)
	if !ok {
		s.Send("ERR_MODAL", "Failed to change event: Bad request")
		return
	}

	var evtId int64 = int64(evtIdFloat)

	evt, err := orm.GET.Events.GetEvent(evtId)
	if err != nil {
		s.Send("ERR_MODAL", "Failed to change event: "+err.Error())
		return
	}

	services.GET.CurrentState.CurrentEvent = &evtId
	services.GET.CurrentState.CurrentEventObj = evt

	err = orm.GET.AppState.SetState(services.GET.CurrentState)
	if err != nil {
		s.Send("ERR_MODAL", "Failed to save state: "+err.Error())
		return
	}

	socket.SOCKETS.BroadcastState()
}
