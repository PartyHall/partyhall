package message_handler

import (
	"fmt"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
)

type MqttSenderHandler struct {
	Type       string
	Topic      string
	GetPayload func(interface{}) interface{}
}

func (h MqttSenderHandler) GetType() string {
	return h.Type
}

func (h MqttSenderHandler) Do(s *socket.Socket, payload interface{}) {
	token := (*services.GET.MqttClient).Publish(
		h.Topic,
		2,
		false,
		h.GetPayload(payload),
	)

	token.Wait()
	if token.Error() != nil {
		logs.Error(token.Error())
		// @TODO: Tell the sender
	}
}

func payloadedString(payload interface{}) interface{} {
	return fmt.Sprintf("%v", payload)
}

func getMqttHandler(hType string, topic string, getPayload func(interface{}) interface{}) MqttSenderHandler {
	return MqttSenderHandler{
		Type:       hType,
		Topic:      config.GetMqttTopic(topic),
		GetPayload: getPayload,
	}
}

func getButtonPressHandler(hType string, button string) MqttSenderHandler {
	return getMqttHandler(
		hType,
		"button_press",
		func(payload interface{}) interface{} {
			return fmt.Sprintf("%v", button)
		},
	)
}
