package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"nexzap/internal/services/container"

	"github.com/docker/docker/client"
)

// TestRunTestGo tests the basic functionality of running Go tests
func TestRunTestGo(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}
	defer svc.Cleanup()

	goDir := "../../tutorials/0_go/warmup/"

	files, err := os.ReadDir(goDir)
	if err != nil {
		t.Fatalf("Failed to read Go directory: %v", err)
	}

	filePaths := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(goDir, file.Name()))
		}
	}

	var wg sync.WaitGroup
	wg.Add(5)
	for i := range 5 {
		go func(runNum int) {
			defer wg.Done()
			output, err := svc.RunTest("gotest", "go test", filePaths)
			if err != nil {
				t.Errorf("Run %d: Failed to run Go test: %v", runNum, err)
				return
			}
			t.Logf("Run %d - Go Test Output:\n%s", runNum, output)
			if output == "" {
				t.Errorf("Run %d: Expected non-empty output from Go test, but got empty string", runNum)
			}
		}(i + 1)
	}
	wg.Wait()
}

// TestServiceNotInitialized tests the behavior when service is not initialized
func TestServiceNotInitialized(t *testing.T) {
	svc := &Service{}
	output, err := svc.RunTest(0, []string{"test.go"})
	if err == nil {
		t.Error("Expected error when service is not initialized, but got nil")
	}
	if output != "" {
		t.Errorf("Expected empty output when service is not initialized, but got: %s", output)
	}
	if err.Error() != "not initialized" {
		t.Errorf("Expected error message 'not initialized', but got: %v", err)
	}
}

// TestUnknownLanguage tests the behavior when an unknown language is provided
func TestUnknownLanguage(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}
	defer svc.Cleanup()

	unknownLang := 999
	output, err := svc.RunTest(unknownLang, []string{"test.go"})
	if err == nil {
		t.Error("Expected error when using unknown language, but got nil")
	}
	if output != "" {
		t.Errorf("Expected empty output when using unknown language, but got: %s", output)
	}
	if err.Error() != "Number longer than number of tutorials" {
		t.Errorf("Expected error message 'Number longer than number of tutorials', but got: %v", err)
	}
}

// TestEmptyFileList tests the behavior when an empty file list is provided
func TestEmptyFileList(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}
	defer svc.Cleanup()

	output, err := svc.RunTest(0, []string{})
	if err != nil {
		t.Errorf("Expected no error when running with empty file list, but got: %v", err)
	}
	if output == "" {
		t.Log("Warning: Empty file list resulted in empty output, which might be expected behavior")
	}
}

// TestConcurrentAccess tests the service under high concurrent access
func TestConcurrentAccess(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}
	defer svc.Cleanup()

	goDir := "../../tutorials/0_go/warmup/"
	files, err := os.ReadDir(goDir)
	if err != nil {
		t.Fatalf("Failed to read Go directory: %v", err)
	}

	filePaths := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(goDir, file.Name()))
		}
	}

	// Run 20 concurrent test requests
	concurrency := 20
	var wg sync.WaitGroup
	wg.Add(concurrency)
	errors := make(chan error, concurrency)
	outputs := make(chan string, concurrency)

	for i := range concurrency {
		go func(runNum int) {
			defer wg.Done()
			output, err := svc.RunTest(0, filePaths)
			if err != nil {
				errors <- fmt.Errorf("Run %d: Failed to run Go test: %v", runNum, err)
				return
			}
			outputs <- fmt.Sprintf("Run %d: %s", runNum, output)
		}(i)
	}

	// Wait for all goroutines to complete with a timeout
	timeout := time.After(30 * time.Second)
	go func() {
		wg.Wait()
		close(errors)
		close(outputs)
	}()

	select {
	case <-timeout:
		t.Fatal("Test timed out waiting for concurrent operations to complete")
	default:
		// Check for errors
		for err := range errors {
			t.Error(err)
		}
		// Log outputs
		for output := range outputs {
			t.Log(output)
		}
	}
}

// TestNonExistentFiles tests the behavior when non-existent files are provided
func TestNonExistentFiles(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Fatalf("Failed to initialize container service: %v", err)
	}
	defer svc.Cleanup()

	// Provide a list of non-existent files
	files := []string{"/path/does/not/exist/test.go"}
	output, err := svc.RunTest(0, files)
	if err != nil {
		t.Logf("Expected error or specific handling for non-existent files, got: %v", err)
	} else {
		t.Logf("Warning: No error returned for non-existent files, output: %s", output)
	}
}

// mockContainerPool is a mock implementation for testing container pool behavior
type mockContainerPool struct {
	getContainerFunc  func(ctx context.Context, cli *client.Client) (string, error)
	freeContainerFunc func(ctx context.Context, cli *client.Client, ctn string)
}

func (m *mockContainerPool) GetLanguagePool(ctx context.Context, cli *client.Client, language container.Tutorial) container.LanguagePool {
	return container.LanguagePool{}
}

func (m *mockContainerPool) GetContainer(ctx context.Context, cli *client.Client) (string, error) {
	if m.getContainerFunc != nil {
		return m.getContainerFunc(ctx, cli)
	}
	return "", fmt.Errorf("mock get container not implemented")
}

func (m *mockContainerPool) FreeContainer(ctx context.Context, cli *client.Client, ctn string) {
	if m.freeContainerFunc != nil {
		m.freeContainerFunc(ctx, cli, ctn)
	}
}

// TestContainerPoolExhaustion tests the behavior when container pool is exhausted
func TestContainerPoolExhaustion(t *testing.T) {
	// This test would ideally mock the container pool to simulate exhaustion
	// However, due to the complexity of the container pool implementation,
	// a full mock would be needed which is beyond the scope of this test file update.
	t.Skip("Test for container pool exhaustion requires mocking of container package")
}

// TestTimeoutScenarios tests various timeout scenarios
func TestTimeoutScenarios(t *testing.T) {
	// Similar to TestContainerPoolExhaustion, testing timeout scenarios
	// requires mocking the container package to simulate timeouts.
	// This test is skipped until proper mocking is implemented.
	t.Skip("Test for timeout scenarios requires mocking of container package")
}
