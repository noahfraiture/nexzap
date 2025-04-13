package services

import (
	"os"
	"zapbyte/internal/models"
	"zapbyte/internal/services/container"
)

func GetTutorials() (*[]models.Tutorial, error) { return nil, nil }

func getTutorial(tutorialMeta container.Tutorial) (*models.Tutorial, error) {
	tutorial := models.Tutorial{
		Language: tutorialMeta.Language,
		Sheets:   []models.Sheet{},
	}
	files, err := os.ReadDir()
	return nil, nil
}
