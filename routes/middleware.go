package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/partyhall/partyhall/config"
)

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
