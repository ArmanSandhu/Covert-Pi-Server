package parsing

import (
	"strings"
	"regexp"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
)

func GetAirmonInterfaces(output string) models.Airmon_Result {
	var airmonResult models.Airmon_Result
	airmonInterfaces := []models.Airmon_Interface{}
	var airmonInterface *models.Airmon_Interface
	interfaceModes := utils.GetWifiInterfacesMode()
	interfaceFlag := false
	regex := regexp.MustCompile(`\s+`)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			line = regex.ReplaceAllString(line, " ")
			if strings.Contains(line, "PHY") {
				interfaceFlag = true
				continue
			}
			if interfaceFlag {
				airmonInterface = parseInterface(line, interfaceModes)
				airmonInterfaces = append(airmonInterfaces, *airmonInterface)
			}
		}
	}

	airmonResult.Airmon_Interfaces = airmonInterfaces
	airmonResult.Result = "success"
	airmonResult.Details = "n/a"
	return airmonResult
}

func CheckAirmon(output string) models.Airmon_Result {
	processes := make(map[string]string)
	processFlag := false
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			split := strings.Split(line, " ")
			if split[0] == "PID" {
				processFlag = true
				continue
			}
			if processFlag {
				processes[split[1]] = split[0]
			}
		}
	}

	airmonResult := models.Airmon_Result{
		Airmon_Interfaces: []models.Airmon_Interface{},
        Details: "n/a",
        Result: "success",
		PIDS: processes,
	}
	
	return airmonResult
}

func StartStopAirmon(output string) models.Airmon_Result {
	var airmonResult models.Airmon_Result
	airmonInterfaces := []models.Airmon_Interface{}
	var airmonInterface *models.Airmon_Interface
	var sb strings.Builder
	interfaceModes := utils.GetWifiInterfacesMode()
	interfaceFlag := false
	detailsFlag := false
	airmonResult.Result = "success"
	airmonResult.Details = "n/a"
	regex := regexp.MustCompile(`\s+`)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			line = regex.ReplaceAllString(line, " ")
			if strings.Contains(line, "PHY") {
				interfaceFlag = true
				detailsFlag = true
				continue
			}
			if interfaceFlag {
				if strings.Contains(line, "phy") {
					airmonInterface = parseInterface(line, interfaceModes)
					airmonInterfaces = append(airmonInterfaces, *airmonInterface)
				} else {
					sb.WriteString(line)
					sb.WriteString("\n")
				}
				continue
			}
			if interfaceFlag && len(line) == 0 {
				interfaceFlag = false
				continue
			}
			if detailsFlag {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}
	airmonResult.Airmon_Interfaces = airmonInterfaces
	if sb.Len() != 0{
		airmonResult.Details = sb.String()
		if strings.Contains(airmonResult.Details, "ERROR") {
			airmonResult.Result = "fail"
		}
		if strings.Contains(airmonResult.Details, "You are trying to stop a device that isn't in monitor mode.") {
			airmonResult.Details = "You are trying to stop a device that isn't in monitor mode."
			airmonResult.Result = "fail"
		}
	}
	return airmonResult
}

func parseInterface(line string, interfaceModes map[string]string) *models.Airmon_Interface {
	airmonInterface := new(models.Airmon_Interface)
	split := strings.Split(line, " ")
	airmonInterface.PHY = split[0]
	airmonInterface.Interface = split[1]
	airmonInterface.Driver = split[2]
	airmonInterface.Chipset = strings.Join(split[3:], " ")
	airmonInterface.Mode = interfaceModes[airmonInterface.Interface]
	return airmonInterface
}