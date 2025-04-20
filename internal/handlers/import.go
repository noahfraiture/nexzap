package handlers

import (
	"encoding/json"
	"net/http"
	"nexzap/internal/services"
)

// ImportData defines the expected JSON structure for ImportHandler
type ImportData struct {
	Tutorial services.InsertTutorialModelInsert `json:"tutorial"`
	Files    []services.FilesPerSheet           `json:"files"`
}

func ImportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var importData ImportData
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&importData); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		if err := services.InsertTutorialAndFiles(importData.Tutorial, importData.Files); err != nil {
			http.Error(w, "Failed to import tutorial: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Tutorial imported successfully"))
	}
}

func RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := services.RefreshTutorials(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Tutorials refreshed successfully"))
	}
}
