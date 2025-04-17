package main

import "testing"

func TestCalculateSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{
			name:     "Empty slice",
			numbers:  []int{},
			expected: 0,
		},
		{
			name:     "Single element",
			numbers:  []int{5},
			expected: 5,
		},
		{
			name:     "Multiple elements",
			numbers:  []int{1, 2, 3, 4},
			expected: 10,
		},
		{
			name:     "Larger slice",
			numbers:  []int{10, 20, 30, 40, 50, 60},
			expected: 210,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSum(tt.numbers)
			if result != tt.expected {
				t.Errorf("CalculateSum(%v) = %d; want %d", tt.numbers, result, tt.expected)
			}
		})
	}
}
