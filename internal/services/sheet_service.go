package services

import (
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"

	"github.com/google/uuid"
)

var md *MarkdownParser

func InitMarkdown() {
	md = NewMarkdownParser()
}

type FindTutorialFirstSheetModelSelect = generated.FindLastTutorialFirstSheetRow
type FindTutorialSheetModelSelect = generated.FindLastTutorialSheetRow

// LastTutorialFirstPage get the last tutorial and parse the content of sheets and tests
func LastTutorialFirstPage() (*FindTutorialFirstSheetModelSelect, error) {
	tutorial, err := db.GetRepository().FindLastTutorialFirstSheet(context.Background())
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent = markdownToHtml(tutorial.GuideContent)
	tutorial.ExerciseContent = markdownToHtml(tutorial.ExerciseContent)
	return &tutorial, nil
}

func LastTutorialPage(page int) (*FindTutorialSheetModelSelect, error) {
	tutorial, err := db.GetRepository().FindLastTutorialSheet(context.Background(), int32(page))
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent = markdownToHtml(tutorial.GuideContent)
	tutorial.ExerciseContent = markdownToHtml(tutorial.ExerciseContent)
	return &tutorial, nil
}

func markdownToHtml(content string) string {
	return md.ParseMarkdown(content)
}

type InsertTutorialModelInsert = generated.InsertTutorialParams
type InsertFilesModelInsert = generated.InsertFilesParams

func InsertTutorial(args InsertTutorialModelInsert) ([]uuid.UUID, error) {
	return db.GetRepository().InsertTutorial(context.Background(), generated.InsertTutorialParams(args))
}

func InsertFile(args InsertFilesModelInsert) error {
	return db.GetRepository().InsertFiles(context.Background(), generated.InsertFilesParams(args))
}
