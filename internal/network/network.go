package network

import (
	"fmt"
	"net"
	"os"
	"github.com/ArmanSandhu/CovertPi/internal/parsing"
	"github.com/ArmanSandhu/CovertPi/internal/security"
)

const (
	ConnHost = "192.168.1.60"
	ConnPort = "8080"
	ConnType = "tcp"
)

func StartServer() {
	//Listen for incoming connection
	listener, err := net.Listen(ConnType, ConnHost + ":" + ConnPort)
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		os.Exit(1)
	}
	// Close listener when program finishes
	defer listener.Close()

	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)
	
	// Wait for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			os.Exit(1)
		}
		go handleInConn(conn)
	}
}

func handleInConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(buf[:cap(buf)])
		buf = buf[:reqLen]
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			break;
		}
		fmt.Printf("Read %d bytes\n", reqLen)
		encCmdString := string(buf)
		fmt.Println("Incoming Enc String: ", encCmdString)
		fmt.Println("Beginning Decryption!")
		cmdString := security.Decrypt(encCmdString)
		fmt.Println("Decrypted Cmd String: ", cmdString)
		parsing.RunCommand(conn, cmdString)
	}
}
