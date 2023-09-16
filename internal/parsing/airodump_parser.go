package parsing

import (
	"strings"
	"os"
	"io"
	"encoding/csv"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func ParseAirodumpOutput(filepath string, result *models.Airodump_Result) {
	if (result.Result != "fail") {
		var aps []models.Airodump_AP
		apMap := make(map[string]models.Airodump_AP)
		apFlag := false
		clientFlag := false
		errorFlag := false

		// Open file for parsing
		file, err := os.Open(filepath)
		if err != nil {
			result.Result = "fail"
			result.Message = err.Error()
			return
		}
		defer file.Close()

		// Create a CSV Reader
		csvReader := csv.NewReader(file)
		csvReader.FieldsPerRecord = -1
		
		for {
			line, err := csvReader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				result.Result = "fail"
				result.Message = err.Error()
				errorFlag = true
				break
			}
			if strings.Contains(line[0], "BSSID") {
				apFlag = true
				continue
			}
			if strings.Contains(line[0], "Station MAC") {
				clientFlag = true
				apFlag = false
				// Add AP to catch (not associated) clients
				ap := models.Airodump_AP {
					Mac: "(not associated)",
				}
				apMap[ap.Mac] = ap
				continue
			}
			if apFlag {
				ap := models.Airodump_AP {
					Mac: strings.TrimSpace(line[0]),
					First_Seen: strings.TrimSpace(line[1]),
					Last_Seen: strings.TrimSpace(line[2]),
					Channel: strings.TrimSpace(line[3]),
					Speed: strings.TrimSpace(line[4]),
					Privacy: strings.TrimSpace(line[5]),
					Cipher: strings.TrimSpace(line[6]),
					Auth: strings.TrimSpace(line[7]),
					Power: strings.TrimSpace(line[8]),
					Num_Packets: strings.TrimSpace(line[9]),
					IV: strings.TrimSpace(line[10]),
					Lan_IP: strings.TrimSpace(line[11]),
					ID_Len: strings.TrimSpace(line[12]),
					Essid: strings.TrimSpace(line[13]),
					Key: strings.TrimSpace(line[14]),
				}
				apMap[ap.Mac] = ap
				continue
			}
			if clientFlag {
				client := models.Airodump_Client {
					Mac: strings.TrimSpace(line[0]),
					First_Seen: strings.TrimSpace(line[1]),
					Last_Seen: strings.TrimSpace(line[2]),
					Power: strings.TrimSpace(line[3]),
					Num_Packets: strings.TrimSpace(line[4]),
					Bssid: strings.TrimSpace(line[5]),
					Probed_Essid: strings.TrimSpace(line[6]),
				}
				ap := apMap[client.Bssid]
				ap.Clients = append(ap.Clients, client)
				apMap[ap.Mac] = ap
				continue
			}
		}

		if !errorFlag {
			for _, ap := range apMap {
				aps = append(aps, ap)
			} 

			result.Airodump_APs = aps
			result.Message = "Data parsed successfully!"
			result.Result = "success"
		}
	}
}