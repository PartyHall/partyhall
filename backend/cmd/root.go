package cmd

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/dunglas/mercure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/cron"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/mercure_client"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/mqtt"
	"github.com/partyhall/partyhall/nexus"
	"github.com/partyhall/partyhall/routes"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/state"
	"github.com/partyhall/partyhall/utils"
	"github.com/partyhall/partyhall/validators"
	"github.com/spf13/cobra"
)

var AppFS embed.FS
var ApplianceFS embed.FS

func closeDb() {
	if services.DB != nil {
		services.DB.Close()
	}
}

func isAppliance(ctx *gin.Context) bool {
	clientIP := ctx.ClientIP()

	return clientIP == "127.0.0.1" || clientIP == "::1"
}

var rootCmd = &cobra.Command{
	Use:   "partyhall",
	Short: "The partyhall main app",
	Run: func(cmd *cobra.Command, args []string) {
		defer log.LOG.Sync()
		defer closeDb()

		err := services.Load()
		if err != nil {
			log.LOG.Error(err)
			return
		}

		err = mqtt.Load()
		if err != nil {
			log.LOG.Error(err)
			return
		}

		_, err = mercure_client.NewClient()
		if err != nil {
			log.LOG.Error(err)
			return
		}

		state.STATE.CurrentMode = state.MODE_PHOTOBOOTH

		dal.DB = services.DB

		event, err := dal.EVENTS.GetCurrent()
		if err != nil {
			log.LOG.Error(err)
			return
		}

		if event == nil {
			evt, err := dal.EVENTS.GetAndSetAny()
			if err != nil {
				log.LOG.Error(err)
				return
			}

			if evt == nil {
				state.STATE.CurrentMode = state.MODE_DISABLED
			} else {
				event = evt
			}
		}

		if len(utils.CURRENT_COMMIT) > 7 {
			utils.CURRENT_COMMIT = utils.CURRENT_COMMIT[:7]
		}

		state.STATE.GuestsAllowed = config.GET.GuestsAllowed
		state.STATE.HardwareId = config.GET.HardwareID
		state.STATE.Version = utils.CURRENT_VERSION
		state.STATE.Commit = utils.CURRENT_COMMIT

		dal.SONGS.WipeInvalidSessions()

		state.STATE.Karaoke = state.KaraokeState{
			Current:   nil,
			IsPlaying: false,
			Timecode:  0,
			Volume:    100,
		}

		state.STATE.KaraokeQueue = []*models.SongSession{}
		services.SetEvent(event)

		nexus.NewClient(
			config.GET.NexusURL,
			config.GET.HardwareID,
			config.GET.ApiKey,
			config.GET.NexusIgnoreSSL,
		)

		if config.GET.IsInDev {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		h, err := mercure.NewHub(
			mercure.WithPublisherJWT([]byte(config.GET.Mercure.PublisherKey), "HS256"),
			mercure.WithSubscriberJWT([]byte(config.GET.Mercure.SubscriberKey), "HS256"),
			mercure.WithWriteTimeout(0),
		)

		if err != nil {
			log.LOG.Fatalw("Failed to create Mercure hub", "err", err)
			os.Exit(1)
		}
		defer h.Stop()

		// Setting up the routes
		r := gin.Default()

		// r.Use(log.GinLogger())

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			// Register the custom validation function
			v.RegisterValidation("iso8601", validators.IsIso8601)
		}

		r.Any("/.well-known/mercure", gin.WrapH(h))

		r.GET("/", func(ctx *gin.Context) {
			if isAppliance(ctx) {
				ctx.Redirect(http.StatusTemporaryRedirect, "/appliance")

				return
			}

			ctx.Redirect(http.StatusTemporaryRedirect, "/app")
		})

		app := r.Group("/app")
		app.Any("/*path", func(ctx *gin.Context) {
			appFs, _ := fs.Sub(AppFS, "frontend/app")

			fs := http.FileServer(http.FS(appFs))
			http.StripPrefix("/app", fs).ServeHTTP(ctx.Writer, ctx.Request)
		})

		token, err := generateApplianceToken()
		if err != nil {
			os.Exit(1)
		}

		appliance := r.Group("/appliance")
		appliance.Any("/*path", func(ctx *gin.Context) {
			if !isAppliance(ctx) {
				ctx.Redirect(http.StatusTemporaryRedirect, "/app")
				return
			}

			appFs, _ := fs.Sub(ApplianceFS, "frontend/appliance")
			fs := http.FileServer(http.FS(appFs))

			// We need to inject the token in the index page
			if ctx.Request.URL.Path == "/appliance" || ctx.Request.URL.Path == "/appliance/" || strings.HasSuffix(ctx.Request.URL.Path, "/index.html") {
				content, err := appFs.Open("index.html")
				if err != nil {
					ctx.Status(http.StatusInternalServerError)
					return
				}
				defer content.Close()

				bytes, err := io.ReadAll(content)
				if err != nil {
					ctx.Status(http.StatusInternalServerError)
					return
				}

				tokenScript := fmt.Sprintf(`<script>window.MERCURE_TOKEN = "%s";</script>`, token)
				modified := strings.Replace(string(bytes), "</head>", tokenScript+"</head>", 1)

				ctx.Header("Content-Type", "text/html")
				ctx.String(http.StatusOK, modified)

				return
			}

			http.StripPrefix("/appliance", fs).ServeHTTP(ctx.Writer, ctx.Request)
		})

		apiRouter := r.Group("/api")

		routes.RegisterApplianceRoutes(apiRouter)
		routes.RegisterWebappRoutes(apiRouter)

		// Running the appliance
		log.LOG.Infow("Now listening...", "addr", config.GET.ListeningAddr)

		cron.RunCron()

		if err := r.Run(config.GET.ListeningAddr); err != nil {
			log.LOG.Fatalw("Failed to run the server", "err", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	userCmd.AddCommand(getUserCmd)
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(getInitializeUserCmd())

	devCmd.AddCommand(generateJwtCmd)
	devCmd.AddCommand(testCmd)

	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(hwHandlerCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(versionCmd)
}

func generateApplianceToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromAppliance())
	tokenString, err := token.SignedString(config.GET.Mercure.SubscriberKey)
	if err != nil {
		log.LOG.Errorw("Failed to generate JWT", "err", err)

		return "", err
	}

	return tokenString, nil
}
