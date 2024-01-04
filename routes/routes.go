package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/modules"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
)

func Register(r *mux.Router) {
	r.HandleFunc("/socket/{type}", remote.EasyWS.Route)
	r.HandleFunc("/settings", settings)
	r.HandleFunc("/picture", picture).Methods(http.MethodPost)

	registerAdminRoutes(r.PathPrefix("/admin").Subrouter())

	modules.RegisterRoutes(r.PathPrefix("/modules").Subrouter())
}

func settings(w http.ResponseWriter, r *http.Request) {
	settings := services.BuildFrontendSettings()

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
	remote.BroadcastState()

	if !isUnattended {
		http.ServeFile(w, r, filepath)
	}
}
