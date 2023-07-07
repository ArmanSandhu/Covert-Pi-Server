package network

import (
	"fmt"
	"net"
	"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/parsing"
	"github.com/ArmanSandhu/CovertPi/internal/security"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

var (
	shutdownChannel = make(chan struct{})
)

type listenerInterface interface {
	Accept() (net.Conn, error)
	Close() error
	Addr() net.Addr
	Listen() (net.Listener, error)
}

type Listener struct {
	ConnHost string
	ConnPort string
	ConnType string
	Listener net.Listener
}

func (l *Listener) Listen() (net.Listener, error) {
	listener, err := net.Listen(l.ConnType, l.ConnHost + ":" + l.ConnPort)
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening on " + l.ConnHost + ":" + l.ConnPort)
	l.Listener = listener
	return listener, nil
}

func (l *Listener) Close() error {
	return l.Listener.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.Listener.Addr()
}

func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func StartServer(listen listenerInterface) (string, error) {
	// Start a routine to watch for a shutdown signal
	go func() {
		// Wait for a shutdown signal
		<-shutdownChannel
		// Stop the Server
		StopServer()
	}()

	//Listen for incoming connection
	listener, err := listen.Listen()
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		return "Error", err
	}
	// Close listener when program finishes
	defer listener.Close()
	
	// Wait for incoming connections
	for {
		select {
		case <-shutdownChannel:
			//Signal server to stop receiving new connections
			fmt.Println("Server Stopped!")
			return "Stop", nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error Accepting: ", err.Error())
				return "Error", err
			}
			go handleInConn(conn)
		}
	}
}

func StopServer() {
	close(shutdownChannel)
}

func handleInConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(buf[:cap(buf)])
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			break;
		}
		if reqLen > 0 {
			buf = buf[:reqLen]
			fmt.Printf("Read %d bytes\n", reqLen)
			encCmdString := string(buf)
			fmt.Println("Incoming Enc String: ", encCmdString)
			fmt.Println("Beginning Decryption!")
			cmdString := security.Decrypt(encCmdString)
			fmt.Println("Decrypted Cmd String: ", cmdString)

			var command models.Cmd
			err := json.Unmarshal([]byte(cmdString), &command)
			if err != nil {
				fmt.Println("Error Unmarshaling Cmd Obj: ", err)
				return
			}
			parsing.RunCommand(conn, command)
		}

		if reqLen == 0 {
			fmt.Println("Connection Closed!")
			break
		}
	}
}
