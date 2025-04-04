package mqtt

import (
	"encoding/json"

	emqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/easymqtt"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/state"
)

const (
	ACTION_TAKE_PICTURE  = "take_picture"
	ACTION_LEFT          = "left"
	ACTION_RIGHT         = "right"
	ACTION_SHUTDOWN      = "shutdown"
	ACTION_DISPLAY_DEBUG = "display_debug"
)

var BUTTON_ACTIONS = []string{
	ACTION_TAKE_PICTURE,
	ACTION_DISPLAY_DEBUG,
	ACTION_SHUTDOWN,
	ACTION_LEFT,
	ACTION_RIGHT,
}

var EasyMqtt *easymqtt.EasyMqtt

func Load() error {
	handlers := map[string]emqtt.MessageHandler{
		"button_press":  OnButtonPress,
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

func OnButtonPress(client emqtt.Client, msg emqtt.Message) {
	var btnPressedMessage struct {
		Button int `json:"button"`
	}

	err := json.Unmarshal(msg.Payload(), &btnPressedMessage)
	if err != nil {
		log.Error("[MQTT] Error parsing payload", "err", err)
		return
	}

	if state.STATE.CurrentMode == state.MODE_BTN_SETUP {
		mercure_client.CLIENT.SendButtonPress(btnPressedMessage.Button)

		return
	}

	action, ok := state.STATE.UserSettings.ButtonMappings[btnPressedMessage.Button]
	if !ok {
		log.Error("Pressed button that is not binded to anything", "btn", btnPressedMessage.Button)
	}

	switch action {
	case ACTION_TAKE_PICTURE:
		OnTakePicture(client, msg)
	case ACTION_SHUTDOWN:
		OnShutdown(client, msg)
	case ACTION_DISPLAY_DEBUG:
		OnDisplayDebug(client, msg)
	case ACTION_RIGHT:
		OnNextBackdrop(client, msg)
	default:
		log.Error("No mqtt action for " + action)
	}
}
