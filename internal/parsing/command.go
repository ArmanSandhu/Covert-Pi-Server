package parsing

import (
	"fmt"
	"net"
	"os/exec"
	"context"
	"time"
	"strings"
	"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/models"
	"github.com/ArmanSandhu/CovertPi/internal/security"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
)


func RunCommand(conn net.Conn, commandObj models.Cmd, stopRoutineChannel chan struct{}, cancelManager *models.CancelManager) {
	stopSignalChan := make(chan struct{})
	cancelFunc := make(chan struct{})
	cancelManager.CancelMutex.Lock()
	cancelManager.CancelCommands[commandObj.Tool] = cancelFunc
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
		case <- cancelFunc:
			fmt.Println("Cancel signal received for:", commandObj.Tool)
			stopSignalChan <- struct{}{}
			return
		}
	}()

	cmdSlices := strings.Fields(commandObj.Command)
	args := cmdSlices[1:]
	cmd := exec.Command(cmdSlices[0], args...)
	var ctx context.Context
	var cancel context.CancelFunc

	if commandObj.Tool == "nmap" {
		fmt.Println("Running Nmap Command!")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		go PrintNmapOutput(stdout, conn, stopSignalChan, stopRoutineChannel)
		cmd.Wait()
	}
	if commandObj.Tool == "airmon" {
		fmt.Println("Running Airmon Command!")
		var jsonRes []byte
		var airmonResult models.Airmon_Result
		cmdOut, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}

		switch {
		case strings.Contains(commandObj.Command, "check"):
			fmt.Println("Checking Airmon Wlan")
			airmonResult = CheckAirmon(string(cmdOut))
			fmt.Println(airmonResult)
		case strings.Contains(commandObj.Command, "start"):
			fmt.Println("Starting Airmon Wlan")
			airmonResult = StartStopAirmon(string(cmdOut))
			fmt.Println(airmonResult)
		case strings.Contains(commandObj.Command, "stop"):
			fmt.Println("Stopping Airmon Wlan")
			airmonResult = StartStopAirmon(string(cmdOut))
			fmt.Println(airmonResult)
		default:
			fmt.Println("Getting Wifi Interfaces!")
			airmonResult = GetAirmonInterfaces(string(cmdOut))
			fmt.Println(airmonResult)
		}

		jsonRes, err = json.Marshal(airmonResult)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			conn.Close()
			return
		}
		
		cipherText := security.Encrypt(jsonRes)
		conn.Write([]byte(cipherText))
		conn.Write([]byte("\n"))
		conn.Close()
		return
	}
	if commandObj.Tool == "airodump" {
		fmt.Println("Running Airodump Command!")
		directory := "/home/kali/Desktop/Airodump_Captures/"
		pattern := "-01"
		var jsonRes []byte
		airodumpResult := models.Airodump_Result{}

		filename := args[1]
		fullFilePath := directory + filename
		args[1] = fullFilePath
		timeout, _ := time.ParseDuration(args[len(args) - 1])
		args = args[:len(args) - 1]

		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd = exec.CommandContext(ctx, cmdSlices[0], args...)
		stdout, stderr := &strings.Builder{}, &strings.Builder{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil && err.Error() != "signal: killed"{
			airodumpResult.Result = "fail"
			airodumpResult.Message = strings.TrimSpace(stdout.String())
		}
		err = utils.RenameCaptureFiles(filename, directory, pattern)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			conn.Close()
			return
		}
		ParseAirodumpOutput(fullFilePath + ".csv", &airodumpResult)
		jsonRes, err = json.Marshal(airodumpResult)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			conn.Close()
			return
		}
		
		cipherText := security.Encrypt(jsonRes)
		conn.Write([]byte(cipherText))
		conn.Write([]byte("\n"))
		conn.Close()
		return
	}
	if commandObj.Tool == "aireplay" {
		fmt.Println("Running Aireplay Command!")
		aireplayResult := models.Aireplay_Result{}
		var jsonRes []byte
		timeout, _ := time.ParseDuration(args[len(args) - 1])
		args = args[:len(args) - 1]

		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd = exec.CommandContext(ctx, cmdSlices[0], args...)
		stdout, stderr := &strings.Builder{}, &strings.Builder{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil && err.Error() != "signal: killed" {
			aireplayResult.Result = "fail"
			aireplayResult.Error = strings.TrimSpace(stderr.String())
			aireplayResult.Message = strings.TrimSpace(stdout.String())
		} else {
			ParseAireplayOutput(stdout, &aireplayResult, cmd.ProcessState.ExitCode())
		}

		jsonRes, err = json.Marshal(aireplayResult)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			conn.Close()
			return
		}
		
		cipherText := security.Encrypt(jsonRes)
		conn.Write([]byte(cipherText))
		conn.Write([]byte("\n"))
		conn.Close()
		return
	}
}


