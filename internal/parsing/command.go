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


func RunCommand(conn net.Conn, cmdString string) {
	cmdSlices := strings.Fields(cmdString)
	cmd := exec.Command(cmdSlices[0], cmdSlices[1])
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	go printOutput(stdout, conn)
	cmd.Wait()
}

func printOutput(stdout io.ReadCloser, conn net.Conn) {
	var hosts []models.Host
	var host *models.Host
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
		}
		if strings.Contains(text, "Host is") {
			fmt.Println("Host Status")
			host.Status = split[2]
			host.Latency = split[3]
		}
		if strings.Contains(text, "STATE SERVICE") {
			portFlag = true
			continue
		}
		if strings.Contains(text, "MAC Address") {
			host.Mac = split[2]
			portFlag = false
		}
		if portFlag && len(split) > 0 {
			fmt.Println("Processing Port")
			port := models.Port {
				Num: split[0],
				Porttype: split[1],
				State: split[2],
				Service: split[3],
			}
			host.Ports = append(host.Ports, port)
		}
		if len(split) == 0 {
			fmt.Println("Adding to results")
			hosts = append(hosts, *host)
			portFlag = false
		}
	}
	jsonHosts, err := json.Marshal(hosts)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	cipherText := security.Encrypt(jsonHosts)
	conn.Write([]byte(cipherText))
	conn.Write([]byte("\n"))
}
