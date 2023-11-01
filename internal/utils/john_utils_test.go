package utils

import (
	"testing"
)

func TestConvertCap2John(t *testing.T) {
	// Test Case 1
	expectedOutput := "/home/kali/Desktop/Airodump_Captures/no_key.john"
	result, _, err := ConvertCap2John("/home/kali/Desktop/Airodump_Captures", "no_key")
	if err != nil {
		t.Errorf("Test Case 1 failed. Error occured: %v", err)
		return
	}
	if result != expectedOutput {
		t.Errorf("Test case 1 failed. Expected File: %v, Got: %v", expectedOutput, result)
		return
	}

	// Test Case 2
	expectedOutput = "/home/kali/Desktop/Airodump_Captures/valid_key.john"
	result, _, err = ConvertCap2John("/home/kali/Desktop/Airodump_Captures", "valid_key")
	if err != nil {
		t.Errorf("Test Case 2 failed. Error occured: %v", err)
		return
	}
	if result != expectedOutput {
		t.Errorf("Test case 2 failed. Expected File: %v, Got: %v", expectedOutput, result)
		return
	}
}

func TestConvertZip2John(t *testing.T) {
	// Test Case 1
	expectedOutput := "/home/kali/Desktop/Airodump_Captures/secureRar.txt"
	result, _, err := ConvertZip2John("/home/kali/Desktop/Airodump_Captures", "secureRar")
	if err != nil {
		t.Errorf("Test Case 1 failed. Error occured: %v", err)
		return
	}
	if result != expectedOutput {
		t.Errorf("Test case 1 failed. Expected File: %v, Got: %v", expectedOutput, result)
		return
	}
}

func TestConvertRar2John(t *testing.T) {
	// Test Case 1
	expectedOutput := "/home/kali/Desktop/Airodump_Captures/secureZip.txt"
	result, _, err := ConvertRar2John("/home/kali/Desktop/Airodump_Captures", "secureZip")
	if err != nil {
		t.Errorf("Test Case 1 failed. Error occured: %v", err)
		return
	}
	if result != expectedOutput {
		t.Errorf("Test case 1 failed. Expected File: %v, Got: %v", expectedOutput, result)
		return
	}
}


func TestGetAvailableFilesForCracking(t *testing.T) {
	//Test Case 1
	result, err := GetAvailableFilesForCracking("/home/kali/Desktop/Airodump_Captures")
	if len(result) == 0 {
		t.Errorf("Test Case 1 failed. No files found in given directory!")
	}
	if err != nil {
		t.Errorf("Test Case 1 failed. Error occured: %v", err)
	}
}