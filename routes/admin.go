package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/remote"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

func registerAdminRoutes(r *mux.Router) {
	r.Use(AuthenticatedMiddleware)

	r.HandleFunc("/login", isPasswordValid)
	r.HandleFunc("/set_mode/{mode}", setMode).Methods(http.MethodPost)

	r.HandleFunc("/exports/{event_id}", getLastExports).Methods(http.MethodGet)
	r.HandleFunc("/exports/{export_id}/download", downloadExport).Methods(http.MethodGet)

	eventRouter := r.PathPrefix("/event/{event}").Subrouter()
	r.HandleFunc("/event", createEvent).Methods(http.MethodPost)
	eventRouter.Use(WithEventMiddleware)

	eventRouter.HandleFunc("", updateEvent).Methods(http.MethodPut)

	imageRouter := eventRouter.PathPrefix("/image/{image}").Subrouter()
	imageRouter.Use(WithImageMiddleware)

	imageRouter.HandleFunc("", serveImage).Methods(http.MethodGet)
}

func isPasswordValid(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("yes"))
}

func setMode(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	(*services.GET.MqttClient).Publish("partyhall/admin/set_mode", 2, false, values["mode"])
}

type UpdatedEvent struct {
	Name     string  `json:"name"`
	Author   *string `json:"author"`
	Location *string `json:"location"`
	Date     *int64  `json:"date"`
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	updatedEvent := UpdatedEvent{}
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		logs.Error("Failed to parse new event: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	evt, err := orm.GET.Events.CreateEvent(updatedEvent.Name, updatedEvent.Author, updatedEvent.Location, updatedEvent.Date)
	if err != nil {
		if err == sqlite3.ErrConstraint {
			w.WriteHeader(http.StatusBadRequest)
			data, _ := json.MarshalIndent("Some data are missing", "", "  ")
			w.Write(data)
			logs.Error("Missing something: ", err)
			return
		}

		logs.Error("Failed to save new event in DB: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	//#region If we just created the first event, we automatically select it
	events, err := orm.GET.Events.GetEvents()
	if err == nil {
		if len(events) == 1 {
			services.GET.CurrentState.CurrentEvent = &events[0].Id
			services.GET.CurrentState.CurrentEventObj = &events[0]

			orm.GET.AppState.SetState(services.GET.CurrentState)
		}
	}
	//#endregion

	remote.BroadcastState()

	data, _ := json.MarshalIndent(evt, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event := ctx.Value(ContextEventKey).(*models.Event)

	updatedEvent := UpdatedEvent{}
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Error(err)
		return
	}

	event.Name = updatedEvent.Name
	event.Author = updatedEvent.Author
	event.Location = updatedEvent.Location
	if updatedEvent.Date != nil {
		date := models.Timestamp(time.Unix(*updatedEvent.Date, 0))
		event.Date = &date
	}

	err = orm.GET.Events.Save(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logs.Error("Failed to save an updated event: ", err)
		return
	}

	remote.BroadcastState()
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event := ctx.Value(ContextEventKey).(*models.Event)
	image := ctx.Value(ContextImageKey).(*models.Image)

	subpath := "pictures"
	if image.Unattended {
		subpath = "unattended"
	}

	path := utils.GetPath("image", fmt.Sprintf("%v", event.Id), "images", subpath, fmt.Sprintf("%v", image.Id)) + ".jpg"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		WriteError(w, errors.New("image not found"), http.StatusNotFound, "Can't find file on filesystem")
		return
	}

	http.ServeFile(w, r, path)
}

func getLastExports(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	event_id := values["event_id"]

	eventId, err := strconv.ParseInt(event_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		// @Todo: handle better
		w.WriteHeader(http.StatusNotFound)
		return
	}

	exportedEvents, err := orm.GET.Events.GetExportedEvents(event, 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logs.Error(err)
		return
	}

	data, _ := json.MarshalIndent(exportedEvents, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func downloadExport(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	exportedEventStr := values["export_id"]

	exportedEventId, err := strconv.ParseInt(exportedEventStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exportedEvent, err := orm.GET.Events.GetExportedEvent(exportedEventId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logs.Error(err)

		return
	}

	path := utils.GetPath("images", fmt.Sprintf("%v", exportedEvent.EventId), "exports", exportedEvent.Filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		WriteError(w, errors.New("export not found"), http.StatusNotFound, "Can't find file on filesystem")
		return
	}

	http.ServeFile(w, r, path)
}
