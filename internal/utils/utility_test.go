package utils

import (
	"reflect"
	"testing"
	"os"
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

func TestGetWifiInterfacesMode(t *testing.T) {
	// Test Case 1
	expectedOutput := make(map[string]string)
	expectedOutput["wlan0"] = "Managed"
	expectedOutput["wlan1"] = "Managed"
	result := GetWifiInterfacesMode()
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedOutput, result)
	}
}

func TestRenameCaptureFiles(t *testing.T) {
	// Test Case 1
	directory := "/home/kali/Desktop/Airodump_Captures/"
	fileName := "test"
	fileNameWPattern := "test-01.txt"
	finalFile := "test.txt"
	pattern := "-01"
	foundFileFlag := false

	_, err := os.Create(directory + fileNameWPattern)
	if err != nil {
		t.Errorf("Test Case 1 failed. Error creating Test File: %v", err)
		return
	}

	err = RenameCaptureFiles(fileName, directory, pattern)
	if err != nil {
		t.Errorf("Test Case 1 failed. Error recieved: %v", err)
		return
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		t.Errorf("Test Case 1 failed. Error reading Directory: %v", err)
		return
	}

	for _, file := range files {
		if file.Name() == finalFile {
			foundFileFlag = true
		}
	}

	if !foundFileFlag {
		t.Errorf("Test Case 1 failed. Test File Not Renamed Correctly: %v", err)
		return
	}

	err = os.Remove(directory + finalFile)
	if err != nil {
		t.Errorf("Test Case 1 failed. Error deleting Test File: %v", err)
		return
	}
}