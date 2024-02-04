package remote

import (
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easymqtt"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"

	"github.com/partyhall/easyws"
)

var EasyWS easyws.EasyWS
var EasyMqtt *easymqtt.EasyMqtt

var onJoins map[string]func(socketType string, s *easyws.Socket)

func RegisterOnJoin(module string, function func(socketType string, s *easyws.Socket)) {
	if onJoins == nil {
		onJoins = make(map[string]func(socketType string, s *easyws.Socket))
	}

	onJoins[module] = function
}

func Initialize() {
	EasyWS = easyws.NewWithTypes(
		utils.SOCKET_TYPES,
		func(socketType string, c *echo.Context) bool {
			// @deprecated: use middleware
			return true
		},
	)

	EasyWS.RegisterMessageHandlers(
		PongHandler{},
		GetStateHandler{},
		SetModeHandler{},
		SetDateTimeHandler{},
		SetEventHandler{},
		asMqtt("DISPLAY_DEBUG"),
		asMqtt("EXPORT_ZIP"),
		asMqtt("SHUTDOWN"),
	)

	EasyWS.OnJoin = func(socketType string, s *easyws.Socket) {
		SendState(s)
		for _, function := range onJoins {
			function(socketType, s)
		}
	}
}

func InitMqtt() error {
	handlers := map[string]mqtt.MessageHandler{
		"display_debug":  OnDisplayDebug,
		"export_zip":     OnExportZip,
		"shutdown":       OnShutdown,
		"admin/set_mode": OnSetMode,
	}

	EasyMqtt = easymqtt.New(
		"partyhall",
		config.GET.Mosquitto.Address,
		"partyhall",
		handlers,
		func(c mqtt.Client, err error) {
			logs.Error("[MQTT] Connection lost...")
		},
		func(c mqtt.Client, co *mqtt.ClientOptions) {
			logs.Info("[MQTT] Reconnecting...")
		},
	)

	return EasyMqtt.Start()
}

func SendState(s *easyws.Socket) error {
	settings := services.BuildFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	return s.Send("APP_STATE", settings)
}

func BroadcastBooth(msgType string, payload interface{}) {
	EasyWS.BroadcastTo(utils.SOCKET_TYPE_BOOTH, msgType, payload)
}

func BroadcastAdmin(msgType string, payload interface{}) {
	EasyWS.BroadcastTo(utils.SOCKET_TYPE_ADMIN, msgType, payload)
}

func Broadcast(msgType string, payload interface{}) {
	EasyWS.BroadcastTo("", msgType, payload)
}

func BroadcastState() error {
	settings := services.BuildFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	EasyWS.BroadcastTo("", "APP_STATE", settings)

	return nil
}
