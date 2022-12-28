package utils

import (
	"os"
	"path/filepath"

	"github.com/partyhall/partyhall/config"
)

func GetPath(path ...string) string {
	path = append([]string{config.GET.RootPath}, path...)

	return filepath.Join(path...)
}

func MakeOrCreateFolder(path string) error {
	if _, err := os.Stat(GetPath(path)); os.IsNotExist(err) {
		err := os.MkdirAll(GetPath(path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
