package utils

import (
	"regexp"
	"os/exec"
	"fmt"
	"strings"
	"os"
	"os/user"
	"path/filepath"
	"bufio"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func RegSplit(text string, delim string) []string {
	reg := regexp.MustCompile(delim)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes) + 1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}

func TrimSlice(slices []string) []string {
	var results []string
	for _, slice := range slices {
		if len(slice) != 0 {
			results = append(results, slice)
		}
	}
	return results
}

func GetWifiInterfacesMode() map[string]string {
	modes := make(map[string]string)
	cmd, err := exec.Command("iwconfig").Output()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return modes
	}
	output := string(cmd)
	regex := regexp.MustCompile(`\s+`)
	lines := strings.Split(string(output), "\n")
	wlan := ""
	mode := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			line = regex.ReplaceAllString(line, " ")
			split := strings.Split(line, " ")
			if strings.Contains(split[0], "wl") {
				wlan = split[0]
			}
			if strings.Contains(split[0], "Mode") {
				parts := strings.SplitN(split[0], ":", 2)
				mode = parts[1]
				modes[wlan] = mode
			}
		}
		
	}

	return modes
}

func RenameCaptureFiles(filename string, directory string, pattern string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error in Reading Directory!")
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), filename) {
			fullSrcPath := filepath.Join(directory, file.Name())
			newDestPath := filepath.Join(directory, removePatternFromCaptureFiles(file.Name(), pattern))
			err := os.Rename(fullSrcPath, newDestPath)
			if err != nil {
				fmt.Println("Error Renaming File!")
				return err
			}
		}
	}

	return nil
}

func removePatternFromCaptureFiles(filename string, pattern string) string {
	baseFileName := strings.TrimSuffix(filename, filepath.Ext(filename))
	fileNameWOPattern := strings.TrimSuffix(baseFileName, pattern)
	newFileName := fileNameWOPattern + filepath.Ext(filename)
	return newFileName
}

func SplitPathandFileName(fullFilePath string) (string, string) {
	directory := filepath.Dir(fullFilePath)
	fileName := filepath.Base(fullFilePath)
	ext := filepath.Ext(fileName)
	fileNameWOExt := fileName[:len(fileName) - len(ext)]
	return directory, fileNameWOExt
}

func ResetRaspberryPiWifiAdapter() {
	unloadCmd := exec.Command("rmmod", "brcmfmac")
	err := unloadCmd.Run()
	if err != nil {
		fmt.Println("There was an error while unloading the Raspberry Pi's driver: ", err)
		return
	}

	loadCmd := exec.Command("modprobe", "brcmfmac")
	err = loadCmd.Run()
	if err != nil {
		fmt.Println("There was an error while loading the Raspberry Pi's driver: ", err)
		return
	}

	fmt.Println("Raspberry Pi 4 driver reset!")
}

func ReadCovertPiConfigFile(filePath string) (*models.Covert_Pi_Config, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	covertPiConfig := &models.Covert_Pi_Config{}
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		configLine := scanner.Text()
		lineSlice := strings.SplitN(configLine, "=", 2)
		if len(lineSlice) != 2 {
			fmt.Println("Unexpected line within Conf File!")
			continue
		}
		configKey := strings.TrimSpace(lineSlice[0])
		configVal := strings.TrimSpace(lineSlice[1])
		switch configKey {
		case "hostIP":
			covertPiConfig.HostIP = configVal
		case "hostPort":
			covertPiConfig.HostPort = configVal
		case "captureDir":
			covertPiConfig.CaptureDir = configVal
		case "serverKeyFilePath":
			covertPiConfig.ServerKeyFilePath = configVal
		case "serverCertFilePath":
			covertPiConfig.ServerCertFilePath = configVal
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return covertPiConfig, nil
}

func GetConfFilePath(username string) (string, error) {
	curr, err := user.Lookup(username)
	if err != nil {
		fmt.Println("Error retrieving current user: ", err)
		return "", err
	}

	return filepath.Join(curr.HomeDir, "Desktop/CovertPiServerDetails/covertpi.conf"), nil
}