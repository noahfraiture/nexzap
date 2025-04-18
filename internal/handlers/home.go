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
		tutorial, err := services.LastTutorial()
		if err != nil {
			tutorial = &services.Tutorial{LanguageName: "Error"}
		}

		sheet := models.SheetTempl{
			SheetContent: tutorial.GuideContents[0],
			TestContent:  tutorial.TestContents[0],
			NbPage:       0,
			MaxPage:      len(tutorial.GuideContents),
		}
		pages.Home(
			isFromHtmx(r),
			tutorial.LanguageName,
			sheet,
		).Render(r.Context(), w)
	}
}

// TODO : cache the tutorial
// TODO : limit page and return error if not present
func SheetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, err := services.LastTutorial()
		if err != nil {
			tutorial = &services.Tutorial{LanguageName: "Error"}
		}
		pageParam := r.URL.Query().Get("page")
		pageIndex := 0
		if pageParam != "" {
			var convErr error
			pageIndex, convErr = strconv.Atoi(pageParam)
			if convErr != nil || pageIndex < 0 || pageIndex >= len(tutorial.GuideContents) {
				pageIndex = 0
			}
		}
		sheet := models.SheetTempl{
			SheetContent: tutorial.GuideContents[pageIndex],
			TestContent:  tutorial.TestContents[pageIndex],
			NbPage:       pageIndex,
			MaxPage:      len(tutorial.GuideContents),
		}
		pages.NextContent(sheet).Render(r.Context(), w)
	}
}
