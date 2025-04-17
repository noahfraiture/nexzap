package main

import (
	"strings"
	"testing"
)

func TestReadFileContent(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectedContent string
		expectedErrorMsg string
	}{
		{
			name:          "Empty filename",
			filename:      "",
			expectedContent: "",
			expectedErrorMsg: "filename cannot be empty",
		},
		{
			name:          "Non-existent file",
			filename:      "nonexistent.txt",
			expectedContent: "",
			expectedErrorMsg: "no such file or directory",
		},
		// Note: We can't easily test a successful read without creating a file during testing.
		// A more complete test would involve creating a temporary file.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ReadFileContent(tt.filename)
			if content != tt.expectedContent {
				t.Errorf("ReadFileContent() content = %q, expected %q", content, tt.expectedContent)
			}
			if err == nil && tt.expectedErrorMsg != "" {
				t.Errorf("ReadFileContent() expected error with message containing %q, got no error", tt.expectedErrorMsg)
			} else if err != nil && tt.expectedErrorMsg == "" {
				t.Errorf("ReadFileContent() unexpected error: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tt.expectedErrorMsg) {
				t.Errorf("ReadFileContent() error message = %q, expected to contain %q", err.Error(), tt.expectedErrorMsg)
			}
		})
	}
}
