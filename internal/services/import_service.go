package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/BurntSushi/toml"
)

// RefreshTutorials read the tutorials in the current directory and insert them
func RefreshTutorials() error {
	tutorials, err := os.ReadDir("tutorials/")
	if err != nil {
		return fmt.Errorf("Failed to read tutorials directory: %v", err)
	}

	for _, tutorialDir := range tutorials {
		if !tutorialDir.IsDir() {
			return fmt.Errorf("Invalid entry in tutorials directory: %s", tutorialDir.Name())
		}
		path := filepath.Join("tutorials", tutorialDir.Name())
		meta, sheets, err := readDirectory(path)
		if err != nil {
			return fmt.Errorf("Failed to read tutorial directory: %s. Error: %v", path, err)
		}

		// Construct tutorial and files per sheet
		pages := []int32{}
		guides := []string{}
		exercises := []string{}
		images := []string{}
		commands := []string{}
		submissionName := []string{}
		submissionContent := []string{}
		var filesPerSheet []FilesPerSheet
		for i, sheet := range *sheets {
			pages = append(pages, int32(i+1))
			guides = append(guides, sheet.Guide)
			exercises = append(exercises, sheet.Exercise)
			images = append(images, sheet.Image)
			commands = append(commands, sheet.Command)
			submissionName = append(submissionName, sheet.SubmissionName)
			submissionContent = append(submissionContent, sheet.SubmissionContent)

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

		tutorial := InsertTutorialModelInsert{
			Title:              meta.Title,
			Highlight:          meta.Highlight,
			CodeEditor:         meta.CodeEditor,
			Version:            int32(meta.Version),
			Pages:              pages,
			GuidesContent:      guides,
			ExercisesContent:   exercises,
			DockerImages:       images,
			Commands:           commands,
			SubmissionsName:    submissionName,
			SubmissionsContent: submissionContent,
		}

		if err := InsertTutorialAndFiles(tutorial, filesPerSheet); err != nil {
			return fmt.Errorf("Failed to insert tutorial %s: %v", meta.Title, err)
		}
	}

	return nil
}

// InsertTutorialAndFiles inserts a tutorial and its associated files per sheet
func InsertTutorialAndFiles(tutorial InsertTutorialModelInsert, filesPerSheet []FilesPerSheet) error {
	sheetsID, err := InsertTutorial(tutorial)
	if err != nil {
		return err
	}
	if len(sheetsID) != len(filesPerSheet) {
		return fmt.Errorf("number of sheets (%d) does not match number of files data (%d)", len(sheetsID), len(filesPerSheet))
	}
	for i, sheetID := range sheetsID {
		fileInsert := InsertFilesModelInsert{
			Names:    filesPerSheet[i].Names,
			Contents: filesPerSheet[i].Contents,
			SheetID:  sheetID,
		}
		if err := InsertFile(fileInsert); err != nil {
			return err
		}
	}
	return nil
}

// tutorialMeta holds metadata for a programming language tutorial.
type tutorialMeta struct {
	Title      string `toml:"title"`
	Highlight  string `toml:"highlight"`
	CodeEditor string `toml:"codeEditor"`
	Version    int    `toml:"version"`
}

// correctionFile represents a file with correction content for a tutorial sheet.
type correctionFile struct {
	Name    string
	Content string
}

// sheet represents a tutorial sheet with guide, exercise, and correction content.
type sheet struct {
	Guide             string
	Exercise          string
	SubmissionContent string
	SubmissionName    string `toml:"submission"`
	Image             string `toml:"image"`
	Command           string `toml:"command"`
	Files             []correctionFile
}

// readDirectory reads a tutorial directory, returning metadata and sheets.
// Errors if directory unreadable or files missing.
func readDirectory(path string) (*tutorialMeta, *[]sheet, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	meta, err := extractMeta(dir, path)
	if err != nil {
		return nil, nil, err
	}

	// Guides
	guides := []os.DirEntry{}
	for _, guide := range dir {
		// Check if guide name starts with a number using regex
		matched, err := regexp.MatchString(`^\d+`, guide.Name())
		if err != nil {
			return nil, nil, err
		}
		if matched {
			guides = append(guides, guide)
		}
	}
	sort.Slice(guides, func(i, j int) bool { return guides[i].Name() < guides[j].Name() })
	sheets := []sheet{}
	for _, guide := range guides {
		fmt.Println(guide)
		fmt.Println(path)
		sheet, err := readGuide(guide, path) // FIXME
		if err != nil {
			return nil, nil, err
		}
		sheets = append(sheets, sheet)
	}
	return meta, &sheets, nil
}

// extractMeta extracts metadata from meta.toml in the directory.
// Errors if file missing, unreadable, or fields unset.
func extractMeta(dir []os.DirEntry, path string) (*tutorialMeta, error) {
	var metaFile *os.DirEntry
	for _, f := range dir {
		if f.Name() == "meta.toml" {
			metaFile = &f
			break
		}
	}

	if metaFile == nil {
		return nil, errors.New("meta.toml file not found in directory")
	}

	metaPath := filepath.Join(path, (*metaFile).Name())
	content, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var meta tutorialMeta
	if err := toml.Unmarshal(content, &meta); err != nil {
		return nil, err
	}

	// Check if all required fields are set
	if meta.Title == "" {
		return nil, errors.New("title field is not set in meta.toml")
	}
	if meta.Version == 0 {
		return nil, errors.New("version field is not set in meta.toml")
	}

	return &meta, nil
}

// readGuide processes a guide directory to create a Sheet.
// Errors if required files missing or unreadable.
func readGuide(dir os.DirEntry, basePath string) (sheet, error) {
	dirPath := filepath.Join(basePath, dir.Name())
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return sheet{}, err
	}

	// Read metadata from meta.toml in the guide directory
	var sheetMeta sheet
	metaPath := filepath.Join(dirPath, "meta.toml")
	metaContent, err := os.ReadFile(metaPath)
	if err == nil {
		if err := toml.Unmarshal(metaContent, &sheetMeta); err != nil {
			return sheet{}, err
		}
	} else {
		// If meta.toml is not found, initialize with empty values
		sheetMeta = sheet{}
	}

	var guideFile, exerciseFile, submissionFile string
	var correctionFiles []correctionFile

	// get the markdown content exercise.md and guide.md
	guideFile, exerciseFile, submissionFile = findFiles(entries, dirPath, sheetMeta.SubmissionName)
	if guideFile == "" || exerciseFile == "" {
		return sheet{}, errors.New("guide.md or exercise.md not found in directory")
	}

	// Find each file in the correction/ that will be run
	correctionFiles, err = readCorrectionFiles(entries, dirPath)
	if err != nil {
		return sheet{}, err
	}

	guideContent, err := os.ReadFile(guideFile)
	if err != nil {
		return sheet{}, err
	}

	exerciseContent, err := os.ReadFile(exerciseFile)
	if err != nil {
		return sheet{}, err
	}

	submissionContent, err := os.ReadFile(submissionFile)
	if err != nil {
		return sheet{}, err
	}

	sheet := sheet{
		Guide:             string(guideContent),
		Exercise:          string(exerciseContent),
		Image:             sheetMeta.Image,
		Command:           sheetMeta.Command,
		SubmissionName:    sheetMeta.SubmissionName,
		SubmissionContent: string(submissionContent),
		Files:             correctionFiles,
	}

	return sheet, nil
}

// findGuideAndExerciseFiles finds paths to guide.md and exercise.md files.
func findFiles(entries []os.DirEntry, dirPath string, submissionName string) (string, string, string) {
	var guideFile, exerciseFile, submissionFile string
	for _, entry := range entries {
		switch entry.Name() {
		case "guide.md":
			guideFile = filepath.Join(dirPath, entry.Name())
		case "exercise.md":
			exerciseFile = filepath.Join(dirPath, entry.Name())
		case submissionName:
			submissionFile = filepath.Join(dirPath, entry.Name())
		}
	}
	return guideFile, exerciseFile, submissionFile
}

// readCorrectionFiles reads correction files from a subdirectory.
// Errors if directory unreadable.
func readCorrectionFiles(entries []os.DirEntry, dirPath string) ([]correctionFile, error) {
	var correctionFiles []correctionFile
	var correctionDir string

	for _, entry := range entries {
		if entry.Name() == "correction" && entry.IsDir() {
			correctionDir = filepath.Join(dirPath, entry.Name())
			break
		}
	}

	if correctionDir == "" {
		return nil, errors.New("correction directory not found")
	}
	corrEntries, err := os.ReadDir(correctionDir)
	if err != nil {
		return nil, err
	}

	for _, corrEntry := range corrEntries {
		if !corrEntry.IsDir() {
			content, err := os.ReadFile(filepath.Join(correctionDir, corrEntry.Name()))
			if err != nil {
				return nil, err
			}
			correctionFiles = append(correctionFiles, correctionFile{
				Name:    corrEntry.Name(),
				Content: string(content),
			})
		}
	}

	return correctionFiles, nil
}

// FilesPerSheet defines file data for a single sheet
type FilesPerSheet struct {
	Names    []string
	Contents []string
}
