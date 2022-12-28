package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"golang.org/x/exp/slices"
)

const SOCKET_TYPE_BOOTH = "BOOTH"
const SOCKET_TYPE_ADMIN = "ADMIN"

var SOCKET_TYPES = []string{
	SOCKET_TYPE_BOOTH,
	SOCKET_TYPE_ADMIN,
}

type Sockets []*Socket

func (s Sockets) broadcastTo(to string, msgType string, data interface{}) {
	for _, socket := range s {
		if len(to) > 0 && socket.Type != to {
			continue
		}

		socket.Send(msgType, data)
	}
}

func (s Sockets) BroadcastPartyHall(msgType string, data interface{}) {
	s.broadcastTo(SOCKET_TYPE_BOOTH, msgType, data)
}

func (s Sockets) BroadcastTakePicture() {
	for _, socket := range s {
		if socket.Type != SOCKET_TYPE_BOOTH {
			continue
		}

		socket.TakePicture()
	}
}

func (s Sockets) BroadcastAdmin(msgType string, data interface{}) {
	s.broadcastTo(SOCKET_TYPE_ADMIN, msgType, data)
}

func (s Sockets) BroadcastState() {
	for _, socket := range s {
		socket.sendState()
	}
}

type Socket struct {
	Type string
	Open bool
	Conn *websocket.Conn

	mtx *sync.Mutex
}

func (s *Socket) TakePicture() {
	if GET.PartyHall.IsTakingPicture || s.Type != SOCKET_TYPE_BOOTH {
		return
	}

	GET.PartyHall.IsTakingPicture = true

	logs.Info("Taking picture...")
	go func() {
		timeout := config.GET.PartyHall.DefaultTimer

		for timeout >= 0 {
			s.Send("TIMER", timeout)
			timeout--
			time.Sleep(1 * time.Second)
		}

		GET.PartyHall.IsTakingPicture = false
	}()
}

func (s *Socket) OnMessage(msg models.SocketMessage) {
	switch msg.MsgType {
	case "PONG":
		break
	case "GET_STATE":
		s.sendState()
	case "TAKE_PICTURE":
		s.TakePicture()
	case "REMOTE_TAKE_PICTURE":
		for _, sock := range GET.Sockets {
			sock.TakePicture()
		}
	case "SET_MODE":
		mode, ok := msg.Payload.(string)
		if !ok {
			s.Send("ERR_MODAL", "Bad request")
			break
		}

		if !slices.Contains(config.MODES, mode) {
			s.Send("ERR_MODAL", "Unknown mode")
			break
		}

		GET.PartyHall.CurrentMode = mode
		GET.Sockets.BroadcastState()
	case "SET_DATETIME":
		dt, ok := msg.Payload.(string)
		if !ok {
			s.Send("ERR_MODAL", "Bad request")
			break
		}

		time, err := time.Parse("2006-01-02 15:04:05", dt)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to set date: "+err.Error())
			break
		}
		err = SetSystemDate(time)
		if err != nil {
			logs.Error(err)
		}
	case "SET_EVENT":
		evtIdFloat, ok := msg.Payload.(float64)
		if !ok {
			s.Send("ERR_MODAL", "Failed to change event: Bad request")
			break
		}

		var evtId int64 = int64(evtIdFloat)

		evt, err := orm.GET.Events.GetEvent(evtId)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to change event: "+err.Error())
			break
		}

		GET.PartyHall.CurrentState.CurrentEvent = &evtId
		GET.PartyHall.CurrentState.CurrentEventObj = evt

		err = orm.GET.AppState.SetState(GET.PartyHall.CurrentState)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to save state: "+err.Error())
			break
		}

		GET.Sockets.BroadcastState()
	case "SHOW_DEBUG":
		(*GET.MqttClient).Publish("partyhall/button_press", 2, false, "DISPLAY_DEBUG")
	case "EXPORT_ZIP":
		token := (*GET.MqttClient).Publish("partyhall/export", 2, false, fmt.Sprintf("%v", msg.Payload))
		token.Wait()
		if token.Error() != nil {
			logs.Error(token.Error())
		}
	case "SHUTDOWN":
		(*GET.MqttClient).Publish("partyhall/button_press", 2, false, "SHUTDOWN")
	case "":
		// Probably should be handled in another way
		return
	default:
		logs.Infof("Unhandled socket message: %v => ", msg.MsgType, msg)
	}
}

func (s *Socket) Send(msgType string, data interface{}) error {
	s.mtx.Lock()
	err := s.Conn.WriteJSON(models.SocketMessage{
		MsgType: msgType,
		Payload: data,
	})
	s.mtx.Unlock()

	return err
}

func (s *Socket) sendState() error {
	settings := GET.GetFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	return s.Send("APP_STATE", settings)
}

func (p *Provider) Join(socketType string, socket *websocket.Conn) {
	sock := &Socket{
		Type: socketType,
		Conn: socket,
		Open: true,
		mtx:  &sync.Mutex{},
	}
	p.Sockets = append(p.Sockets, sock)

	go func() {
		for sock.Open {
			time.Sleep(1 * time.Second)
			sock.Send("PING", time.Now().Format("2006-01-02 15:04:05"))
		}
	}()

	go func() {
		for {
			data := models.SocketMessage{}
			err := socket.ReadJSON(&data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logs.Error("Unexpected close error: ", err)
					sock.Open = false
					return
				} else if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logs.Error("Websocket disconnected: ", err)
					sock.Open = false
					return
				}

				logs.Error(err)
				continue
			}

			sock.OnMessage(data)
		}
	}()

	if socketType == SOCKET_TYPE_BOOTH && config.GET.PartyHall.UnattendedInterval > 1 {
		go func() {
			for sock.Open {
				time.Sleep(time.Duration(config.GET.PartyHall.UnattendedInterval) * time.Minute)
				logs.Info("Unattended picture")
				sock.Send("UNATTENDED_PICTURE", nil)
			}
		}()
	}

	sock.sendState()
}
