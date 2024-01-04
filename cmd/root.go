package cmd

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/modules"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/routes"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "partyhall",
	Short: "The partyhall main app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}

		modules.BroadcastFrontendSettings()

		err := remote.InitMqtt()
		if err != nil {
			logs.Error(err)
			os.Exit(1)
		}

		modules.PreInitializeModules()
		remote.Initialize()

		for name, module := range modules.MODULES {
			// Must be not be done on the same map
			handlers := module.GetMqttHandlers()
			for k, v := range module.GetMqttHandlers() {
				handlers[name+"/"+k] = v
				delete(handlers, k)
			}

			remote.EasyMqtt.RegisterHandlers(handlers)
			remote.EasyWS.RegisterMessageHandlers(module.GetWebsocketHandlers()...)
		}

		modules.InitializeModules()

		_, err = orm.GET.AppState.GetState()
		if err != nil {
			if err != sql.ErrNoRows {
				logs.Error("Failed to load appstate: ", err)
				os.Exit(1)
			}

			as := models.AppState{
				HardwareID: uuid.New().String(),
				ApiToken:   nil, // @TODO: The token should be retreived from the API server while setting the partyhall up
			}
			err := orm.GET.AppState.CreateState(as)
			if err != nil {
				logs.Error("Failed to save the state: ", err)
				os.Exit(1)
			}

			logs.Info("Initializing the partyhall with id ", as.HardwareID)
		}

		err = orm.GET.Events.ClearExporting()
		if err != nil {
			logs.Error("Failed to clear event exporting.")
			logs.Error("Some event might be in a wrong state...")
			logs.Error(err)
		}

		r := mux.NewRouter()

		r.PathPrefix("/media/partyhall").Handler(http.StripPrefix("/media/partyhall", http.FileServer(http.Dir(utils.GetPath("images")))))
		r.PathPrefix("/media/karaoke").Handler(http.StripPrefix("/media/karaoke", http.FileServer(http.Dir(utils.GetPath("karaoke")))))

		routes.Register(r.PathPrefix("/api").Subrouter())

		if services.WEBAPP_FS != nil {
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				mode := "booth"
				if utils.IsRemote(r) {
					mode = "admin"
				}

				w.Write([]byte(services.InjectHtmlMode(mode)))
			})

			r.PathPrefix("/").Handler(http.FileServer(http.FS(*services.WEBAPP_FS)))
		} else {
			logs.Error("Failed to embed webapp: not loaded")
		}

		logs.Infof("PartyHall app is listening on %v\n", config.GET.Web.ListeningAddr)
		err = http.ListenAndServe(config.GET.Web.ListeningAddr, r)
		if err != nil {
			logs.Error("Failed to listen on the given address/port", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(hwhandlerCmd)
	rootCmd.AddCommand(versionCmd)
}
