package handlers

import (
	"context"
	"fmt"
	"net/http"
	"nexzap/internal/db"
	"nexzap/internal/services"

	"github.com/google/uuid"
)

func SubmitHandler(s *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		payload := r.FormValue("payload")
		if payload == "" {
			http.Error(w, "No payload provided", http.StatusBadRequest)
			return
		}
		sheetId := r.FormValue("sheet")
		if sheetId == "" {
			http.Error(w, "No sheet id provided", http.StatusBadRequest)
			return
		}
		sheetUUID, err := uuid.Parse(sheetId)
		if err != nil {
			http.Error(w, "Invalid sheet id", http.StatusBadRequest)
			return
		}

		repo := db.GetRepository()
		submissionData, err := repo.FindSubmissionData(context.Background(), sheetUUID)
		if err != nil {
			http.Error(w, "Failed to retrieve submission data", http.StatusInternalServerError)
			return
		}
		fmt.Println(submissionData)

		fmt.Println(submissionData)
		fmt.Println(payload)
		s.RunTest(submissionData, payload)

		// Process the payload as needed
		w.WriteHeader(http.StatusOK)
	}
}
