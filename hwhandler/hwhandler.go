package hwhandler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"go.bug.st/serial"
)

var client mqtt.Client

type HardwareHandler struct {
	PortName string
}

func Load() error {
	opts := mqtt.NewClientOptions().AddBroker(config.GET.Mosquitto.Address).SetClientID("hwhandler").SetPingTimeout(10 * time.Second).SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		fmt.Printf("[MQTT] Connection lost: %s\n" + err.Error())
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Println("[MQTT] Reconnecting...")
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (hh *HardwareHandler) ProcessSerialConn() {
	port, err := serial.Open(hh.PortName, config.GET.HardwareHandler.SerialMode)
	if err != nil {
		fmt.Println("Failed to open device: ")
		fmt.Println(err)
		return
	}

	for {
		line, err := hh.readLine(port)
		if err != nil {
			fmt.Println("Failed to read: ", err)
			continue
		}

		hh.processSerialMessage(line)
	}
}

func (hh *HardwareHandler) readLine(port serial.Port) (string, error) {
	str := ""

	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			return str, err
		}

		if n == 0 {
			fmt.Println("No clue what to do here")
			return "", errors.New("byte = 0")
		}

		str += string(buff[:n])
		if strings.HasSuffix(str, "\n") || strings.HasSuffix(str, "\r") {
			break
		}
	}

	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")

	return str, nil
}

func (hh *HardwareHandler) processSerialMessage(msg string) {
	if strings.HasPrefix(msg, "BTN_") {
		msg = strings.Trim(msg, " \t")
		val, ok := config.GET.HardwareHandler.Mappings[msg]
		if !ok {
			fmt.Println("Unknown button: " + msg)
			return
		}

		client.Publish(config.GetMqttTopic("button_press"), 2, false, val)

		return
	}

	data := strings.Split(msg, " ")
	if len(data) == 0 {
		fmt.Println("Bad message received")
		return
	}

	switch data[0] {
	default:
		fmt.Println("Unhandled arduino message: ", data[0])
		fmt.Println(strings.Join(append([]string{"\targs: "}, data[1:]...), " "))
	}
}
