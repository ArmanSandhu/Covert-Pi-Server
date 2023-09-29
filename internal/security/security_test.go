package security

import (
	"testing"
)

func TestLoadingServerCertAndKey(t *testing.T) {
	// Test Case 1
	serverCertFile := "/home/kali/Desktop/CovertPiKey/server.crt"
	serverKeyFile := "/home/kali/Desktop/CovertPiKey/server.key"
	cert, key, err := LoadTLSCertificate(serverCertFile, serverKeyFile)
	if err != nil {
		t.Errorf("Test case 1 failed. Received an error loading valid certs and keys! Error: %v", err)
	}

	// Test Case 2
	cert, key, err = LoadTLSCertificate(serverKeyFile, serverCertFile)
	if cert != nil && key != nil {
		t.Errorf("Test case 2 failed. Cert and Key were returned even though incorrect files were provided!")
	}
}