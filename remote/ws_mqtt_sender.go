package remote

import (
	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
)

type MqttSenderHandler struct {
	Type       string
	Topic      string
	GetPayload func(interface{}) interface{}
}

func (h MqttSenderHandler) GetType() string {
	return h.Type
}

func (h MqttSenderHandler) Do(s *easyws.Socket, payload interface{}) {
	EasyMqtt.Send(h.Topic, h.GetPayload(payload))
}

func asMqtt(message string) MqttSenderHandler {
	return MqttSenderHandler{
		Type:  message,
		Topic: config.GetMqttTopic("", message),
		GetPayload: func(i interface{}) interface{} {
			return i
		},
	}
}
