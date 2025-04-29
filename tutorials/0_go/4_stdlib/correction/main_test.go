package main

import (
	"reflect"
	"testing"
)

func TestParseProfile(t *testing.T) {
	tests := []struct {
		name            string
		jsonStr         string
		expectedName    string
		expectedAge     int
		expectedHobbies []string
		expectError     bool
	}{
		{
			name:            "Valid JSON",
			jsonStr:         `{"person":{"name":"Bob","age":25},"hobbies":["Swimming","Cooking"]}`,
			expectedName:    "Bob",
			expectedAge:     25,
			expectedHobbies: []string{"Swimming", "Cooking"},
			expectError:     false,
		},
		{
			name:            "Invalid JSON",
			jsonStr:         `{"person":{"name":"Alice","age":30,"hobbies":["Reading"]}`,
			expectedName:    "",
			expectedAge:     0,
			expectedHobbies: nil,
			expectError:     true,
		},
		{
			name:            "Empty JSON",
			jsonStr:         `{"person":{"name":"","age":0},"hobbies":[]}`,
			expectedName:    "",
			expectedAge:     0,
			expectedHobbies: []string{},
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, age, hobbies, err := ParseProfile(tt.jsonStr)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseProfile() except error but didn't find")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseProfile() receive unexpected err %v", err)
			}
			if name != tt.expectedName {
				t.Errorf("ParseProfile() name = %v, expected %v", name, tt.expectedName)
			}
			if age != tt.expectedAge {
				t.Errorf("ParseProfile() age = %v, expected %v", age, tt.expectedAge)
			}
			if !reflect.DeepEqual(hobbies, tt.expectedHobbies) {
				t.Errorf("ParseProfile() hobbies = %v, expected %v", hobbies, tt.expectedHobbies)
			}
		})
	}
}
