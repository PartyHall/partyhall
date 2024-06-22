package module_karaoke

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	_ "embed"

	"github.com/fogleman/gg"
	"github.com/h2non/bimg"
	"github.com/partyhall/partyhall/services"
)

/**
@TODO: Vérifier si on peut pas minimiser l'usage de bimg pour faire un max dans gg
La lib a l'air puissante et peut donc potentiellement gérer pas mal de trucs

Intéressant, à voir: https://github.com/u2takey/ffmpeg-go / https://github.com/csnewman/ffmpeg-go / https://github.com/xfrr/goffmpeg

Pour l'instant ce code est embed dans le photomaton, à terme il faudra un rabbitmq pour que le php délègue le taf a un microservice de génération de story
*/

type OverlayGenerator struct {
	logo         *image.Image
	logoSize     bimg.ImageSize
	finalWidth   int
	finalHeight  int
	coverArtSize int
	fontTtf      string
	fontSize     float64

	Padding int
	FilterY int
}

func NewOverlayGenerator(finalWidth int, coverArtSize int, fontTtf string, fontSize float64) (*OverlayGenerator, error) {
	finalHeight := int(finalWidth * 16 / 9)

	logoResized, err := bimg.NewImage(services.LOGO_IMAGE).Resize(finalWidth/2, 0)
	if err != nil {
		return nil, err
	}

	logoImg, err := png.Decode(bytes.NewReader(logoResized))
	if err != nil {
		return nil, err
	}

	logoSize, err := bimg.NewImage(logoResized).Size()
	if err != nil {
		return nil, err
	}

	return &OverlayGenerator{
		logo:         &logoImg,
		logoSize:     logoSize,
		finalWidth:   finalWidth,
		finalHeight:  finalHeight,
		coverArtSize: coverArtSize,
		fontTtf:      fontTtf,
		fontSize:     fontSize,
		Padding:      20,
		FilterY:      250,
	}, nil
}

func (og *OverlayGenerator) LoadCoverArt(songFile string) (*image.Image, error) {
	coverArt, err := readFileInZip(songFile, "cover.jpg")
	if err != nil {
		return nil, err
	}

	coverArtImage, err := bimg.NewImage(coverArt).Resize(og.coverArtSize, og.coverArtSize)
	if err != nil {
		return nil, err
	}

	coverArtPng, err := bimg.NewImage(coverArtImage).Convert(bimg.PNG)
	if err != nil {
		return nil, err
	}

	coverImg, err := png.Decode(bytes.NewReader(coverArtPng))
	if err != nil {
		return nil, err
	}

	return &coverImg, nil
}

func (og OverlayGenerator) Generate(phk PhkSong, output string) error {
	songFile, err := getModuleFile(phk.Filename)
	if err != nil {
		return err
	}

	coverArtImg, err := og.LoadCoverArt(songFile)
	if err != nil {
		return err
	}

	// We create an empty transparent image
	baseImg := image.NewRGBA(image.Rect(0, 0, og.finalWidth, og.finalHeight))
	draw.Draw(baseImg, baseImg.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src) // Not sure if required

	// We draw the Partyhall logo
	if og.logo != nil {
		draw.Draw(
			baseImg,
			image.Rect(
				(og.finalWidth-og.logoSize.Width)/2,
				0,
				(og.finalWidth+og.logoSize.Width)/2,
				og.logoSize.Height,
			),
			*og.logo,
			image.Point{},
			draw.Over,
		)
	}

	fImg := gg.NewContext(og.finalWidth, og.finalHeight)
	if err := fImg.LoadFontFace(og.fontTtf, og.fontSize); err != nil {
		return err
	}

	// @TODO: Multiline ?
	artistWidth, _ := fImg.MeasureString(phk.Artist)
	titleWidth, _ := fImg.MeasureString(phk.Title)
	maxTextWidth := artistWidth
	if titleWidth > maxTextWidth {
		maxTextWidth = titleWidth
	}

	songFilterWidth := og.coverArtSize + (og.Padding * 3) + int(maxTextWidth)
	songFilterHeight := og.coverArtSize + (og.Padding * 2)

	// Draw black background of the filter at 0.8 opacity
	fImg.SetRGBA(0, 0, 0, 0.8)
	fImg.DrawRoundedRectangle(
		float64((og.finalWidth-songFilterWidth)/2),
		float64(og.finalHeight-og.FilterY-(songFilterHeight/2)),
		float64(songFilterWidth),
		float64(songFilterHeight),
		20,
	)
	fImg.Fill()

	// Draw the cover art
	coverArtX := (og.finalWidth-songFilterWidth)/2 + og.Padding
	coverArtY := og.finalHeight - og.FilterY - (songFilterHeight / 2) + og.Padding
	fImg.DrawImageAnchored(*coverArtImg, coverArtX, coverArtY, 0, 0)

	availableSpace := (songFilterWidth - (3 * og.Padding) - og.coverArtSize)

	textX := float64(coverArtX + og.coverArtSize + og.Padding + (availableSpace / 2))
	textY := float64(coverArtY + (og.coverArtSize / 2) - 20)

	// Draw the artist
	fImg.SetRGB(1, 1, 1)
	fImg.DrawStringAnchored(
		phk.Artist,
		textX,
		textY,
		0.5,
		0.5,
	)

	fImg.DrawStringAnchored(
		phk.Title,
		textX,
		textY+40,
		0.5,
		0.5,
	)

	// Draw the gg thing on the original image
	draw.Draw(baseImg, baseImg.Bounds(), fImg.Image(), image.Point{}, draw.Over)

	// Save the final image
	finalFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	if err := png.Encode(finalFile, baseImg); err != nil {
		return err
	}

	return nil
}
