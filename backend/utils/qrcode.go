package utils

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var AssetsFS embed.FS

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

func GenerateQrCode(txt string, c *gin.Context) {
	qr, err := qrcode.NewWith(
		txt,
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"err":            "Failed to generate QR Code",
			"internal_error": err.Error(),
		})

		return
	}

	var logo *image.Image = nil
	data, err := AssetsFS.ReadFile("assets/logo.png")
	if err == nil {
		logoImg, err := png.Decode(bytes.NewReader(data))
		if err == nil {
			logo = &logoImg
		}
	}

	wc := nopCloser{Writer: c.Writer}

	opts := []standard.ImageOption{
		standard.WithQRWidth(uint8(16)),
		standard.WithBgTransparent(),
		standard.WithBorderWidth(1),
		standard.WithFgColor(color.RGBA{
			R: 82,
			G: 14,
			B: 19,
			A: 255,
		}),
	}

	if logo != nil {
		opts = append(opts, standard.WithLogoImage(*logo))
		opts = append(opts, standard.WithLogoSizeMultiplier(1))
	}

	w := standard.NewWithWriter(
		wc,
		opts...,
	)

	if err := qr.Save(w); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"err":            "Failed to generate QR Code",
			"internal_error": err.Error(),
		})

		return
	}
}

func GenerateQrCodeWithoutLogo(txt string, c *gin.Context) {
	qr, err := qrcode.NewWith(
		txt,
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"err":            "Failed to generate QR Code",
			"internal_error": err.Error(),
		})

		return
	}

	wc := nopCloser{Writer: c.Writer}

	opts := []standard.ImageOption{
		standard.WithQRWidth(uint8(16)),
		standard.WithBgTransparent(),
		standard.WithBorderWidth(1),
		standard.WithFgColor(color.RGBA{
			R: 82,
			G: 14,
			B: 19,
			A: 255,
		}),
	}

	w := standard.NewWithWriter(
		wc,
		opts...,
	)

	if err := qr.Save(w); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"err":            "Failed to generate QR Code",
			"internal_error": err.Error(),
		})

		return
	}
}
