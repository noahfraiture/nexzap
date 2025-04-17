package main

import (
	"errors"
	"os"
)

// ReadFileContent reads the content of a file and returns it as a string, along with any error
func ReadFileContent(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("filename cannot be empty")
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
