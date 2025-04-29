package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (app *App) ImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	fileHeaders, ok := r.MultipartForm.File["tutorial_zip"]
	if !ok || len(fileHeaders) == 0 {
		http.Error(w, "No file uploaded with key 'tutorial_zip'", http.StatusBadRequest)
		return
	}

	fileHeader := fileHeaders[0]
	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Unable to open uploaded file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a temporary file for the zip
	tempFile, err := os.CreateTemp("", "tutorial-*.zip")
	if err != nil {
		http.Error(w, "Unable to create temporary file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // Clean up

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "tutorial-")
	if err != nil {
		http.Error(w, "Unable to create temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Extract the zip file
	err = extractZip(tempFile.Name(), tempDir)
	if err != nil {
		http.Error(w, "Unable to extract zip file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Read all entries in the extracted directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		http.Error(w, "Unable to read extracted directory", http.StatusInternalServerError)
		return
	}

	// Process each directory as a tutorial
	sucess := []string{}
	errs := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			tutorialPath := filepath.Join(tempDir, entry.Name())
			err = app.ImportService.ImportTutorialFromDir(tutorialPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", entry.Name(), err))
			} else {
				sucess = append(sucess, entry.Name())
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	response := "Tutorials uploaded successfully:\n"
	if len(sucess) > 0 {
		response += "Success:\n"
		for _, s := range sucess {
			response += fmt.Sprintf("- %s\n", s)
		}
	}
	if len(errs) > 0 {
		response += "Errors:\n"
		for _, e := range errs {
			response += fmt.Sprintf("- %s\n", e)
		}
	}
	w.Write([]byte(response))
}

// extractZip extracts a zip file to the specified destination, preventing ZipSlip.
func extractZip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Prevent ZipSlip by ensuring path stays within dest
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			rc, err := f.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
