package main

import (
	"fmt"
	"github.com/ArmanSandhu/CovertPi/internal/network"
)

func main() {
	fmt.Println("Starting Server")
	network.StartServer()	
}