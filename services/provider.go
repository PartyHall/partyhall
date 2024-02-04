package services

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/partyhall/easymqtt"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/migrations"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
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

	EchoJWTPrivateKey   []byte
	EchoJWTPublicKey    []byte
	EchoJwtConfig       echojwt.Config
	EchoWsJwtConfig     echojwt.Config
	EchoJwtMiddleware   echo.MiddlewareFunc
	EchoWsJwtMiddleware echo.MiddlewareFunc

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

	private, public, err := utils.GetOrGenerateJwtKeys()
	if err != nil {
		return err
	}

	prv := &Provider{
		CurrentMode:    config.GET.DefaultMode,
		Spotify:        Spotify{},
		ModuleSettings: map[string]interface{}{},
		EchoJwtConfig: echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(models.JwtCustomClaims)
			},
			SigningKey: private,
		},
		EchoWsJwtConfig: echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(models.JwtCustomClaims)
			},
			SigningKey:  private,
			TokenLookup: "query:token:",
			Skipper: func(c echo.Context) bool {
				if strings.ToUpper(c.Param("type")) == utils.SOCKET_TYPE_BOOTH {
					if utils.IsRemote(c) {
						if config.GET.DebugMode || config.IsInDev() {
							logs.Debug("Letting a remote connection")
							return true
						}

						return false
					}

					return true
				}

				return false
			},
		},
		EchoJWTPrivateKey: private,
		EchoJWTPublicKey:  public,
	}

	prv.EchoJwtMiddleware = echojwt.WithConfig(prv.EchoJwtConfig)
	prv.EchoWsJwtMiddleware = echojwt.WithConfig(prv.EchoWsJwtConfig)

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
