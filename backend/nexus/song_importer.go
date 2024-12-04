package nexus

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/utils"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func importSong(phkPath string) error {
	karaokeBasePath := filepath.Join(config.GET.RootPath, "karaoke")
	if !fileExists(karaokeBasePath) {
		err := os.MkdirAll(karaokeBasePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	extractedPath, err := utils.ExtractZip(phkPath)

	// We remove before checking the error
	// As we want it to be remove both when it works and
	// when it doesn't
	os.RemoveAll(phkPath)

	if err != nil {
		return err
	}

	metadataFilePath := filepath.Join(extractedPath, "song.json")
	if _, err := os.Stat(metadataFilePath); os.IsNotExist(err) {
		return errors.New("invalid song: no song.json metadata file")
	}

	var metadata PhkSong
	metadataFileContent, err := os.ReadFile(metadataFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(metadataFileContent, &metadata)
	if err != nil {
		return err
	}

	log.Info("[SongImport] Parsed metadata", "id", metadata.NexusId, "title", metadata.Title, "artist", metadata.Artist)
	destPath := filepath.Join(karaokeBasePath, metadata.NexusId)
	os.MkdirAll(destPath, os.ModePerm)

	dbSong := models.Song{
		NexusId:  metadata.NexusId,
		Title:    metadata.Title,
		Artist:   metadata.Artist,
		Format:   strings.ToLower(metadata.Format),
		Duration: metadata.Duration,
	}

	coverFilePath := filepath.Join(extractedPath, "cover.jpg")
	mp3FilePath := filepath.Join(extractedPath, "instrumental.mp3")
	cdgFilePath := filepath.Join(extractedPath, "lyrics.cdg")
	videoFilePath := filepath.Join(extractedPath, "instrumental.webm")
	vocalsFilePath := filepath.Join(extractedPath, "vocals.mp3")
	combinedFilePath := filepath.Join(extractedPath, "full.mp3")

	dbSong.HasCover = fileExists(coverFilePath)
	dbSong.HasVocals = fileExists(vocalsFilePath)
	dbSong.HasCombined = fileExists(combinedFilePath)

	if metadata.Hotspot != nil {
		dbSong.Hotspot = models.JsonnableNullInt64(sql.NullInt64{
			Int64: *metadata.Hotspot,
			Valid: true,
		})
	} else {
		dbSong.Hotspot = models.JsonnableNullInt64(sql.NullInt64{
			Int64: 0,
			Valid: false,
		})
	}

	// CDG = MP3+CDG
	if dbSong.Format == "cdg" {
		if !fileExists(mp3FilePath) || !fileExists(cdgFilePath) {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return errors.New("the song phk is invalid: missing instrumental.mp3 or lyrics.cdg")
		}

		err = utils.MoveFileCrossDrive(
			mp3FilePath,
			filepath.Join(destPath, "instrumental.mp3"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}

		err = utils.MoveFileCrossDrive(
			cdgFilePath,
			filepath.Join(destPath, "lyrics.cdg"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}

		// Otherwise = video
	} else {
		if !fileExists(videoFilePath) {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return errors.New("the song phk is invalid: missing instrumental.webm")
		}

		err = utils.MoveFileCrossDrive(
			videoFilePath,
			filepath.Join(destPath, "instrumental.webm"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}
	}

	if dbSong.HasCover {
		err = utils.MoveFileCrossDrive(
			coverFilePath,
			filepath.Join(destPath, "cover.jpg"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}
	}

	if dbSong.HasVocals {
		err = utils.MoveFileCrossDrive(
			vocalsFilePath,
			filepath.Join(destPath, "vocals.mp3"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}
	}

	if dbSong.HasCombined {
		err = utils.MoveFileCrossDrive(
			combinedFilePath,
			filepath.Join(destPath, "combined.mp3"),
		)

		if err != nil {
			os.RemoveAll(extractedPath)
			os.RemoveAll(destPath)
			return err
		}
	}

	err = dal.SONGS.Create(&dbSong)
	if err != nil {
		os.RemoveAll(extractedPath)
		os.RemoveAll(destPath)
		return err
	}

	return nil
}
