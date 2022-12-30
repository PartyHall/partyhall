package mqtt_handler

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
)

type MessageHandler interface {
	GetTopic() string
	Do(client mqtt.Client, msg mqtt.Message)
}

var MESSAGE_HANDLERS map[string]MessageHandler = make(map[string]MessageHandler)

func init() {
	handlers := []MessageHandler{
		ButtonPressHandler{},
		SyncRequestedHandler{},
		ExportHandler{},
	}

	for _, h := range handlers {
		MESSAGE_HANDLERS[h.GetTopic()] = h
	}
}

func getClientOptions() *mqtt.ClientOptions {
	opts := mqtt.
		NewClientOptions().
		AddBroker(config.GET.Mosquitto.Address).
		SetClientID("partyhall").
		SetPingTimeout(10 * time.Second).
		SetKeepAlive(10 * time.Second)

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		logs.Errorf("[MQTT] Connection lost: %s\n" + err.Error())
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		logs.Info("[MQTT] Reconnecting...")
	})

	return opts
}

func ConnectAndListen() error {
	client := mqtt.NewClient(getClientOptions())
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for topic, action := range MESSAGE_HANDLERS {
		token := client.Subscribe(config.GetMqttTopic(topic), 2, action.Do)
		waited := token.Wait()

		if !waited {
			return errors.New("failed to wait on subscription for topic " + topic)
		}

		err := token.Error()
		if err != nil {
			return err
		}
	}

	services.GET.MqttClient = &client

	return nil
}
