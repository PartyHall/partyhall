package module_photobooth

import (
	"time"

	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/utils"
)

type Action struct{}

func (a Action) StartUnattended(s *easyws.Socket) {
	if s.Type != utils.SOCKET_TYPE_BOOTH || CONFIG.UnattendedInterval < 1 {
		return
	}

	go func() {
		for s.Open {
			time.Sleep(time.Duration(CONFIG.UnattendedInterval) * time.Minute)
			logs.Info("Unattended picture")
			s.Send("UNATTENDED_PICTURE", nil)
		}
	}()
}

func (a Action) TakePicture(s *easyws.Socket) {
	if s.Type != utils.SOCKET_TYPE_BOOTH {
		return
	}

	logs.Info("Taking picture...")
	s.Send("TAKE_PICTURE", nil)
}
