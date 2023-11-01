package utils

import (
	"regexp"
	"os/exec"
	"fmt"
	"strings"
	"os"
	"path/filepath"
	"io/ioutil"
)

func ConvertCap2John(directory string, filename string) (string, string, error) {
	capFilePath := directory + "/" + filename + ".cap"
	johnFilePath := directory + "/" + filename + ".john"
	
	cmd := exec.Command("wpapcap2john", capFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occured in running wpapcap2john!")
		return "fail", processCap2JohnOutput(string(output)), err
	}

	err = os.WriteFile(johnFilePath, output, 0644)
	if err != nil {
		fmt.Println("Error occured in writing wpapcap2john results to new file!")
		return "fail", processCap2JohnOutput(string(output)), err
	}
	
	return johnFilePath, processCap2JohnOutput(string(output)), nil
}

func processCap2JohnOutput(output string) string {
	pmkPattern := `^\S+:\$WPAPSK\$[\S#\.]+`
	regex := regexp.MustCompile(pmkPattern)
	var sb strings.Builder

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != 0 && !regex.MatchString(line) {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func ConvertZip2John(directory string, filename string) (string, string, error) {
	zipFilePath := directory + "/" + filename + ".zip"
	txtFilePath := directory + "/" + filename + ".txt"

	cmd := exec.Command("zip2john", zipFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occured in running zip2john!")
		return "fail", processCap2JohnOutput(string(output)), err
	}

	err = os.WriteFile(txtFilePath, output, 0644)
	if err != nil {
		fmt.Println("Error occured in writing zip2john results to new file!")
		return "fail", processCap2JohnOutput(string(output)), err
	}
	
	return txtFilePath, string(output), nil
}

func ConvertRar2John(directory string, filename string) (string, string, error) {
	rarFilePath := directory + "/" + filename + ".rar"
	txtFilePath := directory + "/" + filename + ".txt"

	cmd := exec.Command("rar2john", rarFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occured in running rar2john!")
		return "fail", string(output), err
	}

	err = os.WriteFile(txtFilePath, output, 0644)
	if err != nil {
		fmt.Println("Error occured in writing rar2john results to new file!")
		return "fail", string(output), err
	}
	
	return txtFilePath, string(output), nil
}

func GetAvailableFilesForCracking(directory string) ([]string, error) {
	var files []string

	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Error occured in retrieving files! ", err)
		return files, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}

	return files, err
}