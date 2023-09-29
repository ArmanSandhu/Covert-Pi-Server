package network

import (
	"fmt"
	"net"
	"encoding/json"
	"crypto/tls"
	"os"
	"github.com/ArmanSandhu/CovertPi/internal/parsing"
	"github.com/ArmanSandhu/CovertPi/internal/security"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

var (
	shutdownChannel = make(chan struct{})
	stopRoutineChannel = make(chan struct{})
	cancelManager models.CancelManager
	ServerCertFile = "/home/kali/Desktop/CovertPiKey/server.crt"
	ServerKeyFile = "/home/kali/Desktop/CovertPiKey/server.key"
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
	TLSConfig *tls.Config
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

func StartServer(listen *Listener) (string, error) {
	// Load the SSL certificate and private key
	cert, key, err := security.LoadTLSCertificate(ServerCertFile, ServerKeyFile)
	if err != nil {
		fmt.Println("There was an error loading the server's TLS Config! Error: ", err)
		os.Exit(1)
	}

	// // Create a TLS configuration with the certificate and key
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	//Listen for incoming connection
	tlsListener, err:= tls.Listen(listen.ConnType, listen.ConnHost+":"+listen.ConnPort, tlsConfig)
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		return "Error", err
	}
	listen.Listener = tlsListener
	// Close listener when program finishes
	defer tlsListener.Close()

	// Start a routine to watch for a shutdown signal
	go func() {
		// Wait for a shutdown signal
		<-shutdownChannel
		// Stop the Server
		StopServer()
	}()

	fmt.Println("TLS Server is listening on " + listen.ConnHost + ":" + listen.ConnPort)

	cancelManager := models.NewCancelManager()

	// Wait for incoming connections
	for {
		select {
		case <-shutdownChannel:
			//Signal server to stop receiving new connections
			fmt.Println("Server Stopped!")
			return "Stop", nil
		default:
			conn, err := tlsListener.Accept()
			if err != nil {
				fmt.Println("Error Accepting: ", err.Error())
				return "Error", err
			}
			go handleInConn(conn, cancelManager)
		}
	}
}

func StopServer() {
	close(shutdownChannel)
	close(stopRoutineChannel)
	cancelManager.CancelMutex.Lock()
	defer cancelManager.CancelMutex.Unlock()
	for _, cancel := range cancelManager.CancelCommands {
		close(cancel)
	}
}

func handleInConn(conn net.Conn, cancelManager *models.CancelManager) {
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
			cmdString := string(buf)
			fmt.Println("Recieved Cmd String: ", cmdString)

			var command models.Cmd
			err := json.Unmarshal([]byte(cmdString), &command)
			if err != nil {
				fmt.Println("Error Unmarshaling Cmd Obj: ", err)
				return
			}
			if command.Command == "cancel" {
				fmt.Println("Cancel Command Received!")
				cancelManager.CancelMutex.Lock()
				cancel, found := cancelManager.CancelCommands[command.Tool]
				if found {
					close(cancel)
					delete(cancelManager.CancelCommands, command.Tool)
				}
				cancelManager.CancelMutex.Unlock()
				conn.Close()
			} else {
				parsing.RunCommand(conn, command, stopRoutineChannel, cancelManager)
			}
		}

		if reqLen == 0 {
			fmt.Println("Connection Closed!")
			break
		}
	}
}
