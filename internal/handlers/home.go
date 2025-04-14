package handlers

import (
	"net/http"
	"strconv"
	"zapbyte/internal/services"
	"zapbyte/templates/pages"
)

func HomeHandler() http.HandlerFunc {
	// Computed once
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, err := services.LastTutorial()
		if err != nil {
			tutorial = &services.Tutorial{LanguageName: "Error"}
		}
		pages.Home(
			isFromHtmx(r),
			tutorial.LanguageName,
			tutorial.SheetContents[0],
			tutorial.TestContents[0],
		).Render(r.Context(), w)
	}
}

// TODO : cache the tutorial
// TODO : limit page and return error if not present
func PageHandler() http.HandlerFunc {
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
			if convErr != nil || pageIndex < 0 || pageIndex >= len(tutorial.SheetContents) {
				pageIndex = 0
			}
		}
		pages.PageContent(
			tutorial.SheetContents[pageIndex],
			tutorial.TestContents[pageIndex],
		).Render(r.Context(), w)
	}
}
