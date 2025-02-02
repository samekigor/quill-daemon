package utils

import (
	"os"
)

func CreateDir(path string) error {
	if path == "" {
		ErrorLogger.Printf("Path is empty")
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		ErrorLogger.Printf("Failed to create directory: %s, error: %v", path, err)
		return err
	}
	return nil
}
