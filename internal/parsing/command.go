package parsing

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	//"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	//"github.com/ArmanSandhu/CovertPi/internal/security"
)


func RunCommand(conn net.Conn, commandObj models.Cmd, stopRoutineChannel chan struct{}, cancelManager *models.CancelManager) {
	stopSignalChan := make(chan struct{})
	cancel := make(chan struct{})
	cancelManager.CancelMutex.Lock()
	cancelManager.CancelCommands[commandObj.Tool] = cancel
	cancelManager.CancelMutex.Unlock()

	go func() {
		select {
		case <- stopSignalChan:
			fmt.Println("Routine Stopped:", commandObj.Tool)
			return
		case <- stopRoutineChannel:
			fmt.Println("Stopping Routine:", commandObj.Tool)
			stopSignalChan <- struct{}{}
			return
		case <- cancel:
			fmt.Println("Cancel signal received for:", commandObj.Tool)
			stopSignalChan <- struct{}{}
			return
		}
	}()

	cmdSlices := strings.Fields(commandObj.Command)
	args := cmdSlices[1:]
	cmd := exec.Command(cmdSlices[0], args...)
	

	if commandObj.Tool == "nmap" {
		fmt.Println("Running Nmap Command!")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		go PrintNmapOutput(stdout, conn, stopSignalChan, stopRoutineChannel)
		cmd.Wait()
	}
	if commandObj.Tool == "airmon" {
		fmt.Println("Running Airmon Command!")
		//var jsonRes string
		cmdOut, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			conn.Close()
			return
		}

		switch {
		case strings.Contains(commandObj.Command, "check"):
			fmt.Println("Checking Airmon Wlan")
			result := CheckAirmon(string(cmdOut))
			fmt.Println(result)
		case strings.Contains(commandObj.Command, "start"):
			fmt.Println("Starting Airmon Wlan")
			result := StartStopAirmon(string(cmdOut))
			fmt.Println(result)
		case strings.Contains(commandObj.Command, "stop"):
			fmt.Println("Stopping Airmon Wlan")
			result := StartStopAirmon(string(cmdOut))
			fmt.Println(result)
		default:
			fmt.Println("Getting Wifi Interfaces!")
			result := GetAirmonInterfaces(string(cmdOut))
			fmt.Println(result)
		}

			// jsonRes, err := json.Marshal(result)
			// if err != nil {
			// 	fmt.Println("Error: ", err.Error())
			// 	conn.Close()
			// 	return
			// }
		
		// cipherText := security.Encrypt(jsonRes)
		// conn.Write([]byte(cipherText))
		// conn.Write([]byte("\n"))
		// conn.Close()
		// return
	}
}


