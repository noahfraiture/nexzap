package services

import (
	"os"
	"path/filepath"
	"zapbyte/internal/models"

	"github.com/BurntSushi/toml"
)

func GetTutorials() (*[]models.Tutorial, error) {
	dir := "tutorials"
	tutorialsFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	tutorials := []models.Tutorial{}
	for _, file := range tutorialsFiles {
		if file.Type().IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var tutorial models.Tutorial
		_, err = toml.Decode(string(content), &tutorial)
		if err != nil {
			return nil, err
		}
		tutorials = append(tutorials, tutorial)
	}

	return &tutorials, nil
}
