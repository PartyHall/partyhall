package message_handler

import (
	"os/exec"
	"time"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/socket"
)

type SetDateTimeHandler struct{}

func (h SetDateTimeHandler) GetType() string {
	return "SET_DATETIME"
}

func (h SetDateTimeHandler) Do(s *socket.Socket, payload interface{}) {
	dt, ok := payload.(string)
	if !ok {
		s.Send("ERR_MODAL", "Bad request")
		return
	}

	time, err := time.Parse("2006-01-02 15:04:05", dt)
	if err != nil {
		s.Send("ERR_MODAL", "Failed to set date: "+err.Error())
		return
	}

	_, lookErr := exec.LookPath("sudo")
	if lookErr != nil {
		logs.Errorf("Sudo binary not found, cannot set system date: %s\n", lookErr.Error())
		s.Send("ERR_MODAL", "Sudo binary not found")
		return
	}

	dateString := time.Format("2 Jan 2006 15:04:05")

	logs.Errorf("Setting system date to: %s\n", dateString)
	args := []string{"date", "--set", dateString}
	err = exec.Command("sudo", args...).Run()

	if err != nil {
		logs.Errorf("Failed to set the date: ", err)
		s.Send("ERR_MODAL", "Failed to set the date: "+err.Error())
	}
}
