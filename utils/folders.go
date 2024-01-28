package utils

import (
	"io"
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

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func FileExistsForAnyExt(basename string, allowedExtensions []string) string {
	for _, v := range allowedExtensions {
		if FileExists(basename + "." + v) {
			return v
		}
	}

	return ""
}

// ChatGPT cursed thing
// thus it will probably bites me later
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
