package main

import (
	"testing"
)

func TestAddMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Basic addition", 3, 5, 8},
		{"Mid-range addition", 129, 14, 143},
		{"Exceeds max", 128, 250, 255},
		{"Zero values", 0, 0, 0},
		{"One max input", 255, 1, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddMax(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("AddMax(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
