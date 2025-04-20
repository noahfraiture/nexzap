package handlers

import (
	"log"
	"net/http"
	"nexzap/internal/models"
	"nexzap/internal/services"
	"nexzap/templates/pages"
)

func HomeHandler() http.HandlerFunc {
	// Computed once
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, err := services.LastTutorialFirstPage()
		if err != nil {
			tutorial = &services.FindTutorialFirstSheetModelSelect{Title: "Error"}
			log.Println(err)
		}

		sheet := models.NewSheetTempl(
			tutorial.ID.String(),
			tutorial.Title,
			tutorial.Highlight,
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
}
