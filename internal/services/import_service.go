package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"

	"github.com/BurntSushi/toml"
)

type ImportService struct {
	numberRegex *regexp.Regexp
	db          *db.Database
}

func NewImportService(db *db.Database) *ImportService {
	return &ImportService{
		numberRegex: regexp.MustCompile(`^\d+`),
		db:          db,
	}
}

// RefreshTutorials reads all tutorial directories in "tutorials/" and inserts them into the database.
func (s *ImportService) RefreshTutorials() error {
	if os.Getenv("ENV") != "dev" {
		return fmt.Errorf("RefreshTutorials can only be run in development environment")
	}
	tutorials, err := os.ReadDir(os.Getenv("TUTORIALS_PATH"))
	if err != nil {
		return fmt.Errorf("Failed to read tutorials directory: %v", err)
	}

	for _, tutorialDir := range tutorials {
		if !tutorialDir.IsDir() {
			return fmt.Errorf("Invalid entry in tutorials directory: %s", tutorialDir.Name())
		}
		path := filepath.Join("tutorials", tutorialDir.Name())
		if err := s.ImportTutorialFromDir(path); err != nil {
			return err
		}
	}

	return nil
}

// ImportTutorialFromDir reads a single tutorial directory and inserts it into the database.
func (s *ImportService) ImportTutorialFromDir(path string) error {
	meta, sheets, err := s.readDirectory(path)
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
	correctionContent := []string{}
	var filesPerSheet []FilesPerSheet
	for i, sheet := range *sheets {
		pages = append(pages, int32(i+1))
		guides = append(guides, sheet.guide)
		exercises = append(exercises, sheet.exercise)
		images = append(images, sheet.Image)
		commands = append(commands, sheet.Command)
		submissionName = append(submissionName, sheet.SubmissionName)
		submissionContent = append(submissionContent, sheet.submissionContent)
		correctionContent = append(correctionContent, sheet.correctionContent)

		filesName := []string{}
		filesContent := []string{}
		for _, f := range sheet.files {
			filesName = append(filesName, f.Name)
			filesContent = append(filesContent, f.Content)
		}
		filesPerSheet = append(filesPerSheet, FilesPerSheet{
			Names:    filesName,
			Contents: filesContent,
		})
	}

	tutorial := generated.InsertTutorialParams{
		Title:              meta.Title,
		Highlight:          meta.Highlight,
		CodeEditor:         meta.CodeEditor,
		Version:            int32(meta.Version),
		Unlock:             meta.UnlockTime,
		Pages:              pages,
		GuidesContent:      guides,
		ExercisesContent:   exercises,
		DockerImages:       images,
		Commands:           commands,
		SubmissionsName:    submissionName,
		SubmissionsContent: submissionContent,
		CorrectionContent:  correctionContent,
	}

	if err := s.insertTutorialAndFiles(tutorial, filesPerSheet); err != nil {
		return fmt.Errorf("Failed to insert tutorial %s: %v", meta.Title, err)
	}

	return nil
}

// insertTutorialAndFiles inserts a tutorial and its associated files per sheet
func (s *ImportService) insertTutorialAndFiles(
	tutorial generated.InsertTutorialParams,
	filesPerSheet []FilesPerSheet,
) error {
	sheetsID, err := s.db.GetRepository().InsertTutorial(context.Background(), tutorial)
	if err != nil {
		return err
	}
	if len(sheetsID) != len(filesPerSheet) {
		return fmt.Errorf(
			"number of sheets (%d) does not match number of files data (%d)",
			len(sheetsID),
			len(filesPerSheet),
		)
	}
	for i, sheetID := range sheetsID {
		fileInsert := generated.InsertFilesParams{
			Names:    filesPerSheet[i].Names,
			Contents: filesPerSheet[i].Contents,
			SheetID:  sheetID,
		}
		if err := s.db.GetRepository().InsertFiles(context.Background(), fileInsert); err != nil {
			return err
		}
	}
	return nil
}

// tutorialMeta holds metadata for a programming language tutorial.
type tutorialMeta struct {
	Title      string    `toml:"title"`
	Highlight  string    `toml:"highlight"`
	CodeEditor string    `toml:"codeEditor"`
	Version    int       `toml:"version"`
	UnlockTime time.Time `toml:"unlock"`
}

// file represents a file with correction content for a tutorial sheet.
type file struct {
	Name    string
	Content string
}

// toml key must be exported
type sheet struct {
	guide             string
	exercise          string
	submissionContent string
	correctionContent string
	SubmissionName    string `toml:"submission"`
	Image             string `toml:"image"`
	Command           string `toml:"command"`
	files             []file
}

// readDirectory reads a tutorial directory, returning metadata and sheets.
// Errors if directory unreadable or files missing.
func (s *ImportService) readDirectory(path string) (*tutorialMeta, *[]sheet, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	meta, err := s.extractMeta(dir, path)
	if err != nil {
		return nil, nil, err
	}

	// Guides
	guides := []os.DirEntry{}
	for _, guide := range dir {
		// Check if guide name starts with a number using regex
		if s.numberRegex.MatchString(guide.Name()) {
			guides = append(guides, guide)
		}
	}
	sort.Slice(guides, func(i, j int) bool { return guides[i].Name() < guides[j].Name() })
	sheets := []sheet{}
	for _, guide := range guides {
		sheet, err := s.readGuide(guide, path)
		if err != nil {
			return nil, nil, err
		}
		sheets = append(sheets, sheet)
	}
	return meta, &sheets, nil
}

// extractMeta extracts metadata from meta.toml in the sheet directory.
// Errors if file missing, unreadable, or fields unset.
func (s *ImportService) extractMeta(dir []os.DirEntry, path string) (*tutorialMeta, error) {
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
func (s *ImportService) readGuide(dir os.DirEntry, basePath string) (sheet, error) {
	dirPath := filepath.Join(basePath, dir.Name())

	// Read metadata from meta.toml in the guide directory
	var sheetMeta sheet
	metaPath := filepath.Join(dirPath, "meta.toml")
	metaContent, err := os.ReadFile(metaPath)
	if err != nil {
		return sheet{}, err
	}
	err = toml.Unmarshal(metaContent, &sheetMeta)
	if err != nil {
		return sheet{}, err
	}

	var correctionFiles []file

	paths, err := s.findFiles(dirPath, sheetMeta.SubmissionName)
	if err != nil {
		return sheet{}, err
	}

	// Find each file in the correction/ directory that will be run
	correctionFiles, err = s.readCorrectionFiles(dirPath)
	if err != nil {
		return sheet{}, err
	}

	guideContent, err := os.ReadFile(paths.Guide)
	if err != nil {
		return sheet{}, err
	}

	exerciseContent, err := os.ReadFile(paths.Exercise)
	if err != nil {
		return sheet{}, err
	}

	submissionContent, err := os.ReadFile(paths.Submission)
	if err != nil {
		return sheet{}, err
	}

	correctionContent, err := os.ReadFile(paths.Correction)
	if err != nil {
		return sheet{}, err
	}

	sheet := sheet{
		guide:             string(guideContent),
		exercise:          string(exerciseContent),
		Image:             sheetMeta.Image,
		Command:           sheetMeta.Command,
		SubmissionName:    sheetMeta.SubmissionName,
		submissionContent: string(submissionContent),
		correctionContent: string(correctionContent),
		files:             correctionFiles,
	}

	return sheet, nil
}

// FilePaths holds the paths to various files in a tutorial sheet.
type FilePaths struct {
	Guide      string
	Exercise   string
	Submission string
	Correction string
}

// findFiles finds paths to guide.md, exercise.md, submission, and correction files.
func (s *ImportService) findFiles(dirPath string, submissionName string) (FilePaths, error) {
	var paths FilePaths
	paths.Guide = filepath.Join(dirPath, "guide.md")
	paths.Exercise = filepath.Join(dirPath, "exercise.md")
	paths.Submission = filepath.Join(dirPath, filepath.Base(submissionName))
	paths.Correction = filepath.Join(dirPath, "correction", submissionName)

	if _, err := os.Stat(paths.Guide); os.IsNotExist(err) {
		return paths, errors.New("guide.md not found at " + paths.Guide)
	}
	if _, err := os.Stat(paths.Exercise); os.IsNotExist(err) {
		return paths, errors.New("exercise.md not found at " + paths.Exercise)
	}
	if _, err := os.Stat(paths.Submission); os.IsNotExist(err) {
		return paths, errors.New("submission file not found at " + paths.Submission)
	}
	if _, err := os.Stat(paths.Correction); os.IsNotExist(err) {
		return paths, errors.New("correction file not found at " + paths.Correction)
	}

	return paths, nil
}

// readCorrectionFiles reads correction files from a subdirectory.
// Errors if directory unreadable.
func (s *ImportService) readCorrectionFiles(dirPath string) ([]file, error) {
	var correctionFiles []file
	correctionDir := filepath.Join(dirPath, "correction")
	if _, err := os.Stat(correctionDir); os.IsNotExist(err) {
		return correctionFiles, errors.New("correction dir not found at " + correctionDir)
	}

	if correctionDir == "" {
		return nil, errors.New("correction directory not found")
	}
	err := s.readCodeFiles(correctionDir, "", &correctionFiles)
	return correctionFiles, err
}

// readCodeFiles recursively reads code files from a directory and its subdirectories.
// It appends the file information to the provided files slice.
// Errors if directory is unreadable or file operations fail.
func (s *ImportService) readCodeFiles(dirPath, subDir string, files *[]file) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		if entry.IsDir() {
			err = s.readCodeFiles(
				filepath.Join(dirPath, entry.Name()),
				filepath.Join(subDir, entry.Name()),
				files,
			)
			if err != nil {
				return err
			}
		} else {
			content, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				return err
			}
			*files = append(*files, file{
				Name:    filepath.Join(subDir, entry.Name()),
				Content: string(content),
			})
		}
	}
	return nil
}

// FilesPerSheet defines file data for a single sheet
type FilesPerSheet struct {
	Names    []string
	Contents []string
}
