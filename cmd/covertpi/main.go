package main

import (
	"fmt"
	"os"
	"github.com/ArmanSandhu/CovertPi/internal/network"
	"github.com/ArmanSandhu/CovertPi/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide current username!")
		return
	}

	currUser := os.Args[1]
	currUserConfFilePath, err := utils.GetConfFilePath(currUser)
	if err != nil {
		fmt.Println("There was an error retrieving the Conf file location for the user. Error: ", err)
		return
	}

	fmt.Println("Reading Configuration File")
	config, err := utils.ReadCovertPiConfigFile(currUserConfFilePath)
	if err != nil {
		fmt.Println("There was an error reading the Conf file! Error: ", err)
		return
	}

	listener := &network.Listener{ConnHost: config.HostIP, ConnPort: config.HostPort, ConnType: "tcp"}
	fmt.Println("Starting Server")
	network.StartServer(listener, config.ServerKeyFilePath, config.ServerCertFilePath, config.CaptureDir)	
}