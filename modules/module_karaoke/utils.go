package module_karaoke

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

var songFilenameCache = map[string]string{}

func getModuleEventDir() (string, error) {
	evt := services.GET.CurrentState.CurrentEvent
	eventId := -1

	if evt != nil {
		eventId = *evt
	}

	basePath, err := utils.GetEventFolder(eventId)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(basePath, "karaoke")
	err = os.MkdirAll(fullPath, os.ModePerm)

	return fullPath, err
}

func getModuleDir() (string, error) {
	fullPath := filepath.Join(config.GET.RootPath, "karaoke")
	err := os.MkdirAll(fullPath, os.ModePerm)

	return fullPath, err
}

func getModuleFile(filename string) (string, error) {
	baseDir, err := getModuleDir()
	if err != nil {
		return baseDir, err
	}

	return filepath.Join(baseDir, filename), nil
}

func getBestImage(images []services.SpotifyImage) (bestImage *services.SpotifyImage) {
	// Sort images by resolution (descending order)
	sort.Slice(images, func(i, j int) bool {
		areaI := images[i].Width * images[i].Height
		areaJ := images[j].Width * images[j].Height
		return areaI > areaJ
	})

	// Find the first image with size 300x300
	for _, img := range images {
		if img.Width == 300 && img.Height == 300 {
			return &img
		}
	}

	// If no 300x300 image, return the highest resolution image
	if len(images) > 0 {
		return &images[0]
	}

	// Return an empty image if the input array is empty
	return nil
}

func streamFileFromZip(filename string, mimetype string) func(echo.Context) error {
	return func(c echo.Context) error {
		songUuid := c.Param("uuid")
		songFilename, found := songFilenameCache[songUuid]
		if !found {
			song, err := ormLoadSongByUuid(songUuid)
			if err != nil {
				return c.String(http.StatusNotFound, "Song not found: "+err.Error())
			}

			songFilenameCache[songUuid] = song.Filename
		}

		filepath, err := getModuleFile(songFilename)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get the song file: "+err.Error())
		}

		zipFile, err := os.Open(filepath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to open the song file: "+err.Error())
		}
		defer zipFile.Close()

		fileInfo, err := zipFile.Stat()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get the song file size: "+err.Error())
		}

		zipReader, err := zip.NewReader(zipFile, fileInfo.Size())
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to read the song file: "+err.Error())
		}

		for _, file := range zipReader.File {
			if file.Name == filename {
				fileReader, err := file.Open()
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to open file in the song: "+err.Error())
				}

				defer fileReader.Close()

				c.Response().Header().Set(echo.HeaderContentType, mimetype)
				c.Response().WriteHeader(http.StatusOK)
				_, err = io.Copy(c.Response().Writer, fileReader)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to stream file: "+err.Error())
				}
				return nil
			}
		}

		return c.String(http.StatusNotFound, "File not found in the song file")
	}
}

func LoadPhkSong(path string) (*PhkSong, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	phk := PhkSong{}

	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	found := false
	for _, file := range reader.File {
		if file.Name == "song.json" {
			found = true

			rc, err := file.Open()
			if err != nil {
				return nil, errors.New("failed to read PHK Song file: " + err.Error())
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)

			err = json.Unmarshal(buf.Bytes(), &phk)
			if err != nil {
				return nil, errors.New("failed to read PHK Song file: " + err.Error())
			}

			rc.Close()
		}
	}

	if !found {
		return nil, errors.New("failed to read PHK Song file: No song.json found")
	}

	return &phk, nil
}
