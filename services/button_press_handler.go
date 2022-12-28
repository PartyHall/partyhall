package services

import (
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
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
	logs.Debug("HW Display debug request")
	GET.Sockets.BroadcastState()
	GET.Sockets.BroadcastBooth("DISPLAY_DEBUG", nil)
}

func (bph *ButtonPressHandler) onShutdown(client mqtt.Client) {
	if config.IsInDev() {
		logs.Info("[IN DEV] Shutdown requested")
		return
	}

	if err := GET.Shutdown(); err != nil {
		GET.Sockets.broadcastTo("", "ERR_MODAL", "Failed to shutdown: "+err.Error())
		return
	}

	exec.Command("shutdown", "-h", "now").Run()
}
