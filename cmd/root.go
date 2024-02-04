package cmd

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dto"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/modules"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/routes"
	"github.com/partyhall/partyhall/services"
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
			logs.Error("Failed to get AppState")
			return
		}

		err = orm.GET.Events.ClearExporting()
		if err != nil {
			logs.Error("Failed to clear event exporting.")
			logs.Error("Some event might be in a wrong state...")
			logs.Error(err)
		}

		e := echo.New()
		e.Validator = dto.NewValidator()
		e.Use(middleware.Recover())

		routes.Register(e.Group("/api"))

		log.Info("Registered routes: ")
		for _, r := range e.Routes() {
			log.Info("\t- " + r.Path)
		}

		logs.Infof("PartyHall app is listening on %v\n", config.GET.Web.ListeningAddr)
		logs.Error(e.Start(config.GET.Web.ListeningAddr))

		/*
			r.PathPrefix("/media/partyhall").Handler(http.StripPrefix("/media/partyhall", http.FileServer(http.Dir(utils.GetPath("images")))))

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
		*/
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(hwhandlerCmd)
	rootCmd.AddCommand(convertVideoCmd)
	rootCmd.AddCommand(spotifySearchCmd)
	rootCmd.AddCommand(versionCmd)
}
