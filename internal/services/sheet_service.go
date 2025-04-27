package services

import (
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"
)

// SheetService handles operations related to tutorial sheets
type SheetService struct {
	db       *db.Database
	markdown *MarkdownParser
}

// NewSheetService creates a new SheetService with the given database
func NewSheetService(database *db.Database) *SheetService {
	return &SheetService{
		db:       database,
		markdown: NewMarkdownParser(),
	}
}

// FindTutorialSheetModelSelect is an alias for the generated type
type FindTutorialSheetModelSelect = generated.FindLastTutorialSheetRow

// LastTutorialFirstPage gets the last tutorial for page 1 and parses content
func (s *SheetService) LastTutorialFirstPage() (*FindTutorialSheetModelSelect, error) {
	return s.LastTutorialPage(1)
}

// LastTutorialPage gets the last tutorial for the specified page and parses content
func (s *SheetService) LastTutorialPage(page int) (*FindTutorialSheetModelSelect, error) {
	tutorial, err := s.db.GetRepository().FindLastTutorialSheet(context.Background(), int32(page))
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent = s.markdown.ParseMarkdown(tutorial.GuideContent)
	tutorial.ExerciseContent = s.markdown.ParseMarkdown(tutorial.ExerciseContent)
	return &tutorial, nil
}
