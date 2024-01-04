package remote

import (
	"fmt"

	"github.com/partyhall/easyws"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
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
	// @TODO: Do the fmt.Sprintf in EasyMqtt
	err := EasyMqtt.Send(h.Topic, fmt.Sprintf("%v", h.GetPayload(payload)))
	if err != nil {
		logs.Errorf("Failed to send mqtt message on topic %v: %v", h.Topic, err)
	}
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
