package parsing

import (
	"testing"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func TestValidRunJohnCommand(t *testing.T) {
	// Test Case 1
	johnResult := models.John_Result{}
	cmdSlices := []string{"john", "--wordlist=/usr/share/john/password.lst", "--format=wpapsk", "/home/kali/Desktop/Airodump_Captures/valid_key.cap"}
	expectedResult := "success"
	RunJohnCommand(cmdSlices, &johnResult)
	if johnResult.Result != expectedResult {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, johnResult.Result)
		return
	}

	// Test Case 2
	johnResult = models.John_Result{}
	cmdSlices = []string{"john", "--wordlist=/usr/share/wordlists/rockyou.txt", "/home/kali/Desktop/Airodump_Captures/secureZip.zip"}
	RunJohnCommand(cmdSlices, &johnResult)
	if johnResult.Result != expectedResult {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expectedResult, johnResult.Result)
		return
	}

	// Test Case 3
	johnResult = models.John_Result{}
	cmdSlices = []string{"john", "--wordlist=/usr/share/wordlists/rockyou.txt", "/home/kali/Desktop/Airodump_Captures/secureRar.rar"}
	RunJohnCommand(cmdSlices, &johnResult)
	if johnResult.Result != expectedResult {
		t.Errorf("Test case 3 failed. Expected: %v, Got: %v", expectedResult, johnResult.Result)
		return
	}
}

func TestInvalidRunJohnCommand(t *testing.T) {
	// Test Case 1
	johnResult := models.John_Result{}
	cmdSlices := []string{"john", "--wordlist=/usr/share/list/rockyou.txt", "/home/kali/Desktop/Airodump_Captures/secureRar.rar"}
	expectedResult := "fail"
	RunJohnCommand(cmdSlices, &johnResult)
	if johnResult.Result != expectedResult {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, johnResult.Result)
		return
	}

	// Test Case 2
	johnResult = models.John_Result{}
	cmdSlices = []string{"john", "--wordlist=/usr/share/wordlists/rockyou.txt", "/home/kali/Desktop/Airodump_Captures/secure.txt"}
	RunJohnCommand(cmdSlices, &johnResult)
	if johnResult.Result != expectedResult {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expectedResult, johnResult.Result)
		return
	}
}

func TestGetAllCrckdPswdInDir(t *testing.T) {
	// Test Case 1
	directory := "/home/kali/Desktop/Airodump_Captures"
	crckdPswds, _, err := GetAllCrckdPswdInDir(directory)
	if len(crckdPswds) == 0 {
		t.Errorf("Test case 1 failed. No Passwords Found! Make sure that the passwords have been cracked!")
		return
	}

	// Test Case 2
	directory = "/home/kali/Desktop/Captures"
	crckdPswds, _, err = GetAllCrckdPswdInDir(directory)
	if err == nil {
		t.Errorf("Test case 2 failed!")
		return
	}
}