package services_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"
	services "nexzap/internal/services"
)

type TestService struct {
	exercise *services.ExerciseService
	db       *db.Database
}

func setupService(t *testing.T) *TestService {
	db, err := db.NewDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	svc, err := services.NewExerciseService()
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	return &TestService{
		exercise: svc,
		db:       db,
	}
}

// runSingleTest performs a single RunTest execution and assertions for a given row.
func (s *TestService) runSingleTest(t *testing.T, row generated.FindCorrectionSheetRow) {
	output, status, err := s.exercise.RunTest(services.Correction{
		DockerImage:    row.DockerImage,
		Command:        row.Command,
		SubmissionName: row.SubmissionName,
		FilesName:      row.FilesName,
		FilesContent:   row.FilesContent,
	}, row.CorrectionContent)
	if err != nil {
		t.Errorf("RunTest failed: %v", err)
		return
	}
	if status.Error != nil {
		t.Error(status.Error)
		return
	}
	if status.StatusCode != 0 {
		t.Errorf("Test failed with error code %d and output %s", status.StatusCode, output)
		return
	}
	if output == "" {
		t.Errorf("Expected non-empty output, got empty")
	}
}

func TestRunTest_Sequential(t *testing.T) {
	svc := setupService(t)
	defer func() {
		if err := svc.exercise.Cleanup(); err != nil {
			t.Error(err)
		}
	}()

	dbRepo := svc.db.GetRepository()
	rows, err := dbRepo.FindCorrectionSheet(context.Background())
	if err != nil {
		t.Fatalf("Failed to fetch correction sheets: %v", err)
	}

	if len(rows) == 0 {
		t.Skip("No data in database; skipping test")
	}

	for _, row := range rows[:10] {
		svc.runSingleTest(t, row)
	}
}

func TestRunTest_AllData(t *testing.T) {
	svc := setupService(t)
	defer func() {
		if err := svc.exercise.Cleanup(); err != nil {
			t.Error(err)
		}
	}()

	dbRepo := svc.db.GetRepository()
	rows, err := dbRepo.FindCorrectionSheet(context.Background())
	if err != nil {
		t.Fatalf("Failed to fetch correction sheets: %v", err)
	}

	if len(rows) == 0 {
		t.Skip("No data in database; skipping test")
	}

	for _, row := range rows {
		t.Run(fmt.Sprintf("TestItem_%s_%d", row.Title, row.Page), func(t *testing.T) {
			svc.runSingleTest(t, row)
		})
	}
}

func (s *TestService) runFailureTest(t *testing.T, title string, page int32, code int64, old, new string) {
	dbRepo := s.db.GetRepository()
	row, err := dbRepo.FindSpecificCorrectionSheet(context.Background(), generated.FindSpecificCorrectionSheetParams{
		Title: title,
		Page:  page,
	})
	if err != nil {
		t.Fatalf("Failed to fetch correction sheets: %v", err)
	}

	correction := services.Correction{
		DockerImage:    row.DockerImage,
		Command:        row.Command,
		SubmissionName: row.SubmissionName,
		FilesName:      row.FilesName,
		FilesContent:   row.FilesContent,
	}
	badPayload := strings.ReplaceAll(row.CorrectionContent, old, new)
	output, status, err := s.exercise.RunTest(correction, badPayload)
	if err != nil {
		t.Errorf("RunTest failed: %v", err)
	}
	if status.Error != nil {
		t.Error(status.Error)
	}
	if status.StatusCode != code {
		t.Errorf("Test expected code %d but got %d with output %s", code, status.StatusCode, output)
	}
	if output == "" {
		t.Errorf("Expected non-empty output, got empty")
	}
}

func TestRunTest_Fail(t *testing.T) {
	svc := setupService(t)
	defer func() {
		if err := svc.exercise.Cleanup(); err != nil {
			t.Error(err)
		}
	}()

	tests := []struct {
		title string
		page  int32
		code  int64
		old   string
		new   string
	}{
		{"Go", 1, 1, "return", "return 0 //"},
		{"Go", 3, 1, "return", "return 0 //"},
		{"Go", 4, 1, "return", `return "", 0, []string{}, nil //`},
		{"Go", 5, 1, "return", `return "", nil //`},
		// Add more test cases here as needed, e.g., {"AnotherTitle", 2}
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("TestFail_%s_%d", tt.title, tt.page), func(t *testing.T) {
			svc.runFailureTest(t, tt.title, tt.page, tt.code, tt.old, tt.new)
		})
	}
}

func TestRunTest_Concurrent(t *testing.T) {
	svc := setupService(t)
	defer func() {
		if err := svc.exercise.Cleanup(); err != nil {
			t.Error(err)
		}
	}()

	dbRepo := svc.db.GetRepository()
	rows, err := dbRepo.FindCorrectionSheet(context.Background())
	if err != nil {
		t.Fatalf("Failed to fetch correction sheets: %v", err)
	}

	if len(rows) == 0 {
		t.Skip("No data in database; skipping test")
	}

	concurrencyLevels := []int{3, 5, 10, 20}

	for _, level := range concurrencyLevels {
		t.Run(fmt.Sprintf("Concurrent_Level_%d", level), func(t *testing.T) {
			var wg sync.WaitGroup
			errChan := make(chan error, level)

			for i := range level {
				wg.Add(1)
				go func(index int) {
					defer wg.Done()
					row := rows[index%len(rows)]
					svc.runSingleTest(t, row)
				}(i)
			}

			wg.Wait()
			close(errChan)

			for err := range errChan {
				if err != nil {
					t.Errorf("Error in concurrent test at level %d: %v", level, err)
				}
			}
		})
	}
}
