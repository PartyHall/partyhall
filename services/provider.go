package services

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/migrations"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
)

var (
	WEBAPP_FS     *fs.FS
	ADMIN_FS      *fs.FS
	DB_SCRIPTS_FS embed.FS
)

var GET *Provider

type Provider struct {
	Sockets    Sockets
	MqttClient *mqtt.Client

	Admin     Admin
	PartyHall PartyHall
}

func (p *Provider) GetFrontendSettings() *models.FrontendSettings {
	settings := models.FrontendSettings{
		AppState:     p.PartyHall.CurrentState,
		Photobooth:   config.GET.PartyHall,
		DebugDisplay: p.PartyHall.DisplayDebug,
		CurrentMode:  p.PartyHall.CurrentMode,

		IPAddress:  map[string][]string{},
		KnownModes: config.MODES,

		PartyHallVersion: utils.CURRENT_VERSION,
		PartyHallCommit:  utils.CURRENT_COMMIT,
	}

	events, err := orm.GET.Events.GetEvents()
	if err != nil {
		logs.Error("Failed to get events: ", err)
		return nil
	}

	if GET.PartyHall.CurrentState.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*GET.PartyHall.CurrentState.CurrentEvent)
		if err == nil {
			GET.PartyHall.CurrentState.CurrentEventObj = evt
			settings.AppState.CurrentEventObj = evt
		}
	}

	settings.KnownEvents = events

	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		shouldSkip := false
		for _, ignored := range []string{"lo", "br-", "docker", "vmnet", "veth"} { // Ignoring docker / vmware networks for in-dev purposes
			if strings.HasPrefix(inter.Name, ignored) {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			continue
		}

		settings.IPAddress[inter.Name] = []string{}

		addrs, _ := inter.Addrs()
		for _, ip := range addrs {
			settings.IPAddress[inter.Name] = append(settings.IPAddress[inter.Name], ip.String())
		}
	}

	return &settings
}

func SetSystemDate(newTime time.Time) error {
	_, lookErr := exec.LookPath("sudo")
	if lookErr != nil {
		logs.Errorf("Sudo binary not found, cannot set system date: %s\n", lookErr.Error())
		return lookErr
	} else {
		dateString := newTime.Format("2 Jan 2006 15:04:05")
		logs.Errorf("Setting system date to: %s\n", dateString)
		args := []string{"date", "--set", dateString}
		return exec.Command("sudo", args...).Run()
	}
}

func (prv *Provider) loadState() error {
	state, err := orm.GET.AppState.GetState()
	if err != nil {
		return err
	}

	if state.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*state.CurrentEvent)
		if err != nil {
			return err
		}

		state.CurrentEventObj = evt
	}

	prv.PartyHall.CurrentState = state

	return nil
}

func Load() error {
	err := utils.MakeOrCreateFolder("")
	if err != nil {
		fmt.Println("Failed to create root folder")
		os.Exit(1)
	}

	if err := migrations.CheckDbExists(DB_SCRIPTS_FS); err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	for _, folder := range []string{"images"} {
		if err := utils.MakeOrCreateFolder(folder); err != nil {
			return err
		}
	}

	err = orm.Load()
	if err != nil {
		return err
	}

	opts := mqtt.NewClientOptions().AddBroker(config.GET.Mosquitto.Address).SetClientID("partyhall").SetPingTimeout(10 * time.Second).SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		logs.Errorf("[MQTT] Connection lost: %s\n" + err.Error())
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		logs.Info("[MQTT] Reconnecting...")
	})

	prv := &Provider{
		Sockets: []*Socket{},
	}

	prv.Admin = Admin{prv: prv}
	prv.PartyHall = PartyHall{
		prv:         prv,
		CurrentMode: config.GET.DefaultMode,
	}

	err = prv.loadState()
	if err != nil {
		return err
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	initButtonHandler(&prv.PartyHall)

	actions := map[string]mqtt.MessageHandler{
		"partyhall/button_press":   BPH.OnButtonPress,
		"partyhall/sync":           prv.PartyHall.OnSyncRequested,
		"partyhall/export":         prv.PartyHall.OnExportEvent,
		"partyhall/admin/set_mode": prv.Admin.OnSetMode,
	}

	for topic, action := range actions {
		if token := client.Subscribe(topic, 2, action); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	prv.MqttClient = &client
	GET = prv

	return nil
}

func (p *Provider) Shutdown() error {
	return orm.GET.DB.Close()
}
