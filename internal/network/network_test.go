package network

import (
	"testing"
	"fmt"
	"net"
	"time"
	"sync"
)

// Setup mock structs for testing
type mockConnection struct{
	remoteAddr net.Addr
}
type mockListener struct {
	ConnHost string
	ConnPort string
	ConnType string
}
// ----MOCK LISTENER FUNCTIONS----
func (ml *mockListener) Accept() (net.Conn, error) {
	fmt.Println("Accepted Incoming Connection!")
	time.Sleep(time.Millisecond)
	return &mockConnection{}, nil
}

func (ml *mockListener) Close() error {
	fmt.Println("Closing Client Connection!")
	time.Sleep(time.Millisecond)
	return nil
}

func (ml *mockListener) Addr() net.Addr {
	return &net.IPAddr{
		IP: net.ParseIP(ml.ConnHost),
	}
}

func (ml *mockListener) Listen() (net.Listener, error) {
	return &mockListener{}, nil
} 

// ----MOCK CONNECTION FUNCIONS----
func (mc *mockConnection) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (mc *mockConnection) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (mc *mockConnection) Close() error {
	return nil
}

func (mc *mockConnection) LocalAddr() net.Addr {
	return nil
}

func (mc *mockConnection) RemoteAddr() net.Addr {
	return mc.remoteAddr
}

func (mc *mockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestStartServer(t *testing.T) {
	mockListener := &mockListener{
		ConnHost: "192.168.1.60",
		ConnPort: "8037",
		ConnType: "tcp",
	}

	listener := &Listener{
		ConnHost: mockListener.ConnHost,
		ConnPort: mockListener.ConnPort,
		ConnType: mockListener.ConnType,
	}
	
	var err error
	var actualOutput string
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		actualOutput, err = StartServer(listener, "/home/kali/Desktop/CovertPiKey/server.key", "/home/kali/Desktop/CovertPiKey/server.crt", "")
	}()

	time.Sleep(time.Millisecond)

	shutdownChannel <- struct{}{}
	fmt.Println("Shutdown Signal Sent")
	wg.Wait()
	
	if err != nil {
		t.Errorf("Error occured in StartServer: %v", err)
	}

	expectedOutput := "Stop"
	if actualOutput != expectedOutput {
		t.Errorf("Expected Server Output: %s, but got: %s\n", expectedOutput, actualOutput)
	}
}