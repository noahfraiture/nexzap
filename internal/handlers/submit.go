package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

func (app *App) SubmitHandler(w http.ResponseWriter, r *http.Request) {
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

	submissionData, err := app.Database.GetRepository().FindSubmissionData(context.Background(), sheetUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve submission data. Error : %s", err), http.StatusInternalServerError)
		return
	}

	output, status, err := app.ExerciseService.RunTest(submissionData, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to run test. Error: %s", err), http.StatusInternalServerError)
	}

	// Respond with JSON containing the output and status code
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Output     string `json:"output"`
		StatusCode int    `json:"statusCode"`
	}{
		Output:     sanitize(output),
		StatusCode: int(status.StatusCode),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// TODO : move that shit
// Define the regex pattern globally for efficiency
var sanitizeRegex = regexp.MustCompile(`[^\x20-\x7E\n\t]`)

// Function to sanitize the output
func sanitize(output string) string {
	return sanitizeRegex.ReplaceAllString(output, "")
}
