package hwhandler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"go.bug.st/serial"
)

var client mqtt.Client

var lastButtonPress string = ""
var lastButtonPressTime time.Time = time.Now()

// @TODO:
// PINGPONG with the arduino & autorestart when its not working

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

	// We read char by char
	buff := make([]byte, 1)
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

func (hh *HardwareHandler) DebugMsg(msg string, hx string) {
	fmt.Printf("Unknown message: %v\n\t=> 0x%v\n", msg, hx)
}

func (hh *HardwareHandler) processSerialMessage(msg string) {
	if len(msg) == 0 {
		return
	}

	hx := hex.EncodeToString([]byte(msg))

	if strings.HasPrefix(msg, "BTN_") {
		msg = strings.Trim(msg, " \t")
		val, ok := config.GET.HardwareHandler.Mappings[msg]
		if !ok {
			hh.DebugMsg(msg, hx)
			return
		}

		currTime := time.Now()
		diff := currTime.Sub(lastButtonPressTime).Seconds()

		// Debounce
		if lastButtonPress != val || diff > 1 {
			topic := config.GetMqttTopic("", strings.ToLower(val))
			fmt.Println("Button pressed: ", msg, " sending ", topic)
			client.Publish(topic, 2, false, "press")

			lastButtonPress = val
			lastButtonPressTime = currTime
		}

		return
	}

	if len(msg) == 0 {
		fmt.Print("Bad message: ")
		hh.DebugMsg(msg, hx)
		return
	}

	data := strings.Split(msg, " ")

	switch data[0] {
	case "STARTING_UP":
		fmt.Println("Arduino's starting up...")
	case "OK_RF24":
		fmt.Println("Wireless device detected & ready to be used")
	default:
		hh.DebugMsg(msg, hx)
		fmt.Println(strings.Join(append([]string{"\targs: "}, data[1:]...), " "))
	}
}
