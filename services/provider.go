package services

import (
	"embed"
	"io/fs"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/partyhall/easymqtt"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/migrations"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

var (
	WEBAPP_FS              *fs.FS
	DB_SCRIPTS_FS          embed.FS
	KARAOKE_FALLBACK_IMAGE []byte
)

var GET *Provider

type Provider struct {
	MqttClient *mqtt.Client
	EasyMqtt   *easymqtt.EasyMqtt

	CurrentState   models.AppState
	CurrentMode    string
	Spotify        Spotify
	ModuleSettings map[string]interface{}
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
	err := migrations.CheckDbExists(DB_SCRIPTS_FS)
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}

	err = orm.Load()
	if err != nil {
		return err
	}

	prv := &Provider{
		CurrentMode:    config.GET.DefaultMode,
		Spotify:        Spotify{},
		ModuleSettings: map[string]interface{}{},
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
