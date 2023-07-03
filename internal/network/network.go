package network

import (
	"fmt"
	"net"
	"os"
	"github.com/ArmanSandhu/CovertPi/internal/parsing"
	"github.com/ArmanSandhu/CovertPi/internal/security"
)

var (
	//ConnHost = "192.168.1.60"
	//ConnPort = "8080"
	//ConnType = "tcp"
	shutdownChannel = make(chan struct{})
	shutdownChannelPtr = &shutdownChannel
)

type listenerInterface interface {
	Accept() (net.Conn, error)
	Close() error
	Addr() net.Addr
	//Listen(network, address string) (net.Listener, error)
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

func StartServer(listen listenerInterface) {
	//Listen for incoming connection
	//listener, err := net.Listen(ConnType, ConnHost + ":" + ConnPort)
	//listener, err := listen.Listen(listen.ConnType, listen.ConnHost + ":" + listen.ConnPort)
	listener, err := listen.Listen()
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		os.Exit(1)
	}
	// Close listener when program finishes
	defer listener.Close()

	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)
	
	// Wait for incoming connections
	for {
		select {
		case <- shutdownChannel:
			//Signal server to stop receiving new connections
			StopServer()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error Accepting: ", err.Error())
				os.Exit(1)
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
