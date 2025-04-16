package handlers

import (
	"encoding/json"
	"net/http"
	"zapbyte/internal/services"
)

// InsertTutorialHandler handles the HTTP request to insert a new complete tutorial
func InsertTutorialHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params services.TutorialParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Basic validation
		if params.LanguageName == "" {
			http.Error(w, "Language name is required", http.StatusBadRequest)
			return
		}
		if len(params.DockerImages) == 0 {
			http.Error(w, "Docker image is required", http.StatusBadRequest)
			return
		}
		if len(params.TestContents) == 0 {
			http.Error(w, "At least one test content is required", http.StatusBadRequest)
			return
		}
		if len(params.GuideContents) == 0 {
			http.Error(w, "At least one sheet content is required", http.StatusBadRequest)
			return
		}

		tutorialID, err := services.InsertCompleteTutorial(r.Context(), params)
		if err != nil {
			http.Error(w, "Failed to insert tutorial: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"tutorial_id": tutorialID})
	}
}
