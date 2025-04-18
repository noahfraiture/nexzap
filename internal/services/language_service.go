package services

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/BurntSushi/toml"
)

// TutorialMeta holds metadata for a programming language tutorial.
type TutorialMeta struct {
	title       string `toml:"title"`
	highlight   string `toml:"highlight"`
	codeEditor  string `toml:"code_editor"`
	dockerImage string `toml:"docker_image"`
	version     int    `toml:"version"`
}

// CorrectionFile represents a file with correction content for a tutorial sheet.
type CorrectionFile struct {
	name    string
	content string
}

// Sheet represents a tutorial sheet with guide, exercise, and correction content.
type Sheet struct {
	guide      string
	exercise   string
	correction []CorrectionFile
}

// readDirectory reads a tutorial directory, returning metadata and sheets.
// Errors if directory unreadable or files missing.
func readDirectory(path string) (*TutorialMeta, *[]Sheet, error) {
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
	sheets := []Sheet{}
	for _, guide := range guides {
		sheet, err := readGuide(guide)
		if err != nil {
			return nil, nil, err
		}
		sheets = append(sheets, sheet)
	}
	return meta, &sheets, nil
}

// extractMeta extracts metadata from meta.toml in the directory.
// Errors if file missing, unreadable, or fields unset.
func extractMeta(dir []os.DirEntry, path string) (*TutorialMeta, error) {
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

	var meta TutorialMeta
	if err := toml.Unmarshal(content, &meta); err != nil {
		return nil, err
	}

	// Check if all required fields are set
	if meta.title == "" {
		return nil, errors.New("title field is not set in meta.toml")
	}
	if meta.dockerImage == "" {
		return nil, errors.New("docker_image field is not set in meta.toml")
	}
	if meta.version == 0 {
		return nil, errors.New("version field is not set in meta.toml")
	}

	return &meta, nil
}

// readGuide processes a guide directory to create a Sheet.
// Errors if required files missing or unreadable.
func readGuide(dir os.DirEntry) (Sheet, error) {
	dirPath := dir.Name()
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return Sheet{}, err
	}

	var guideFile, exerciseFile string
	var correctionFiles []CorrectionFile

	guideFile, exerciseFile = findGuideAndExerciseFiles(entries, dirPath)
	if guideFile == "" || exerciseFile == "" {
		return Sheet{}, errors.New("guide.md or exercise.md not found in directory")
	}

	correctionFiles, err = readCorrectionFiles(entries, dirPath)
	if err != nil {
		return Sheet{}, err
	}

	guideContent, err := os.ReadFile(guideFile)
	if err != nil {
		return Sheet{}, err
	}

	exerciseContent, err := os.ReadFile(exerciseFile)
	if err != nil {
		return Sheet{}, err
	}

	sheet := Sheet{
		guide:      string(guideContent),
		exercise:   string(exerciseContent),
		correction: correctionFiles,
	}

	return sheet, nil
}

// findGuideAndExerciseFiles finds paths to guide.md and exercise.md files.
func findGuideAndExerciseFiles(entries []os.DirEntry, dirPath string) (string, string) {
	var guideFile, exerciseFile string
	for _, entry := range entries {
		switch entry.Name() {
		case "guide.md":
			guideFile = filepath.Join(dirPath, entry.Name())
		case "exercise.md":
			exerciseFile = filepath.Join(dirPath, entry.Name())
		}
	}
	return guideFile, exerciseFile
}

// readCorrectionFiles reads correction files from a subdirectory.
// Errors if directory unreadable.
func readCorrectionFiles(entries []os.DirEntry, dirPath string) ([]CorrectionFile, error) {
	var correctionFiles []CorrectionFile
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
			correctionFiles = append(correctionFiles, CorrectionFile{
				name:    corrEntry.Name(),
				content: string(content),
			})
		}
	}

	return correctionFiles, nil
}
