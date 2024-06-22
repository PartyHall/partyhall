package services

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
)

var RequestModuleExport func(string, *models.Event) (map[string]any, error)

/**
 * @TODO: Each module should process its own export
 * i.e. each module should receive the zip file and add its own file as it wants
 * This would permit some more info to be exported
 * And this will be easier once Prowty is released to integrate with it
 */

type EventExporter struct {
	event *models.Event
}

func NewEventExporter(event *models.Event) EventExporter {
	return EventExporter{event}
}

func (ee EventExporter) setEventExporting(exp bool) error {
	ee.event.Exporting = exp
	err := orm.GET.Events.SaveEvent(ee.event)
	if err != nil {
		logs.Error("Failed to set the exporting state")
		return err
	}

	return nil
}

/**
 * This will be rewritten in the Partyhall server
 * Don't waste too much time on this
 **/
func (ee EventExporter) Export() (*models.ExportedEvent, error) {
	if ee.event.Exporting {
		return nil, errors.New("can't export an event that is already exporting")
	}

	exportTime := time.Now()

	if err := ee.setEventExporting(true); err != nil {
		return nil, err
	}

	tempPath, err := os.MkdirTemp("", "phexport")
	if err != nil {
		return nil, err
	}

	metadata, err := RequestModuleExport(tempPath, ee.event)
	if err != nil {
		fmt.Println("Failed to fully export stuff: ", err)
	}

	//#region Adding a json file with some data about the event
	data := map[string]interface{}{
		"id":       ee.event.Id,
		"name":     ee.event.Name,
		"author":   ee.event.Author,
		"date":     ee.event.Date,
		"location": ee.event.Location,
	}

	for k, v := range metadata {
		data[k] = v
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	err = os.WriteFile(filepath.Join(tempPath, "infos.json"), jsonData, os.ModePerm)
	if err != nil {
		logs.Errorf("Failed to copy the info json for the event %v in the zip file: %v\n", ee.event.Id, err)
		return nil, errors.Join(err, ee.setEventExporting(false))
	}
	//#endregion

	basepath := fmt.Sprintf("events/%v/exports", ee.event.Id)
	err = utils.MakeOrCreateFolder(basepath)
	if err != nil {
		return nil, errors.Join(err, ee.setEventExporting(false))
	}

	filename := exportTime.Format("20060201-150405") + ".zip"
	fullpath := utils.GetPath(basepath, filename)
	err = zipDir(tempPath, fullpath)
	if err != nil {
		ee.setEventExporting(false)

		return nil, err
	}

	exportTimestamp := models.Timestamp(exportTime)
	ee.event.LastExport = &exportTimestamp
	err = ee.setEventExporting(false)
	if err != nil {
		return nil, err
	}

	// Insert the built
	exportedEvent, err := orm.GET.Events.InsertExportedEvent(ee.event, filename)
	if err != nil {
		return nil, err
	}

	return exportedEvent, nil
}

func zipDir(sourceDir, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Obtenir le chemin relatif pour conserver la structure des dossiers dans l'archive ZIP
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Ignorer le répertoire source lui-même
		if relativePath == "." {
			return nil
		}

		if info.IsDir() {
			_, err = zipWriter.Create(relativePath + "/")
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			writer, err := zipWriter.Create(relativePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
