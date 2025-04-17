package services

import (
	"bytes"
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"

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

type Tutorial generated.FindLastTutorialRow

// LastTutorial get the last tutorial and parse the content of sheets and tests
func LastTutorial() (*Tutorial, error) {
	row, err := db.GetRepository().FindLastTutorial(context.Background())
	if err != nil {
		return nil, err
	}
	tutorial := Tutorial(row)
	for i, sheet := range tutorial.GuideContents {
		newSheet, err := markdownToHtml(sheet)
		if err != nil {
			return nil, err
		}
		tutorial.GuideContents[i] = newSheet
	}
	for i, test := range tutorial.TestContents {
		newTest, err := markdownToHtml(test)
		if err != nil {
			return nil, err
		}
		tutorial.TestContents[i] = newTest
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
