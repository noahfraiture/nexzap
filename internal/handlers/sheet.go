package handlers

import (
	"fmt"
	"log"
	"net/http"
	"nexzap/internal/models"
	"nexzap/internal/services"
	"nexzap/templates/pages"
	"strconv"
)

// TODO : cache the tutorial
func (app *App) SheetHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	pageIndex := 1
	if pageParam != "" {
		var convErr error
		pageIndex, convErr = strconv.Atoi(pageParam)
		if convErr != nil || pageIndex < 1 {
			pageIndex = 1
		}
	}

	tutorialId := r.URL.Query().Get("tutorial")

	var sheet models.SheetTempl
	if tutorialId != "" {
		tutorial, err := app.SheetService.SpecificTutorialPage(tutorialId, pageIndex)
		if err != nil {
			fmt.Println(err)
			tutorial = &services.FindSpecificTutorialSheetModelSelect{Title: "Error"}
		}
		sheet = models.NewSheetTempl(
			tutorial.SheetID.String(),
			tutorial.TutorialID.String(),
			tutorial.Title,
			tutorial.CodeEditor,
			tutorial.GuideContent,
			tutorial.ExerciseContent,
			tutorial.SubmissionContent,
			pageIndex,
			int(tutorial.TotalPages),
			false,
		)
	} else {
		tutorial, err := app.SheetService.LastTutorialPage(pageIndex)
		if err != nil {
			tutorial = &services.FindLastTutorialSheetModelSelect{Title: "Error"}
		}
		sheet = models.NewSheetTempl(
			tutorial.SheetID.String(),
			tutorial.TutorialID.String(),
			tutorial.Title,
			tutorial.CodeEditor,
			tutorial.GuideContent,
			tutorial.ExerciseContent,
			tutorial.SubmissionContent,
			pageIndex,
			int(tutorial.TotalPages),
			true,
		)
	}

	var tutorialsTempl []models.ListTutorialTempl
	if !isFromHtmx(r) {
		tutorials, err := app.HistoryService.ListTutorials()
		if err == nil {
			tutorialsTempl = make([]models.ListTutorialTempl, len(tutorials))
			for i, tuto := range tutorials {
				tutorialsTempl[i] = models.NewListTutorial(tuto.ID.String(), tuto.Title)
			}
		} else {
			log.Println(err)
		}
	}

	// Set headers for sheet ID and page number
	w.Header().Set("X-Sheet-ID", sheet.Id)
	err := pages.NextContent(
		isFromHtmx(r),
		sheet,
		tutorialsTempl,
	).Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}

}
