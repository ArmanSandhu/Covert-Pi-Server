package utils

import (
	"regexp"
	"os/exec"
	"fmt"
	"os"
	//"io/ioutil"
	"strings"
	"bufio"
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

func GetWifiInterfaces() map[string]string {
	interfaces := make(map[string]string)
	filepath := "/home/kali/Desktop/CovertPi Server/internal/utils/interfaces.txt"
	cmd, err := exec.Command("/bin/sh", "/home/kali/Desktop/CovertPi Server/internal/utils/findInterfaces.sh").Output()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return interfaces
	}
	fmt.Println("Script Executed!")
	output := string(cmd)
	fmt.Println(output)
	cleanFile(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return interfaces
	}
	defer file.Close()
	
	device := ""
	wlan := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Device") {
			split := strings.Split(line, " ")
			if (len(split) > 1) {
				device = strings.Join(split[1:], " ")
			} else {
				device = split[1]
			}
		}
		if strings.Contains(line, "IF") {
			split := strings.Split(line, " ")
			wlan = split[1]
			if (strings.Contains(wlan, "wl")) {
				fmt.Println("Found Wireless Interface")
				interfaces[wlan] = device 
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}


	return interfaces
}

func cleanFile(filepath string) {
	cleanedLines := processFile(filepath)
	err := writeLinesToFile(filepath, cleanedLines)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("File cleaned!")
}

func removeControlCharacters(line string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 || r == '\n' {
			return r
		}
		return -1
	}, line)
}

func processFile(filepath string) []string {
	var cleanedLines []string

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error:", err)
		return cleanedLines
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cleanedLine := removeControlCharacters(line)
		cleanedLines = append(cleanedLines, cleanedLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	return cleanedLines
}

func writeLinesToFile(filepath string, lines []string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}