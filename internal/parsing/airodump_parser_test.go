package parsing

import (
	"reflect"
	"testing"
	"os/exec"
	"os"
	"context"
	"time"
	"strings"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func TestParseAirodumpOutputValid(t *testing.T) {
	// Test Case 1
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filepath := "/home/kali/Desktop/Airodump_Captures/captured"
	filename := "/home/kali/Desktop/Airodump_Captures/captured.csv"

	cmd := exec.CommandContext(ctx, "airodump-ng", "--output-format", "csv", "-w", filepath, "-a", "wlan1")

	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	airodumpResult := models.Airodump_Result{}

	err := cmd.Run()
	if err != nil && err.Error() != "signal: killed"{
		airodumpResult.Result = "fail"
		airodumpResult.Message = strings.TrimSpace(stdout.String())
	}

	err = os.Rename("/home/kali/Desktop/Airodump_Captures/captured-01.csv", filename)
	if err != nil {
		t.Errorf("Test Case 1 failed. Error Renaming File: %v", err)
		return
	}

	ParseAirodumpOutput(filename, &airodumpResult)
	expectedOutput := "success"
	if !reflect.DeepEqual(expectedOutput, airodumpResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedOutput, airodumpResult.Result)
	}

	if len(airodumpResult.Airodump_APs) == 0 {
		t.Errorf("Test case 1 failed. Num of APs found is %v", len(airodumpResult.Airodump_APs))
	}
}

func TestParseAirodumpOutputInValid(t *testing.T) {
	// Test Case 1
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filepath := "/home/kali/Desktop/Airodump_Captures/captured"

	cmd := exec.CommandContext(ctx, "airodump-ng", "--output-format", "csv", "--write", filepath, "-a", "wlan0")

	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	airodumpResult := models.Airodump_Result{}

	err := cmd.Run()
	if err != nil && err.Error() != "signal: killed"{
		airodumpResult.Result = "fail"
		airodumpResult.Message = strings.TrimSpace(stdout.String())
	}

	expectedOutput := "fail"
	if !reflect.DeepEqual(expectedOutput, airodumpResult.Result) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedOutput, airodumpResult.Result)
	}
}