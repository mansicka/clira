package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mansicka/rtpms/internal/event"
)

type Storage struct {
	RootDir string
}

var instance *Storage

func NewStorage(rootDir string) (*Storage, error) {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to access working directory: %w", err)
	}

	instance = &Storage{
		RootDir: rootDir,
	}
	return instance, nil
}

func GetStorage() (*Storage, error) {
	if instance == nil {
		return nil, fmt.Errorf("storage has not been initialized")
	}
	return instance, nil
}

func (s *Storage) ReadFile(filePath string) ([]byte, error) {
	fullPath := filepath.Join(s.RootDir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", fullPath, err)
	}
	return data, nil
}

func (s *Storage) WriteFile(filePath string, data []byte) error {
	fullPath := filepath.Join(s.RootDir, filePath)
	err := os.WriteFile(fullPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %w", fullPath, err)
	}

	event.TriggerFileSaveEvent(fmt.Sprintf("File saved: %s", filePath))

	return nil
}

func (s *Storage) FileExists(filePath string) bool {
	fullPath := filepath.Join(s.RootDir, filePath)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

func (s *Storage) InitializeDirectoryStructure() error {
	dirs := []string{
		"users",
		"projects",
		"clients",
		"configuration",
	}

	for _, dir := range dirs {
		path := fmt.Sprintf("%s/%s", s.RootDir, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
			fmt.Printf("Created directory: %s\n", path)
		} else {
			fmt.Printf("Directory already exists: %s\n", path)
		}
	}

	return nil
}

func (s *Storage) ReadDir(dirPath string) ([]os.DirEntry, error) {
	fullPath := filepath.Join(s.RootDir, dirPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", fullPath, err)
	}
	return entries, nil
}

func (s *Storage) CreateDir(dirPath string) error {
	fullPath := filepath.Join(s.RootDir, dirPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", fullPath, err)
		}
	}
	return nil
}
