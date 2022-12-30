package socket

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
)

const SOCKET_TYPE_BOOTH = "BOOTH"
const SOCKET_TYPE_ADMIN = "ADMIN"

var SOCKET_TYPES = []string{
	SOCKET_TYPE_BOOTH,
	SOCKET_TYPE_ADMIN,
}

type Socket struct {
	Type string
	Open bool
	Conn *websocket.Conn

	mtx *sync.Mutex
}

func Join(socketType string, socket *websocket.Conn) *Socket {
	sock := &Socket{
		Type: socketType,
		Conn: socket,
		Open: true,
		mtx:  &sync.Mutex{},
	}

	SOCKETS = append(SOCKETS, sock)

	sock.StartPings()
	sock.StartUnattended()
	sock.SendState()

	return sock
}

func (s *Socket) TakePicture() {
	if services.GET.IsTakingPicture || s.Type != SOCKET_TYPE_BOOTH {
		return
	}

	services.GET.IsTakingPicture = true

	logs.Info("Taking picture...")
	go func() {
		timeout := config.GET.Photobooth.DefaultTimer

		for timeout >= 0 {
			s.Send("TIMER", timeout)
			timeout--
			time.Sleep(1 * time.Second)
		}

		services.GET.IsTakingPicture = false
	}()
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

func (s *Socket) SendState() error {
	settings := services.GET.GetFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	return s.Send("APP_STATE", settings)
}

func (s *Socket) StartPings() {
	go func() {
		for s.Open {
			time.Sleep(1 * time.Second)
			s.Send("PING", time.Now().Format("2006-01-02 15:04:05"))
		}
	}()
}

func (s *Socket) StartUnattended() {
	if s.Type != SOCKET_TYPE_BOOTH || config.GET.Photobooth.UnattendedInterval < 1 {
		return
	}

	go func() {
		for s.Open {
			time.Sleep(time.Duration(config.GET.Photobooth.UnattendedInterval) * time.Minute)
			logs.Info("Unattended picture")
			s.Send("UNATTENDED_PICTURE", nil)
		}
	}()
}
