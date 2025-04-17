package main

import "testing"

func TestStudent(t *testing.T) {
	student := Student{Name: "John Doe", ID: 123}

	if got := student.GetName(); got != "John Doe" {
		t.Errorf("GetName() = %v; want John Doe", got)
	}

	if got := student.GetID(); got != 123 {
		t.Errorf("GetID() = %v; want 123", got)
	}
}
