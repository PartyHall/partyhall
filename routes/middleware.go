package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
)

type ContextKey string

const ContextEventKey ContextKey = "EVENT"
const ContextImageKey ContextKey = "IMAGE"

func WriteError(w http.ResponseWriter, err error, errCode int, customTxt string) bool {
	if err != nil {
		metadata := map[string]string{
			"error":   strings.ReplaceAll(err.Error(), "\"", "\\\""),
			"details": customTxt,
		}

		data, _ := json.Marshal(metadata)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errCode)
		w.Write(data)

		return true
	}

	return false
}

func AuthenticatedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := r.Header.Get("Authorization")

		if len(password) == 0 || config.GET.Web.AdminPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func WithEventMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if val, ok := vars["event"]; !ok || len(val) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(vars["event"], 10, 64)
		if WriteError(w, err, http.StatusBadRequest, "Invalid event id") {
			return
		}

		evt, err := orm.GET.Events.GetEvent(id)
		if WriteError(w, err, http.StatusNotFound, "Failed to get event") {
			return
		}

		ctx := context.WithValue(r.Context(), ContextEventKey, evt)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithImageMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		evt := r.Context().Value(ContextEventKey).(*models.Event)

		vars := mux.Vars(r)
		if val, ok := vars["image"]; !ok || len(val) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(vars["image"], 10, 64)
		if WriteError(w, err, http.StatusBadRequest, "Invalid picture ID") {
			return
		}

		picture, err := orm.GET.Events.GetImage(id)
		if WriteError(w, err, http.StatusNotFound, "Failed to get picture") {
			return
		}

		if picture.EventId != evt.Id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`
				{
					"error": "not found",
					"details": "Image not found for this event"
				}
			`))
			return
		}

		ctx := context.WithValue(r.Context(), ContextImageKey, picture)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
