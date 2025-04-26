package handlers

import (
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

		tutorial, err := services.LastTutorialPage(pageIndex)
		if err != nil {
			tutorial = &services.FindTutorialSheetModelSelect{Title: "Error"}
		}
		sheet := models.NewSheetTempl(
			tutorial.ID.String(),
			tutorial.Title,
			tutorial.Highlight,
			tutorial.CodeEditor,
			tutorial.GuideContent,
			tutorial.ExerciseContent,
			tutorial.SubmissionContent,
			pageIndex,
			int(tutorial.TotalPages),
		)

		// Set headers for sheet ID and page number
		w.Header().Set("X-Sheet-ID", sheet.Id)
		w.Header().Set("X-Page", strconv.Itoa(pageIndex))
		pages.NextContent(sheet).Render(r.Context(), w)
	}
}
