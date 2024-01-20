package module_karaoke

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VideoConverter struct {
	defaultFormat string
}

type StreamInfo struct {
	CodecName string `json:"codec_name"`
}

type FFprobeOutput struct {
	Streams []StreamInfo `json:"streams"`
}

func (vc *VideoConverter) ConvertForRaspberryPi(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	isAlreadyOk, err := vc.CheckFormat(path)
	if err != nil {
		return "", err
	}

	if isAlreadyOk {
		return path, nil
	}

	ext := filepath.Ext(filepath.Base(path))
	outpath := strings.TrimSuffix(path, ext) + "_converted"

	ffmpegCmd, newExt := getCommandForFormat(vc.defaultFormat)
	ffmpegCmd = fmt.Sprintf(
		ffmpegCmd,
		strings.ReplaceAll(path, `"`, `\"`),
		strings.ReplaceAll(outpath, `"`, `\"`),
	)

	cmd := exec.Command("bash", "-c", ffmpegCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return outpath + "." + newExt, cmd.Run()
}

func (vc *VideoConverter) CheckFormat(path string) (bool, error) {
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			`ffprobe -v quiet -print_format json -show_entries stream=codec_name "%v"`,
			strings.ReplaceAll(path, `"`, `\"`),
		),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	var data FFprobeOutput
	err = json.Unmarshal(output, &data)
	if err != nil {
		return false, err
	}

	foundVideoFormat := false
	for _, stream := range data.Streams {
		if stream.CodecName == vc.defaultFormat {
			foundVideoFormat = true
			break
		}
	}

	return foundVideoFormat, nil
}

func getCommandForFormat(format string) (string, string) {
	if format == "vp9" {
		return `ffmpeg -i "%v" -c:v libvpx-vp9 -vf "scale=min(720\,iw):-2" -c:a libopus -b:a 192k "%v.webm"`, "webm"
	}

	if format == "h264" {
		return `ffmpeg -i "%v" -c:v libx264 -vf "scale=-2:min(720\,ih)" -c:a aac -strict experimental -b:a 192k "%v.mp4"`, "mp4"
	}

	if format == "h265" {
		return `ffmpeg -i "%v" -c:v libx265 -vf "scale=min(720\,iw):-2" -preset medium -crf 28 -c:a aac -strict experimental -b:a 128k "%v.mp4"`, "mp4"
	}

	return "", ""
}
