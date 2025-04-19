package handlers

import (
	"fmt"
	"net/http"
	"nexzap/internal/models"
	"nexzap/internal/services"
	"nexzap/templates/pages"
	"strconv"
)

// TODO : cache the tutorial
// TODO : limit page and return error if not present
func SheetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageParam := r.URL.Query().Get("page")
		pageIndex := 1
		if pageParam != "" {
			var convErr error
			pageIndex, convErr = strconv.Atoi(pageParam)
			if convErr != nil || pageIndex < 1 {
				pageIndex = 1
			}
		}

		fmt.Println(pageIndex)
		tutorial, err := services.LastTutorialPage(pageIndex)
		if err != nil {
			tutorial = &services.FindTutorialSheetModelSelect{Title: "Error"}
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
