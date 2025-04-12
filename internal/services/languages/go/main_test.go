package main

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        inputA   int
        inputB   int
        expected int
    }{
        {
            name:     "Positive numbers",
            inputA:   2,
            inputB:   3,
            expected: 5,
        },
        {
            name:     "Negative numbers",
            inputA:   -1,
            inputB:   -1,
            expected: -2,
        },
        {
            name:     "Zero and positive",
            inputA:   0,
            inputB:   5,
            expected: 5,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.inputA, tt.inputB)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.inputA, tt.inputB, result, tt.expected)
            }
        })
    }
}
