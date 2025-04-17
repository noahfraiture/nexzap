package services

import (
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"
)

type TutorialParams generated.InsertCompleteTutorialParams

func InsertCompleteTutorial(ctx context.Context, params TutorialParams) (string, error) {
	tutorial, err := db.GetRepository().InsertCompleteTutorial(ctx, generated.InsertCompleteTutorialParams{
		LanguageName:  params.LanguageName,
		TestContents:  params.TestContents,
		DockerImages:  params.DockerImages,
		GuideContents: params.GuideContents,
	})
	if err != nil {
		return "", err
	}
	return tutorial.String(), nil
}
