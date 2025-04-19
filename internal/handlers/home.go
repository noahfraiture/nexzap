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

		sheet := models.SheetTempl{
			SheetContent:    tutorial.GuideContent,
			ExerciseContent: tutorial.ExerciseContent,
			NbPage:          1,
			MaxPage:         int(tutorial.TotalPages),
		}
		pages.Home(
			isFromHtmx(r),
			tutorial.Title,
			sheet,
		).Render(r.Context(), w)
	}
}
