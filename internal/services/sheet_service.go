package services

import (
	"bytes"
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func InitMarkdown() {
	md = goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

type FindTutorialModelSelect = generated.FindLastTutorialFirstSheetRow

// LastTutorialFirstPage get the last tutorial and parse the content of sheets and tests
func LastTutorialFirstPage() (*FindTutorialModelSelect, error) {
	tutorial, err := db.GetRepository().FindLastTutorialFirstSheet(context.Background())
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent, err = markdownToHtml(tutorial.GuideContent)
	if err != nil {
		return nil, err
	}
	tutorial.GuideContent, err = markdownToHtml(tutorial.GuideContent)
	if err != nil {
		return nil, err
	}
	return &tutorial, nil
}

func markdownToHtml(content string) (string, error) {
	var buf bytes.Buffer
	err := md.Convert([]byte(content), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type InsertTutorialModelInsert = generated.InsertTutorialParams
type InsertFilesModelInsert = generated.InsertFilesParams

func InsertTutorial(args InsertTutorialModelInsert) ([]uuid.UUID, error) {
	return db.GetRepository().InsertTutorial(context.Background(), generated.InsertTutorialParams(args))
}

func InsertFile(args InsertFilesModelInsert) error {
	return db.GetRepository().InsertFiles(context.Background(), generated.InsertFilesParams(args))
}
