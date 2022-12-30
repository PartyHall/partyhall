package message_handler

import (
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/socket"
)

type MessageHandler interface {
	GetType() string
	Do(s *socket.Socket, payload interface{})
}

var MESSAGE_HANDLERS map[string]MessageHandler = make(map[string]MessageHandler)

func init() {
	handlers := []MessageHandler{
		PongHandler{},
		GetStateHandler{},
		TakePictureHandler{},
		RemoteTakePictureHandler{},
		SetModeHandler{},
		SetDateTimeHandler{},
		SetEventHandler{},
		getButtonPressHandler("DISPLAY_DEBUG", "DISPLAY_DEBUG"),
		getMqttHandler("EXPORT_ZIP", "export", payloadedString),
		getButtonPressHandler("SHUTDOWN", "SHUTDOWN"),
	}

	for _, h := range handlers {
		MESSAGE_HANDLERS[h.GetType()] = h
	}
}

func ProcessMessage(socket *socket.Socket, msg models.SocketMessage) {
	for name, handler := range MESSAGE_HANDLERS {
		if name == msg.MsgType {
			handler.Do(socket, msg.Payload)
			return
		}
	}

	logs.Infof("Unhandled socket message: %v", msg.MsgType)
}
