package module_photobooth

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/utils"
)

type RecapParams struct {
	PictureFolder  string
	OutputFilename string
	Framerate      int
}

type Exporter struct {
	basePath string
	event    *models.Event
}

func (e Exporter) BuildFfmpegCommand(params RecapParams) *exec.Cmd {
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

func (e Exporter) Export() (map[string]any, error) {
	metadata := map[string]any{}

	//#region Exporting non-unattended images
	images, err := orm.GET.Events.GetImages(e.event, false)
	if err != nil {
		return nil, err
	}

	for _, i := range images {
		imagePath := utils.GetPath(fmt.Sprintf("events/%v/photobooth/pictures/%v.jpg", e.event.Id, i.Id))
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			logs.Errorf("Failed to locate image %v from event %v\n", i.Id, e.event.Id)
			continue
		}

		fr, err := os.Open(imagePath)
		if err != nil {
			logs.Errorf("Failed to open the image %v for the event %v: %v\n", i.Id, e.event.Id, err)
			continue
		}

		output, err := os.Create(filepath.Join(e.basePath, (time.Time(i.Date)).Format("20060201-150405")+".jpg"))
		if err != nil {
			logs.Errorf("Failed to create the image %v for the event %v in the zip file: %v\n", i.Id, e.event.Id, err)
			continue
		}

		if _, err := io.Copy(output, fr); err != nil {
			logs.Errorf("Failed to copy the image %v for the event %v in the zip file: %v\n", i.Id, e.event.Id, err)
			continue
		}

		output.Close()
		fr.Close()
	}
	//#endregion

	//#region Exporting unattended images
	unattendedRoot := utils.GetPath(fmt.Sprintf("events/%v/photobooth/unattended/", e.event.Id))
	if _, err := os.Stat(unattendedRoot); !os.IsNotExist(err) {
		outvid := filepath.Join(e.basePath, "000_recap.mp4")
		if _, err := os.Stat(outvid); !os.IsNotExist(err) {
			os.Remove(outvid)
		}

		cmd := e.BuildFfmpegCommand(RecapParams{
			PictureFolder:  unattendedRoot,
			OutputFilename: outvid,
			Framerate:      6,
		})
		err = cmd.Run()
		if err != nil {
			if _, err := os.Stat(outvid); !os.IsNotExist(err) {
				os.Remove(outvid)
			}

			return nil, err
		}
	} else {
		logs.Warn("Unattended folder doesn't exists, skipping...")
	}
	//#endregion

	metadata["amt_images_handtaken"] = e.event.AmtImagesHandtaken
	metadata["amt_images_unattended"] = e.event.AmtImagesUnattended

	return metadata, nil
}
