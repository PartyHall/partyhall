package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/message_handler"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/socket"
	"github.com/partyhall/partyhall/utils"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Register(r *mux.Router) {
	r.HandleFunc("/socket/{type}", startSocket)
	r.HandleFunc("/settings", settings)
	r.HandleFunc("/picture", picture).Methods(http.MethodPost)

	registerAdminRoutes(r.PathPrefix("/admin").Subrouter())
}

func startSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	socketType := strings.ToUpper(vars["type"])
	if !slices.Contains(socket.SOCKET_TYPES, socketType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// PartyHall should not be allowed from another computer
	if socketType == socket.SOCKET_TYPE_BOOTH {
		if utils.IsRemote(r) {
			if !config.GET.DebugMode && !config.IsInDev() {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			logs.Debug("Letting a remote connection")
		}
	} else {
		pwd := r.URL.Query().Get("password")
		if pwd != config.GET.Web.AdminPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("Failed to upgrade connection: ", err)
		return
	}

	sock := socket.Join(socketType, conn)

	go func() {
		for {
			data := models.SocketMessage{}
			err := conn.ReadJSON(&data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logs.Error("Unexpected close error: ", err)
					sock.Open = false
					return
				} else if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logs.Error("Websocket disconnected: ", err)
					sock.Open = false
					return
				}

				logs.Error("Failed to connect websocket: ", err)
				continue
			}

			message_handler.ProcessMessage(sock, data)
		}
	}()
}

func settings(w http.ResponseWriter, r *http.Request) {
	settings := services.GET.GetFrontendSettings()

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(settings)
	w.Write(data)
}

func getEventAndFilename(event string, isUnattended bool) (int64, string) {
	var err error
	var eventId int64 = -1
	var imageName string = fmt.Sprintf("%v.jpg", time.Now().Format("20060102-150405"))

	eventId, err = strconv.ParseInt(event, 10, 64)
	if err != nil {
		logs.Error("Failed to get event id: ", err)
		logs.Error("Fallingback to id -1")
		eventId = -1
	}

	if eventId == -1 {
		return -1, imageName
	}

	evt, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		logs.Error("No event for the given id")
		return -1, imageName
	}

	img, err := orm.GET.Events.InsertImage(evt.Id, isUnattended)
	if err != nil {
		logs.Error("Failed to insert image: ", err)
		logs.Error("Defaulting name to current timestamp in the root folder for the event")
	} else {
		imageName = fmt.Sprintf("%v.jpg", img.Id)
	}

	return evt.Id, imageName
}

func picture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(256 * 1024) // Max picture size = 256mo, we should be good.
	if err != nil {
		logs.Error("Unable to save picture: Parse form error => ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event := r.FormValue("event")
	unattended := r.FormValue("unattended")
	image := r.FormValue("image")

	if len(event) == 0 || len(unattended) == 0 || len(image) == 0 {
		logs.Error("Failed to save picture: bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isUnattended, err := strconv.ParseBool(unattended)
	if err != nil {
		logs.Error("Failed to parse unattended var: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventId, filename := getEventAndFilename(event, isUnattended)

	path, err := config.GET.GetImageFolder(eventId, isUnattended)
	if err != nil {
		logs.Error("Failed to create path: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filepath := filepath.Join(path, filename)
	f, err := os.Create(filepath)
	if err != nil {
		logs.Error("Failed to create image file...")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	image = image[len("data:image/jpeg;base64,"):]
	data, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		logs.Error("Failed to decode image, writing it to file as-is")
		_, err = f.Write([]byte(image))
		if err != nil {
			logs.Error("Even failed to write the b64... sad")
		}
	} else {
		_, err = f.Write(data)
		if err != nil {
			logs.Error("Failed to write the image to disk")
		}
	}

	if err = f.Sync(); err != nil {
		logs.Error("Failed to sync the data ! be careful")
	}

	// Broadcasting the state so that the current event is refreshed on the admin panel
	socket.SOCKETS.BroadcastState()

	if !isUnattended {
		http.ServeFile(w, r, filepath)
	}
}
