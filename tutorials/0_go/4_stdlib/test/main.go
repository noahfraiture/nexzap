package main

import (
	"encoding/json"
)

// Person represents the nested person data in JSON
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Profile represents the full structure of the JSON data
type Profile struct {
	Person  Person   `json:"person"`
	Hobbies []string `json:"hobbies"`
}

// ParseProfile decodes a JSON string into a Profile struct and returns name, age, hobbies, and any error
func ParseProfile(jsonStr string) (string, int, []string, error) {
	var profile Profile
	err := json.Unmarshal([]byte(jsonStr), &profile)
	if err != nil {
		return "", 0, nil, err
	}
	return profile.Person.Name, profile.Person.Age, profile.Hobbies, nil
}
