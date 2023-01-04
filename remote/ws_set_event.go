package remote

import (
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
)

type SetEventHandler struct{}

func (h SetEventHandler) GetType() string {
	return "SET_EVENT"
}

func (h SetEventHandler) Do(s *easyws.Socket, payload interface{}) {
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

	BroadcastState()
}
