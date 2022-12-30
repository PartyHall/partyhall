package services

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"

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
	DB_SCRIPTS_FS embed.FS
)

var GET *Provider

type Provider struct {
	MqttClient *mqtt.Client

	CurrentState    models.AppState
	IsTakingPicture bool
	CurrentMode     string
}

func (p *Provider) GetFrontendSettings() *models.FrontendSettings {
	settings := models.FrontendSettings{
		AppState:    p.CurrentState,
		Photobooth:  config.GET.Photobooth,
		CurrentMode: p.CurrentMode,

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

	if GET.CurrentState.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*GET.CurrentState.CurrentEvent)
		if err == nil {
			GET.CurrentState.CurrentEventObj = evt
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

	prv.CurrentState = state

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

	prv := &Provider{
		CurrentMode: config.GET.DefaultMode,
	}

	err = prv.loadState()
	if err != nil {
		return err
	}

	GET = prv

	return nil
}

func (p *Provider) Shutdown() error {
	return orm.GET.DB.Close()
}
