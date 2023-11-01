package parsing

import (
	"fmt"
	"os/exec"
	"os"
	"strings"
	"regexp"
	"path/filepath"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
)


func RunJohnCommand(cmdSlices []string, result *models.John_Result) {
	args := cmdSlices[1:]
	cmdFilePath := args[len(args) - 1]
	fmt.Println("File to Crack: ", cmdFilePath)

	// Get the file extension type of the file we wish to crack
	ext := filepath.Ext(cmdFilePath)
	fileTypeSupported := true

	// Get the base directory and base filename for further processing
	directory, filename := utils.SplitPathandFileName(cmdFilePath)

	var finalFilePath string
	var convDetails string
	var convErr error

	// Based on the extension, further conversion may be required
	switch ext {
	case ".txt":
		fmt.Println("Text File passed to John!")
		finalFilePath = cmdFilePath
	case ".cap":
		fmt.Println("Cap File passed to John! Conversion Required!")
		finalFilePath, convDetails, convErr = utils.ConvertCap2John(directory, filename)
	case ".zip":
		fmt.Println("Zip File passed to John! Conversion Required!")
		finalFilePath, convDetails, convErr = utils.ConvertZip2John(directory, filename)
	case ".rar":
		fmt.Println("Zip File passed to John! Conversion Required!")
		finalFilePath, convDetails, convErr = utils.ConvertRar2John(directory, filename)
	default:
		fmt.Println("Unknown File Type Passed!")
		fileTypeSupported = false
	}

	// If an unsupported filetype is passed, let the client know
	if !fileTypeSupported {
		result.Result = "fail"
		result.Details = "Unsupported Filetype was passed for cracking!"
		result.Error = "Unsupported Filetype!"
		return
	}

	// If an error occured while converting the file, let the client know
	if convErr != nil {
		result.Result = "fail"
		result.Details = "Error Occured during Conversion Process!"
		result.Error = convErr.Error()
		return
	}

	// Replace the filepath with the converted (if neccessary) filepath
	args[len(args) - 1] = finalFilePath

	
	// Run the command for John
	cmd := exec.Command(cmdSlices[0], args...)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occured while running John!")
		result.Result = "fail"
		result.Details = processJohnDetails(convDetails, string(cmdOut))
		result.Error = err.Error()
		return
	}

	result.Result = "success"
	username, password, details := showCrackedPassword(finalFilePath)
	result.Details = processJohnDetails(convDetails, string(cmdOut)) + details
	result.Passwords = make(map[string]string)
	result.Passwords[username] = password

}

func processJohnDetails(convDetails string, johnDetails string) string {
	fmt.Println("Processing Details from John File")
	var sb strings.Builder

	// Process Conversion Details if they exist
	if len(convDetails) != 0 {
		sb.WriteString("Conversion Details:")
		sb.WriteString("\n")
		lines := strings.Split(convDetails, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}


	sb.WriteString("John Details:")
	sb.WriteString("\n")

	// Process Details from running John
	lines := strings.Split(johnDetails, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Ignore certain lines to keep details concise
		if !strings.Contains(line, "OpenMP threads") || 
		   !strings.Contains(line, "Note:") || 
		   !strings.Contains(line, "Proceeding with wordlist") || 
		   !strings.Contains(line, "Press") || 
		   !strings.Contains(line, "Use the \"--show\"") ||
		   !strings.Contains(line, "Session completed.") {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func showCrackedPassword(fullFilePath string) (string, string, string) {
	fmt.Println("Inside showCrackedPassword")
	var sb strings.Builder
	var username string
	var password string
	pswdDisPattern := `(?m)([^:]+):([^:]+)(:.*)?`
	regex := regexp.MustCompile(pswdDisPattern)

	sb.WriteString("John Show Password Details:")
	sb.WriteString("\n")

	output, _ := exec.Command("john", "--show", fullFilePath).Output()
	cmdString := string(output)

	lines := strings.Split(cmdString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sb.WriteString(line)
			sb.WriteString("\n")
			lineMatches := regex.FindStringSubmatch(line)
			if len(lineMatches) >= 3 {
				username = lineMatches[1]
				password = lineMatches[2]
			}
		}
	}

	return username, password, sb.String()
}

func GetAllCrckdPswdInDir(directory string) (map[string]string, string, error) {
	var sb strings.Builder
	crckdPswds := make(map[string]string)

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error in Reading Directory!")
		return crckdPswds, sb.String(), err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" || filepath.Ext(file.Name()) == ".john" {
			username, password, details := showCrackedPassword(filepath.Join(directory, file.Name()))
			crckdPswds[username] = password
			sb.WriteString(details)
			sb.WriteString("\n")
		}
	}

	return crckdPswds, sb.String(), nil
}