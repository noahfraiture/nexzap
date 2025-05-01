package handlers

import (
	"log"
	"net/http"
	"nexzap/internal/models"
	"nexzap/internal/services"
	"nexzap/templates/pages"
)

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	tutorial, err := app.SheetService.LastTutorialFirstPage()
	if err != nil {
		tutorial = &services.FindLastTutorialSheetModelSelect{Title: "Error"}
		log.Println(err)
	}

	sheet := models.NewSheetTempl(
		tutorial.SheetID.String(),
		tutorial.TutorialID.String(),
		tutorial.Title,
		tutorial.CodeEditor,
		tutorial.GuideContent,
		tutorial.ExerciseContent,
		tutorial.SubmissionContent,
		1,
		int(tutorial.TotalPages),
		true,
	)
	tutorials, err := app.HistoryService.ListTutorials()
	var tutorialsTempl []models.ListTutorialTempl
	if err == nil {
		tutorialsTempl = make([]models.ListTutorialTempl, len(tutorials))
		for i, tuto := range tutorials {
			tutorialsTempl[i] = models.NewListTutorial(tuto.ID.String(), tuto.Title)
		}
	} else {
		log.Println(err)
	}

	err = pages.Home(
		isFromHtmx(r),
		sheet,
		tutorialsTempl,
	).Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}
