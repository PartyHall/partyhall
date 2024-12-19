package hwhandler

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"go.bug.st/serial"
)

type Device struct {
	Handler             *HardwareHandler
	PortName            string
	Port                *serial.Port
	LastPing            time.Time
	LastButtonPress     string
	LastButtonPressTime time.Time
}

type HardwareHandler struct {
	Devices []*Device
	Mqtt    *mqtt.Client
}

func Load() (*HardwareHandler, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}

	if len(ports) == 0 {
		return nil, errors.New("no arduino/esp32 found")
	}

	devices := []*Device{}
	for _, p := range ports {
		if !strings.HasPrefix(p, "/dev/ttyUSB") {
			continue
		}

		devices = append(devices, &Device{
			PortName:            p,
			LastPing:            time.Now(),
			LastButtonPressTime: time.Now(),
			LastButtonPress:     "",
		})
	}

	return &HardwareHandler{
		Devices: devices,
	}, nil
}

func (hh *HardwareHandler) Start() error {
	opts := mqtt.
		NewClientOptions().
		AddBroker(config.GET.MosquittoAddr).
		SetClientID("hwhandler").
		SetPingTimeout(10 * time.Second).
		SetKeepAlive(10 * time.Second).
		SetAutoReconnect(true).
		SetMaxReconnectInterval(10 * time.Second).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			fmt.Printf("[MQTT] Connection lost: %s\n", err.Error())
		}).
		SetReconnectingHandler(func(c mqtt.Client, co *mqtt.ClientOptions) {
			fmt.Println("[MQTT] Reconnecting...")
		})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	hh.Mqtt = &client

	for _, d := range hh.Devices {
		d.Handler = hh

		err := d.subscribe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = hh.ProcessSerialConn(d)
		if err != nil {
			fmt.Printf("Failed to init device %v: %v\n", d.PortName, err)
			os.Exit(1)
		}
	}

	return nil
}

func (d *Device) subscribeToMqtt(topic string, fn func(c mqtt.Client, m mqtt.Message)) error {
	token := (*d.Handler.Mqtt).Subscribe(topic, 1, fn)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %v", topic, err)
	}

	fmt.Println("Subscribed to " + topic)

	return nil
}

func (d *Device) subscribe() error {
	if err := d.subscribeToMqtt("partyhall/flash", d.OnFlash); err != nil {
		return err
	}

	return nil
}

func (hh *HardwareHandler) ProcessSerialConn(d *Device) error {
	port, err := serial.Open(d.PortName, &serial.Mode{
		BaudRate: config.GET.HardwareHandler.BaudRate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})

	if err != nil {
		return err
	}

	d.Port = &port
	d.LastPing = time.Now()
	d.HandlePing()

	go func() {
		for {
			err := d.ProcessMessage()

			if err != nil {
				fmt.Println("Failed to read: ", err)
				continue
			}

		}
	}()

	return nil
}

func (d *Device) HandlePing() {
	go func() {
		for {
			time.Sleep(time.Second)
			if d.LastPing.Add(15 * time.Second).Before(time.Now()) {
				fmt.Println("No ping received for the last 15 seconds. Crashing !")
				os.Exit(1)
			}
		}
	}()
}

func (d *Device) OnFlash(c mqtt.Client, m mqtt.Message) {
	data := strings.ToLower(string(m.Payload()))
	fmt.Println("[MQTT] Flash " + data)

	var err error = nil
	if data == "on" {
		err = d.WriteMessage("FLASH_ON")
	} else if data == "off" {
		err = d.WriteMessage("FLASH_OFF")
	} else {
		fmt.Println("\t=> Bad request")
	}

	if err != nil {
		fmt.Println("\t=> Failed to send message to device: " + err.Error())
	}
}

func (d *Device) WriteMessage(msg string) error {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	_, err := (*d.Port).Write([]byte(msg))

	return err
}

func (d *Device) ProcessMessage() error {
	msg, err := readLine(*d.Port)
	if err != nil {
		return err
	}

	if len(msg) == 0 {
		return nil
	}

	// If its a button press, special case
	if strings.HasPrefix(msg, "BTN_") {
		msg = strings.Trim(msg, " \t")
		val, ok := config.GET.HardwareHandler.Mappings[msg]
		if !ok {
			debugMsg(msg)

			return nil
		}

		currTime := time.Now()
		diff := currTime.Sub(d.LastButtonPressTime).Seconds()

		// Debounce
		if d.LastButtonPress != val || diff > 1 {
			topic := "partyhall/" + strings.ToLower(val)
			fmt.Println("Button pressed: ", msg, " sending ", topic)
			(*d.Handler.Mqtt).Publish(topic, 2, false, "press")

			d.LastButtonPress = val
			d.LastButtonPressTime = currTime
		}

		return nil
	}

	data := strings.Split(msg, " ")
	switch data[0] {
	case "STARTING_UP":
		fmt.Println("Arduino's starting up...")
	case "PING":
		d.LastPing = time.Now()
	default:
		debugMsg(msg)
	}

	return nil
}
