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
		tutorial = &services.FindTutorialSheetModelSelect{Title: "Error"}
		log.Println(err)
	}

	sheet := models.NewSheetTempl(
		tutorial.ID.String(),
		tutorial.Title,
		tutorial.CodeEditor,
		tutorial.GuideContent,
		tutorial.ExerciseContent,
		tutorial.SubmissionContent,
		1,
		int(tutorial.TotalPages),
	)
	pages.Home(
		isFromHtmx(r),
		sheet,
	).Render(r.Context(), w)
}
