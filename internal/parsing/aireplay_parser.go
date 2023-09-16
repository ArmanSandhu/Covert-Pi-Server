package parsing

import (
	"strings"
	"fmt"
	"regexp"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func ParseAireplayOutput(output *strings.Builder, result *models.Aireplay_Result, exitCode int) {
	timeStampPattern := `^\d{2}:\d{2}:\d{2}`
	regex := regexp.MustCompile(timeStampPattern)
	deauthSuccess := false

	if exitCode == 0 {
		lines := strings.Split(output.String(), "\n")

		// Iterate to find timestamps which are an indication of a succesful deauth attack
		for _, line := range lines {
			timestamp := regex.FindString(line)
			if timestamp != "" {
				fmt.Println("Found Timestamp! Deauths Were Sent!")
				deauthSuccess = true
			}
		}
	}
	if exitCode == -1 {
		fmt.Println("Exit Code of -1 for Aireplay Cmd with Ctx! Deauths were probably being sent continuously and timeout occured!")
		deauthSuccess = true
	}

	if deauthSuccess {
		result.Message = "Deauth Sent Successfully!"
		result.Result = "success"
	} else {
		fmt.Println("No Timestamps Found! Error occured!")
		result.Message = output.String()
		result.Result = "fail"
	}
}