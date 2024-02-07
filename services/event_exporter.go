package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
)

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
	err := orm.GET.Events.Save(ee.event)
	if err != nil {
		logs.Error("Failed to set the exporting state")
		return err
	}

	return nil
}

type RecapParams struct {
	PictureFolder  string
	OutputFilename string
	Framerate      int
}

func (ee EventExporter) BuildFfmpegCommand(params RecapParams) *exec.Cmd {
	args := []string{"ffmpeg"}

	args = append(args, "-framerate", fmt.Sprintf("%v", params.Framerate))
	args = append(args, "-pattern_type", "glob")
	args = append(args, "-i", "'*.jpg'")
	args = append(args, "-c:v", "libx264")

	args = append(args, params.OutputFilename)

	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	cmd.Dir = params.PictureFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (ee EventExporter) Export() (*models.ExportedEvent, error) {
	if ee.event.Exporting {
		return nil, errors.New("can't export an event that is already exporting")
	}

	if err := ee.setEventExporting(true); err != nil {
		return nil, err
	}

	exportTime := time.Now()

	basepath := fmt.Sprintf("images/%v/exports/", ee.event.Id)
	err := utils.MakeOrCreateFolder(basepath)
	if err != nil {
		return nil, errors.Join(err, ee.setEventExporting(false))
	}

	basefilepath := exportTime.Format("20060201-150405") + ".zip"
	filepath := utils.GetPath(basepath + basefilepath)

	archive, err := os.Create(filepath)
	if err != nil {
		return nil, errors.Join(err, ee.setEventExporting(false))
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	//#region Exporting non-unattended images
	images, err := orm.GET.Events.GetImages(ee.event, false)
	if err != nil {
		return nil, errors.Join(err, ee.setEventExporting(false))
	}

	for _, i := range images {
		imagePath := utils.GetPath(fmt.Sprintf("images/%v/pictures/%v.jpg", ee.event.Id, i.Id))
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			logs.Errorf("Failed to locate image %v from event %v\n", i.Id, ee.event.Id)
			continue
		}

		fr, err := os.Open(imagePath)
		if err != nil {
			logs.Errorf("Failed to open the image %v for the event %v: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		fw, err := zipWriter.Create((time.Time(i.Date)).Format("20060201-150405") + ".jpg")
		if err != nil {
			logs.Errorf("Failed to create the image %v for the event %v in the zip file: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		if _, err := io.Copy(fw, fr); err != nil {
			logs.Errorf("Failed to copy the image %v for the event %v in the zip file: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		fr.Close()
	}
	//#endregion

	//#region Exporting unattended images
	unattendedRoot := utils.GetPath(fmt.Sprintf("images/%v/unattended/", ee.event.Id))
	if _, err := os.Stat(unattendedRoot); !os.IsNotExist(err) {
		outvid := unattendedRoot + "/000_recap.mp4"
		if _, err := os.Stat(outvid); !os.IsNotExist(err) {
			os.Remove(outvid)
		}

		cmd := ee.BuildFfmpegCommand(RecapParams{
			PictureFolder:  unattendedRoot,
			OutputFilename: "000_recap.mp4",
			Framerate:      6,
		})
		err = cmd.Run()
		if err != nil {
			return nil, errors.Join(err, ee.setEventExporting(false))
		}

		fr, err := os.Open(outvid)
		if err != nil {
			logs.Errorf("Failed to open the recap video for the event %v: %v\n", ee.event.Id, err)
			return nil, errors.Join(err, ee.setEventExporting(false))
		}

		fw, err := zipWriter.Create("000_recap.mp4")
		if err != nil {
			logs.Errorf("Failed to create the recap video for the event %v in the zip file: %v\n", ee.event.Id, err)
			return nil, errors.Join(err, ee.setEventExporting(false))
		}

		if _, err := io.Copy(fw, fr); err != nil {
			logs.Errorf("Failed to copy the recap video for the event %v in the zip file: %v\n", ee.event.Id, err)
			return nil, errors.Join(err, ee.setEventExporting(false))
		}

		fr.Close()

		if _, err := os.Stat(outvid); !os.IsNotExist(err) {
			os.Remove(outvid)
		}
	} else {
		logs.Warn("Unattended folder doesn't exists, skipping...")
	}
	//#region

	//#region Adding a json file with some data about the event
	data := map[string]interface{}{
		"id":                    ee.event.Id,
		"name":                  ee.event.Name,
		"author":                ee.event.Author,
		"date":                  ee.event.Date,
		"location":              ee.event.Location,
		"amt_images_handtaken":  ee.event.AmtImagesHandtaken,
		"amt_images_unattended": ee.event.AmtImagesUnattended,
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	fw, err := zipWriter.Create("001_infos.json")
	if err != nil {
		logs.Errorf("Failed to add the info json")
		return nil, errors.Join(err, ee.setEventExporting(false))
	}

	if _, err := io.Copy(fw, bytes.NewReader(jsonData)); err != nil {
		logs.Errorf("Failed to copy the info json for the event %v in the zip file: %v\n", ee.event.Id, err)
		return nil, errors.Join(err, ee.setEventExporting(false))
	}
	//#endregion

	exportTimestamp := models.Timestamp(exportTime)
	ee.event.LastExport = &exportTimestamp
	err = ee.setEventExporting(false)
	if err != nil {
		return nil, err
	}

	// Insert the built
	exportedEvent, err := orm.GET.Events.InsertExportedEvent(ee.event, basefilepath)
	if err != nil {
		return nil, err
	}

	return exportedEvent, nil
}
