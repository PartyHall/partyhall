package utils

import (
	"os"
	"path/filepath"
)

var ROOT_PATH string = ""

func GetPath(path ...string) string {
	path = append([]string{ROOT_PATH}, path...)

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
