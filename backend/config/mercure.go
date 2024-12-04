package config

import (
	"crypto/rand"
	"os"
	"path/filepath"
)

func loadMercureKey(rootPath string, filename string) ([]byte, error) {
	keyPath := filepath.Join(rootPath, filename+".key")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		mercureKey, err := generateMercureKey()
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(keyPath, mercureKey, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func generateMercureKey() ([]byte, error) {
	key := make([]byte, 512)
	_, err := rand.Read(key)

	return key, err
}
