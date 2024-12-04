package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Warning:
 * AI generated code, it can break or so but I'm too lazy to re-do it properly myself
 **/
func ExtractZip(zipPath string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	tempDir, err := os.MkdirTemp("", "zipextract-*")
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
		destPath := filepath.Join(tempDir, file.Name)
		if !strings.HasPrefix(destPath, filepath.Clean(tempDir)+string(os.PathSeparator)) {
			os.RemoveAll(tempDir)
			return "", errors.New("invalid file path in zip: " + file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				os.RemoveAll(tempDir)
				return "", err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			os.RemoveAll(tempDir)
			return "", err
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()
		if err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}
	}

	return tempDir, nil
}

/**
 * Caution, this does not work with folders
 */
func MoveFileCrossDrive(src string, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	// Delete the source file
	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("could not remove source file: %w", err)
	}

	return nil
}
