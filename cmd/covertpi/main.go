package main

import (
	"fmt"
	"github.com/ArmanSandhu/CovertPi/internal/network"
)

func main() {
	fmt.Println("Starting Server")
	listener := &network.Listener{ConnHost: "192.168.1.60", ConnPort: "8080", ConnType: "tcp"}
	network.StartServer(listener)	
}