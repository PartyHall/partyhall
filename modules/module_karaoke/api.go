package module_karaoke

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type ApiResponse struct {
	Meta struct {
		Page    int `json:"page"`
		MaxPage int `json:"max_page"`
		Total   int `json:"total"`
	} `json:"meta"`
	Results []models.Song `json:"results"`
}

type SongResult struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
	Cover  string `json:"cover"`
}

var VALID_FORMATS = []string{"CDG", "WEBM", "MP4"}
var VALID_COVER_TYPE = []string{"NO_COVER", "UPLOADED", "LINK"}
var nonAsciiRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func createSong(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100 << 20)

	errors := map[string]string{}

	title := r.FormValue("title")
	if len(strings.TrimSpace(title)) == 0 || len(strings.TrimSpace(title)) > 32 {
		errors["title"] = "Title too long or not found"
	}

	artist := r.FormValue("artist")
	if len(strings.TrimSpace(artist)) == 0 || len(strings.TrimSpace(artist)) > 64 {
		errors["title"] = "Artist too long or not found"
	}

	format := r.FormValue("format")
	if !slices.Contains(VALID_FORMATS, format) {
		errors["format"] = "Invalid format, must be CDG, WEBM or MP4"
	}

	coverType := r.FormValue("cover_type")
	if !slices.Contains(VALID_COVER_TYPE, coverType) {
		errors["cover_type"] = "Invalid cover type, must be LINK, UPLOADED or NO_COVER"
	}

	songFile, _, err := r.FormFile("song")
	if err != nil {
		errors["song"] = "Missing song file"
	} else {
		defer songFile.Close()
	}

	var cdgFile multipart.File
	if format == "CDG" {
		cdgFile, _, err = r.FormFile("cdg")
		if err != nil {
			errors["cdg"] = "Missing CDG file"
		} else {
			defer cdgFile.Close()
		}
	}

	var coverFile multipart.File
	if coverType == "UPLOADED" {
		coverFile, _, err = r.FormFile("cover")
		if err != nil {
			errors["cover"] = "Missing cover file"
		} else {
			defer coverFile.Close()
		}
	}

	if len(errors) > 0 {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(errors)
		w.Write(data)

		return
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	titleFoldername, _, _ := transform.String(t, title)
	artistFoldername, _, _ := transform.String(t, artist)

	titleFoldername = strings.ReplaceAll(titleFoldername, " ", "")
	artistFoldername = strings.ReplaceAll(artistFoldername, " ", "")

	titleFoldername = nonAsciiRegex.ReplaceAllString(titleFoldername, "")
	artistFoldername = nonAsciiRegex.ReplaceAllString(artistFoldername, "")

	foldername := artistFoldername + "_" + titleFoldername

	err = ormCreateSong(foldername, artist, title, strings.ToLower(format))
	if err != nil {
		fmt.Println(err)

		errors["non_fields"] = "Failed to create song: " + err.Error()
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(errors)
		w.Write(data)

		return
	}

	tempDir, err := os.MkdirTemp("", "phkaraoke")
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to create temporary directory\"}"))

		return
	}

	songFilename := "song." + strings.ToLower(format)
	if format == "CDG" {
		songFilename = "song.mp3"
	}

	outputSong, err := os.Create(filepath.Join(tempDir, songFilename))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to create temporary song file\"}"))

		return
	}

	_, err = io.Copy(outputSong, songFile)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to copy song file\"}"))

		return
	}
	outputSong.Close()

	if format == "CDG" {
		outputCdg, err := os.Create(filepath.Join(tempDir, "song.cdg"))
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"non_fields\": \"Failed to create temporary cdg file\"}"))

			return
		}

		_, err = io.Copy(outputCdg, cdgFile)
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"non_fields\": \"Failed to copy cdg file\"}"))

			return
		}
		outputCdg.Close()
	}

	if coverType != "NO_COVER" {
		outputCover, err := os.Create(filepath.Join(tempDir, "cover.jpg"))
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"non_fields\": \"Failed to create temporary cover file\"}"))

			return
		}

		if coverType == "LINK" {
			resp, err := http.Get(r.FormValue("cover"))
			if err != nil {
				fmt.Println(err)

				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"non_fields\": \"Failed to create download cover file\"}"))

				return
			}
			defer resp.Body.Close()

			_, err = io.Copy(outputCover, resp.Body)
			if err != nil {
				fmt.Println(err)

				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"non_fields\": \"Failed to create copy downloaded cover file\"}"))

				return
			}
		} else if coverType == "UPLOADED" {
			_, err = io.Copy(outputCover, coverFile)
			if err != nil {
				fmt.Println(err)

				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"non_fields\": \"Failed to copy cover file\"}"))

				return
			}
		}

		outputCover.Close()

		// Be sure its 300x300 as jpg
		// Save it back on top of it
	}

	infoFile, err := os.Create(filepath.Join(tempDir, "info.txt"))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to create info.txt\"}"))

		return
	}
	infoFile.WriteString(artist + "\n" + title + "\n" + strings.ToLower(format) + "\n0")
	infoFile.Close()

	err = utils.CopyDir(tempDir, filepath.Join(config.GET.RootPath, "karaoke", foldername))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to move temp dir\"}"))

		return
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"non_fields\": \"Failed to remove temp dir\"}"))

		return
	}
}

func spotifySearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tracks, err := services.GET.Spotify.SearchSong(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	songs := []SongResult{}
	for _, t := range tracks {
		var bestImageUrl string = ""
		image := getBestImage(t.Album.Images)

		if image != nil {
			bestImageUrl = image.URL
		}

		song := SongResult{
			Artist: t.Artists[0].Name,
			Song:   t.Name,
			Cover:  bestImageUrl,
		}

		songs = append(songs, song)
	}

	data, _ := json.Marshal(songs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
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
