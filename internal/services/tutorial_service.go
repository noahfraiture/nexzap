package services

import (
	"os"
	"regexp"
	"zapbyte/internal/models"
)

var (
	textRegex *regexp.Regexp
	testRegex *regexp.Regexp
)

func GetTutorials() (*[]models.Tutorial, error) { return nil, nil }

func getTutorial(tutorialsDir string, dir os.DirEntry) (*models.Tutorial, error) { return nil, nil }
