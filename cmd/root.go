package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
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

		_, err := orm.GET.AppState.GetState()
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

		routes.Register(r.PathPrefix("/api").Subrouter())

		if services.WEBAPP_FS != nil {
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				mode := "booth"
				if utils.IsRemote(r) {
					mode = "admin"
				}

				// @TODO: Load the index.html in the embedfs + parse it + inject my script in the correct spot
				// Currently this does not work as vite generates the filepath for the js/css files
				script := fmt.Sprintf(`<script>window.SOCKET_TYPE = '%v';</script>`, mode)

				w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<link rel="icon" type="image/svg+xml" href="/partyhall.svg" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>PartyHall</title>
		%v
		<script type="module" crossorigin src="/assets/index.10ad5d49.js"></script>
		<link rel="stylesheet" href="/assets/index.6797710c.css">
	</head>

	<body>
		<div id="root"></div>
	</body>
</html>
				`, script)))
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
