package parsing

import (
	"reflect"
	"testing"
	"os/exec"
	"context"
	"time"
	"strings"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)


func TestParseAireplayOutputInvalid(t *testing.T) {
	// Test Case 1
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "aireplay-ng", "--deauth", "0", "-a", "94:83:C4:01:6B:3E", "-D", "wlan0")

	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	aireplayResult := models.Aireplay_Result{}

	err := cmd.Run()
	if err != nil && err.Error() != "signal: killed" {
		aireplayResult.Result = "fail"
		aireplayResult.Error = strings.TrimSpace(stderr.String())
		aireplayResult.Message = strings.TrimSpace(stdout.String())
	} else {
		ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
	}
	
	expectedResult := "fail"
	if !reflect.DeepEqual(expectedResult, aireplayResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, aireplayResult.Result)
	}


	// Test Case 2
	cmd = exec.CommandContext(ctx, "aireplay-ng", "--deauth", "0", "-a", "94:83:C4:01:6B:3E", "-D", "wlan2")

	stdout, stderr = &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	aireplayResult = models.Aireplay_Result{}

	err = cmd.Run()
	if err != nil && err.Error() != "signal: killed" {
		aireplayResult.Result = "fail"
		aireplayResult.Error = strings.TrimSpace(stderr.String())
		aireplayResult.Message = strings.TrimSpace(stdout.String())
	} else {
		ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
	}

	if !reflect.DeepEqual(expectedResult, aireplayResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, aireplayResult.Result)
	}

	// Test Case 3
	cmd = exec.CommandContext(ctx, "aireplay-ng", "--deauth", "0", "-a", "94:83:C4:01:6B:3E-D", "wlan1")

	stdout, stderr = &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	aireplayResult = models.Aireplay_Result{}

	err = cmd.Run()
	if err != nil && err.Error() != "signal: killed" {
		aireplayResult.Result = "fail"
		aireplayResult.Error = strings.TrimSpace(stderr.String())
		aireplayResult.Message = strings.TrimSpace(stdout.String())
	} else {
		ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
	}

	if !reflect.DeepEqual(expectedResult, aireplayResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, aireplayResult.Result)
	}
}

func TestParseAireplayOutputValid(t *testing.T) {
	// Test Case 1
	timeout := 5 * time.Second
	ctx1, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx1, "aireplay-ng", "--deauth", "0", "-a", "94:83:C4:01:6B:3E", "-D", "wlan1")

	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	aireplayResult := models.Aireplay_Result{}

	err := cmd.Run()
	if err != nil  && err.Error() != "signal: killed" {
		aireplayResult.Result = "fail"
		aireplayResult.Error = strings.TrimSpace(stderr.String())
		aireplayResult.Message = strings.TrimSpace(stdout.String())
	} else {
		ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
	}
	
	expectedResult := "success"
	if !reflect.DeepEqual(expectedResult, aireplayResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedResult, aireplayResult.Result)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Test Case 2
	cmd = exec.CommandContext(ctx2, "aireplay-ng", "--deauth", "5", "-a", "94:83:C4:01:6B:3E", "-D", "wlan1")

	stdout, stderr = &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	aireplayResult = models.Aireplay_Result{}

	err = cmd.Run()
	if err != nil  && err.Error() != "signal: killed" {
		aireplayResult.Result = "fail"
		aireplayResult.Error = strings.TrimSpace(stderr.String())
		aireplayResult.Message = strings.TrimSpace(stdout.String())
	} else {
		ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
	}
	
	if !reflect.DeepEqual(expectedResult, aireplayResult.Result) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expectedResult, aireplayResult.Result)
	}
}