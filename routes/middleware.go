package routes

import (
	"encoding/json"
	"net/http"
	"strings"
)

// @TODO deprecate or rework for echo
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
