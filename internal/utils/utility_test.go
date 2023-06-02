package utils

import (
	"reflect"
	"testing"
)

func TestTrimSlice(t *testing.T) {
	// Test Case 1
	input := []string{}
	expectedOutput := []string{}
	result := TrimSlice(input)
	if len(result) != 0 && len(expectedOutput) != 0 {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 2
	input = []string{"", ""}
	expectedOutput = []string{}
	result = TrimSlice(input)
	if len(result) != 0 && len(expectedOutput) != 0 {
		t.Errorf("Test case 2 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 3
	input = []string{"test1", "test2"}
	expectedOutput = []string{"test1", "test2"}
	result = TrimSlice(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 3 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 4
	input = []string{"", "test1", "", "test2"}
	expectedOutput = []string{"test1", "test2"}
	result = TrimSlice(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 4 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}

func TestRegSplit(t *testing.T) {
	// Test Case 1
	input := "This is a test"
	delim := " "
	expectedOutput := []string{"This", "is", "a", "test"}
	result := RegSplit(input, delim)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 2
	input = ""
	delim = " "
	expectedOutput = []string{}
	result = RegSplit(input, delim)
	if len(result) != 0 && len(expectedOutput) != 0 {
		t.Errorf("Test case 2 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 3
	input = "This12is34a56custom78test90case"
	delim = "[0-9]+"
	expectedOutput = []string{"This","is","a","custom","test","case"}
	result = RegSplit(input, delim)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 3 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}