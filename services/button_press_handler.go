package services

import (
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/logs"
)

var BPH *ButtonPressHandler = nil

type ButtonPressHandler struct {
	handlers map[string]func(client mqtt.Client)
	pb       *PartyHall
}

func initButtonHandler(pb *PartyHall) {
	BPH = &ButtonPressHandler{}
	BPH.pb = pb

	BPH.handlers = map[string]func(client mqtt.Client){
		"TAKE_PICTURE":  BPH.onTakePicture,
		"DISPLAY_DEBUG": BPH.onDisplayDebug,
		"SHUTDOWN":      BPH.onShutdown,
	}
}

func (bph *ButtonPressHandler) OnButtonPress(client mqtt.Client, msg mqtt.Message) {
	handler, ok := bph.handlers[string(msg.Payload())]
	if ok {
		handler(client)
	} else {
		logs.Error("Unknown button pressent: ", string(msg.Payload()))
	}

	msg.Ack()
}

func (bph *ButtonPressHandler) onTakePicture(client mqtt.Client) {
	bph.pb.prv.Sockets.BroadcastTakePicture()
}

func (bph *ButtonPressHandler) onDisplayDebug(client mqtt.Client) {
	if bph.pb.DisplayDebug {
		return
	}

	bph.pb.DisplayDebug = true
	GET.Sockets.BroadcastState()
	go func() {
		time.Sleep(30 * time.Second)
		bph.pb.DisplayDebug = false
		GET.Sockets.BroadcastState()
	}()
}

func (bph *ButtonPressHandler) onShutdown(client mqtt.Client) {
	if err := GET.Shutdown(); err != nil {
		GET.Sockets.broadcastTo("", "ERR_MODAL", "Failed to shutdown: "+err.Error())
		return
	}

	exec.Command("shutdown", "-h", "now").Run()
}
