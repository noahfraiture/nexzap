package handlers

import (
	"net/http"
	"nexzap/internal/models"
	"nexzap/internal/services"
	"nexzap/templates/pages"
	"strconv"
)

func HomeHandler() http.HandlerFunc {
	// Computed once
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, err := services.LastTutorialFirstPage()
		if err != nil {
			tutorial = &services.FindTutorialModelSelect{Title: "Error"}
		}

		sheet := models.SheetTempl{
			SheetContent:    tutorial.GuideContent,
			ExerciseContent: tutorial.ExerciseContent,
			NbPage:          0,
			MaxPage:         int(tutorial.TotalPages),
		}
		pages.Home(
			isFromHtmx(r),
			tutorial.Title,
			sheet,
		).Render(r.Context(), w)
	}
}

// TODO : cache the tutorial
// TODO : limit page and return error if not present
func SheetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, err := services.LastTutorialFirstPage()
		if err != nil {
			tutorial = &services.FindTutorialModelSelect{Title: "Error"}
		}
		pageParam := r.URL.Query().Get("page")
		pageIndex := 0
		if pageParam != "" {
			var convErr error
			pageIndex, convErr = strconv.Atoi(pageParam)
			if convErr != nil || pageIndex < 0 || pageIndex >= int(tutorial.TotalPages) {
				pageIndex = 0
			}
		}
		sheet := models.SheetTempl{
			SheetContent:    tutorial.GuideContent,
			ExerciseContent: tutorial.ExerciseContent,
			NbPage:          pageIndex,
			MaxPage:         int(tutorial.TotalPages),
		}
		pages.NextContent(sheet).Render(r.Context(), w)
	}
}
