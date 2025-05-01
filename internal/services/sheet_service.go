package services

import (
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"
	"regexp"

	"github.com/google/uuid"
)

// SheetService handles operations related to tutorial sheets
type SheetService struct {
	db          *db.Database
	markdown    *MarkdownParser
	sanitizeReg *regexp.Regexp
}

// NewSheetService creates a new SheetService with the given database
func NewSheetService(database *db.Database) *SheetService {
	return &SheetService{
		db:          database,
		markdown:    NewMarkdownParser(),
		sanitizeReg: regexp.MustCompile(`[^\x20-\x7E\n\t]`),
	}
}

// FindLastTutorialSheetModelSelect is an alias for the generated type
type FindLastTutorialSheetModelSelect = generated.FindLastTutorialSheetRow
type FindSpecificTutorialSheetModelSelect = generated.FindSpecificTutorialSheetRow

// LastTutorialFirstPage gets the last tutorial for page 1 and parses content
func (s *SheetService) LastTutorialFirstPage() (*FindLastTutorialSheetModelSelect, error) {
	return s.LastTutorialPage(1)
}

// LastTutorialPage gets the last tutorial for the specified page and parses content
func (s *SheetService) LastTutorialPage(page int) (*FindLastTutorialSheetModelSelect, error) {
	tutorial, err := s.db.GetRepository().FindLastTutorialSheet(context.Background(), int32(page))
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent = s.markdown.ParseMarkdown(tutorial.GuideContent)
	tutorial.ExerciseContent = s.markdown.ParseMarkdown(tutorial.ExerciseContent)
	return &tutorial, nil
}

func (s *SheetService) SpecificTutorialPage(id string, page int) (
	*FindSpecificTutorialSheetModelSelect,
	error,
) {

	tutorialId, err := uuid.Parse(id)
	if err != nil {
		return &generated.FindSpecificTutorialSheetRow{}, err
	}
	tutorial, err := s.db.GetRepository().FindSpecificTutorialSheet(
		context.Background(),
		generated.FindSpecificTutorialSheetParams{
			Page:       int32(page),
			TutorialID: tutorialId,
		},
	)
	tutorial.GuideContent = s.markdown.ParseMarkdown(tutorial.GuideContent)
	tutorial.ExerciseContent = s.markdown.ParseMarkdown(tutorial.ExerciseContent)
	return &tutorial, err
}

func (s *SheetService) Sanitize(content string) string {
	return s.sanitizeReg.ReplaceAllString(content, "")
}
