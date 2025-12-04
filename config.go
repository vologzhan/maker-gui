//go:build !dev

package main

import (
	"os"
	"path/filepath"
)

func getSourceDir() (string, error) {
	filePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(filePath), nil
}
