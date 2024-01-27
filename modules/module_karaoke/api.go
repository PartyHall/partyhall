package module_karaoke

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/partyhall/partyhall/models"
)

type ApiResponse struct {
	Meta struct {
		Page    int `json:"page"`
		MaxPage int `json:"max_page"`
		Total   int `json:"total"`
	} `json:"meta"`
	Results []models.Song `json:"results"`
}

func searchSong(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	songs, err := ormSearchSong(query)
	if err != nil {
		jsonErr, _ := json.Marshal(map[string]interface{}{
			"err":     "Failed to search song",
			"details": err,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonErr)
		return
	}

	data, _ := json.Marshal(songs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func listSong(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	var page int64 = 1

	if len(pageStr) > 0 {
		var err error = nil
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			page = 1
		}
	}

	songs, err := ormListSongs((page-1)*int64(CONFIG.AmtSongsPerPage), int64(CONFIG.AmtSongsPerPage))
	if err != nil {
		jsonErr, _ := json.Marshal(map[string]interface{}{
			"err":     "Failed to list songs",
			"details": err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonErr)
		return
	}

	count, err := ormCountSongs()
	if err != nil {
		count = len(songs)
	}

	data, _ := json.Marshal(models.ContextualizedResponse{
		Results: songs,
		Meta: models.ResponseMetadata{
			LastPage: int(math.Ceil(float64(count) / float64(CONFIG.AmtSongsPerPage))),
			Total:    count,
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
