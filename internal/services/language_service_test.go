package services

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestRunTestGo(t *testing.T) {
	// Initialize the container service
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}

	// Define the directory containing the Go test files
	goDir := "languages/go"

	// Read all files from the directory
	files, err := os.ReadDir(goDir)
	if err != nil {
		t.Fatalf("Failed to read Go directory: %v", err)
	}

	// Prepare a slice of file paths for the test
	filePaths := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(goDir, file.Name()))
		}
	}

	// Run the test for Go language
	var wg sync.WaitGroup
	wg.Add(5)
	for i := range 5 {
		go func(runNum int) {
			defer wg.Done()
			output, err := svc.RunTest(GO, filePaths)
			if err != nil {
				t.Errorf("Run %d: Failed to run Go test: %v", runNum, err)
				return
			}
			// Optionally, print the output for debugging purposes
			t.Logf("Run %d - Go Test Output:\n%s", runNum, output)
			// Check if output is not empty
			if output == "" {
				t.Errorf("Run %d: Expected non-empty output from Go test, but got empty string", runNum)
			}
		}(i + 1)
	}
	wg.Wait()
}
