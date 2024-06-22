package module_karaoke

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"

	"github.com/partyhall/partyhall/models"
)

type Exporter struct {
	basePath string
	event    *models.Event
}

type TimelapseParams struct {
	PictureFolder string
	OutputPath    string
	Framerate     int
	OverlayImage  string
	Contains      bool
}

// I don't have time to do this properly for now
// It should be a single ffmpeg command, not multiple weird stuff
func (e Exporter) BuildFfmpegCommand(params TimelapseParams) *exec.Cmd {
	args := []string{"ffmpeg"}

	args = append(args, "-framerate", fmt.Sprintf("%v", params.Framerate))
	args = append(args, "-pattern_type", "glob")
	args = append(args, "-i", "'*.jpg'")
	args = append(args, "-c:v", "libx264")
	args = append(args, "-r", fmt.Sprintf("%v", params.Framerate))
	args = append(args, "-t", "10")

	args = append(args, params.OutputPath)

	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	cmd.Dir = params.PictureFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// #region chatgpt / Claude stupid ass stuff
func processImage(inputPath, overlayPath, outputPath string) error {
	// Open the input image
	inputFile, err := os.Open(inputPath) // Replace with your input image path
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return err
	}
	defer inputFile.Close()

	// Decode the input image
	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("Error decoding input image: %v\n", err)
		return err
	}

	overlayFile, err := os.Open(overlayPath) // Replace with your foreground image path
	if err != nil {
		fmt.Printf("Error opening foreground file: %v\n", err)
		return err
	}
	defer overlayFile.Close()

	// Decode the foreground image
	overlayImg, _, err := image.Decode(overlayFile)
	if err != nil {
		fmt.Printf("Error decoding foreground image: %v\n", err)
		return err
	}

	// Create a new context with the desired dimensions
	dc := gg.NewContext(1080, 1920)

	// Get the dimensions of the input image
	inputWidth := inputImg.Bounds().Dx()
	inputHeight := inputImg.Bounds().Dy()

	// Calculate the scaling factor to fit the image within the canvas
	scale := min(1080.0/float64(inputWidth), 1920.0/float64(inputHeight))

	// Calculate the new dimensions
	newWidth := int(float64(inputWidth) * scale)
	newHeight := int(float64(inputHeight) * scale)

	// Calculate the position to center the image
	x := (1080 - newWidth) / 2
	y := (1920 - newHeight) / 2

	dc.SetColor(color.Black)
	dc.DrawRectangle(0, 0, 1080, 1920)
	dc.Fill()
	// Draw the scaled image onto the context
	dc.DrawImage(inputImg, x, y)
	dc.DrawImage(overlayImg, 0, 0)

	// Save the result
	baseImg := image.NewRGBA(image.Rect(0, 0, 1080, 1920))
	draw.Draw(baseImg, baseImg.Bounds(), dc.Image(), image.Point{}, draw.Over)

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Encode the RGBA image to JPEG and write it to the file
	if err := jpeg.Encode(outFile, baseImg, nil); err != nil {
		return err
	}

	return nil
}

func processImages(tempDir, dir, overlayFile string) error {
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		os.RemoveAll(tempDir)
	}

	os.MkdirAll(tempDir, os.ModePerm)

	defer os.RemoveAll(tempDir) // Supprimer le dossier temporaire à la fin

	// Lire toutes les images du répertoire original
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du répertoire original : %s\n", err)
		return err
	}

	// Parcourir chaque fichier image dans le répertoire original
	for _, file := range files {
		if !file.IsDir() {
			// Chemin complet de l'image d'origine
			imagePath := filepath.Join(dir, file.Name())

			// Chemin complet pour sauvegarder l'image modifiée dans le dossier temporaire avec le même nom
			outputPath := filepath.Join(tempDir, file.Name())

			// Redimensionner et ajouter l'overlay
			err := processImage(imagePath, overlayFile, outputPath)
			if err != nil {
				fmt.Printf("Erreur lors du traitement de l'image %s : %s\n", file.Name(), err)
			}
		}
	}

	return nil
}

//#endregion

func (e Exporter) Export() (map[string]any, error) {
	metadata := map[string]any{}

	fmt.Println("Exporting karaoke to " + e.basePath)
	songs, err := ormLoadCompleteSessions(e.event.Id)
	if err != nil {
		return nil, err
	}

	basePath, err := getModuleEventDir(e.event.Id)
	if err != nil {
		return nil, err
	}

	og, err := NewOverlayGenerator(1080, 200, "JetBrainsMono-Regular.ttf", 30)
	if err != nil {
		return nil, err
	}

	metadata["amt_songs_played"] = len(songs)
	songsMetadata := []map[string]any{}
	for _, s := range songs {
		// We add the metadata to the main manifest
		songsMetadata = append(songsMetadata, s.AsExportMetadata())

		sessionPath := filepath.Join(basePath, fmt.Sprintf("%v", s.Id))
		if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
			fmt.Printf("Session %v has no pictures!\n", s.Id)
			continue
		}

		// We build the overlay
		tempOverlayFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
		err := og.Generate(s.Song, tempOverlayFile)
		if err != nil {
			fmt.Printf("Failed to generate overlay for session %v: %v", s.Id, err)
			continue
		}

		// We build the timelapse
		/*
			pictPath, err = e.buildTempFiles(TimelapseParams{
				PictureFolder: sessionPath,
				OutputPath:    filepath.Join(e.basePath, fmt.Sprintf("%v.mp4", s.Id)),
				Framerate:     6,
				OverlayImage:  tempOverlayFile,
				Contains:      true,
			})
			if err != nil {
				fmt.Printf("Failed to generate pictures for session %v: %v", s.Id, err)
				continue
			}
		*/

		tempDir := os.TempDir()
		tempDir = filepath.Join(tempDir, "img_proc")
		err = processImages(tempDir, sessionPath, tempOverlayFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		cmd := e.BuildFfmpegCommand(TimelapseParams{
			PictureFolder: tempDir,
			OutputPath:    filepath.Join(e.basePath, fmt.Sprintf("%v.mp4", s.Id)),
			Framerate:     6,
			OverlayImage:  tempOverlayFile,
			Contains:      true,
		})

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			continue
		}

		if _, err := os.Stat(tempOverlayFile); !os.IsNotExist(err) {
			os.Remove(tempOverlayFile)
		}
	}

	metadata["songs"] = songsMetadata

	return metadata, nil
}
