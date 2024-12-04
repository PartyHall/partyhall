package mqtt

import (
	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/easymqtt"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
)

var EasyMqtt *easymqtt.EasyMqtt

func Load() error {
	handlers := map[string]emqtt.MessageHandler{
		"display_debug": OnDisplayDebug,
		"shutdown":      OnShutdown,
		"take_picture":  OnTakePicture,
	}

	EasyMqtt = easymqtt.New(
		"partyhall",
		config.GET.MosquittoAddr,
		"partyhall",
		handlers,
		func(c emqtt.Client, err error) {
			log.Error("[MQTT] Connection lost...")
		},
		func(c emqtt.Client, co *emqtt.ClientOptions) {
			log.Info("[MQTT] Reconnecting...")
		},
	)

	return EasyMqtt.Start()
}
