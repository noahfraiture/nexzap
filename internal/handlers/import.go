package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexzap/internal/services"
	"os"
	"path/filepath"
)

// FilesPerSheet defines file data for a single sheet
type FilesPerSheet struct {
	Names    []string
	Contents []string
}

// insertTutorialAndFiles inserts a tutorial and its associated files per sheet
func insertTutorialAndFiles(tutorial services.InsertTutorialModelInsert, filesPerSheet []FilesPerSheet) error {
	sheetsID, err := services.InsertTutorial(tutorial)
	if err != nil {
		return err
	}
	if len(sheetsID) != len(filesPerSheet) {
		return fmt.Errorf("number of sheets (%d) does not match number of files data (%d)", len(sheetsID), len(filesPerSheet))
	}
	for i, sheetID := range sheetsID {
		fileInsert := services.InsertFilesModelInsert{
			Names:    filesPerSheet[i].Names,
			Contents: filesPerSheet[i].Contents,
			SheetID:  sheetID,
		}
		if err := services.InsertFile(fileInsert); err != nil {
			return err
		}
	}
	return nil
}

// ImportData defines the expected JSON structure for ImportHandler
type ImportData struct {
	Tutorial services.InsertTutorialModelInsert `json:"tutorial"`
	Files    []FilesPerSheet                    `json:"files"`
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

		if err := insertTutorialAndFiles(importData.Tutorial, importData.Files); err != nil {
			http.Error(w, "Failed to import tutorial: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Tutorial imported successfully"))
	}
}

func RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tutorials, err := os.ReadDir("tutorials/")
		if err != nil {
			http.Error(w, "Failed to read tutorials directory", http.StatusInternalServerError)
			return
		}

		for _, tutorialDir := range tutorials {
			if !tutorialDir.IsDir() {
				http.Error(w, "Invalid entry in tutorials directory: "+tutorialDir.Name(), http.StatusInternalServerError)
				return
			}
			path := filepath.Join("tutorials", tutorialDir.Name())
			meta, sheets, err := services.ReadDirectory(path)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read tutorial directory: %s. Error : %s", path, err), http.StatusInternalServerError)
				return
			}

			// Construct tutorial and files per sheet
			pages := []int32{}
			guides := []string{}
			exercises := []string{}
			images := []string{}
			commands := []string{}
			submission := []string{}
			var filesPerSheet []FilesPerSheet
			for i, sheet := range *sheets {
				pages = append(pages, int32(i+1))
				guides = append(guides, sheet.Guide)
				exercises = append(exercises, sheet.Exercise)
				images = append(images, sheet.Image)
				commands = append(commands, sheet.Command)
				submission = append(submission, sheet.SubmissionFile)

				filesName := []string{}
				filesContent := []string{}
				for _, f := range sheet.Files {
					filesName = append(filesName, f.Name)
					filesContent = append(filesContent, f.Content)
				}
				filesPerSheet = append(filesPerSheet, FilesPerSheet{
					Names:    filesName,
					Contents: filesContent,
				})
			}

			tutorial := services.InsertTutorialModelInsert{
				Title:            meta.Title,
				Highlight:        meta.Highlight,
				CodeEditor:       meta.CodeEditor,
				Version:          int32(meta.Version),
				Pages:            pages,
				GuidesContent:    guides,
				ExercisesContent: exercises,
				DockerImages:     images,
				Commands:         commands,
				SubmissionFile:   submission,
			}

			if err := insertTutorialAndFiles(tutorial, filesPerSheet); err != nil {
				http.Error(w, "Failed to insert tutorial "+meta.Title+": "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Tutorials refreshed successfully"))
	}
}
