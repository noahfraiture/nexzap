package handlers

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"zapbyte/internal/models"
	"zapbyte/internal/services"
	"zapbyte/templates/pages"
)

func HomeHandler() http.HandlerFunc {
	// Computed once
	return func(w http.ResponseWriter, r *http.Request) {
		tutorial, _ := getRandomTutorial()
		pages.Home(fromHtmx(r), *tutorial).Render(r.Context(), w)
	}
}

func getRandomTutorial() (*models.Tutorial, error) {

	tutorials, err := services.GetTutorials()
	if err != nil {
		return nil, err
	}
	if len(*tutorials) == 0 {
		return nil, fmt.Errorf("No tutorial found")
	}
	randomIndex := rand.IntN(len(*tutorials))
	return &(*tutorials)[randomIndex], nil
}
