package parsing

import (
	"fmt"
	"net"
	"strings"
	"io"
	"bufio"
	"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
)

func PrintNmapOutput(stdout io.ReadCloser, conn net.Conn, stopSignalChan, stopRoutineChannel chan struct{}) {
	var hosts []models.Host
	var host *models.Host
	var payload models.Nmap_Payload
	var sb strings.Builder
	var portFlag bool = false
	scanner := bufio.NewScanner(stdout)
	Loop:
		for scanner.Scan() {
			select {
			case <- stopSignalChan:
				fmt.Println("Received Stop Signal")
				break Loop
			case <- stopRoutineChannel:
				fmt.Println("Stopping printNmapOutput Routine!")
				break Loop
			default:
				text := scanner.Text()
				split := utils.RegSplit(text, "[^0-9A-Za-z.:-]+")
				split = utils.TrimSlice(split)
				if strings.Contains(text, "Nmap scan report") {
					fmt.Println("New Host")
					host = new(models.Host)
					host.Ip = split[len(split) - 1]
					if len(split) == 6 {
						host.Name = split[len(split) - 2]
					}
					continue
				}
				if strings.Contains(text, "Host is") {
					fmt.Println("Host Status")
					host.Status = split[2]
					host.Latency = split[3]
					continue
				}
				if len(split) == 3 && split[0] == "PORT" {
					portFlag = true
					continue
				}
				if strings.Contains(text, "MAC Address") {
					host.Mac = split[2]
					portFlag = false
					continue
				}
				if portFlag && len(split) > 0 {
					fmt.Println("Processing Port")
					port := models.Port {
						Num: split[0],
						Porttype: split[1],
						State: split[2],
						Service: split[3],
					}
					if (len(split) > 4) {
						fmt.Println("Adding Extra Port Data to Details")
						extraPortData := split[4:]
						extraPortDataString := strings.Join(extraPortData, " ")
						sb.WriteString(extraPortDataString)
					}
					host.Ports = append(host.Ports, port)
					continue
				}
				if strings.Contains(text, "Running:") {
					host.Running = split[1]
					continue
				}
				if strings.Contains(text, "OS CPE:") {
					oscpeData := strings.Split(text, " ")
					host.OSCPE = oscpeData[2]
					continue
				}
				if len(split) == 0 {
					fmt.Println("Adding to results")
					hosts = append(hosts, *host)
					portFlag = false
					continue
				}
				sb.WriteString(text)
				sb.WriteString("\n")
			}
		}
		payload.Hosts = hosts
		payload.Verbose = sb.String()
		nmap_payload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		conn.Write(nmap_payload)
		conn.Write([]byte("\n"))
		fmt.Println("End of Nmap Parse Function")
		conn.Close()
		return
}