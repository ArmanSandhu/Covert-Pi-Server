package parsing

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"io"
	"bufio"
	"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
	"github.com/ArmanSandhu/CovertPi/internal/security"
)


func RunCommand(conn net.Conn, commandObj models.Cmd) {
	cmdSlices := strings.Fields(commandObj.Command)
	args := cmdSlices[1:]
	cmd := exec.Command(cmdSlices[0], args...)
	stdout, _ := cmd.StdoutPipe()
	if commandObj.Tool == "nmap" {
		fmt.Println("Running Nmap Command!")
		cmd.Start()
		go printNmapOutput(stdout, conn)
		cmd.Wait()
	}
}

func printNmapOutput(stdout io.ReadCloser, conn net.Conn) {
	var hosts []models.Host
	var host *models.Host
	var payload models.Nmap_Payload
	var sb strings.Builder
	var portFlag bool = false
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
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
		if strings.Contains(text, "STATE SERVICE") {
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
	//fmt.Println(sb.String())
	payload.Hosts = hosts
	payload.Verbose = sb.String()
	//jsonHosts, err := json.Marshal(hosts)
	nmap_payload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	//cipherText := security.Encrypt(jsonHosts)
	cipherText := security.Encrypt(nmap_payload)
	conn.Write([]byte(cipherText))
	conn.Write([]byte("\n"))
}
