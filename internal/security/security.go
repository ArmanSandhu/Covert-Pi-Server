package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"fmt"
)

func LoadTLSCertificate(serverCertFile, serverKeyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {

	var key *rsa.PrivateKey

	// Load the server's certificate
	certPem, err := ioutil.ReadFile(serverCertFile)
	if err != nil {
		fmt.Println("There was an error loading the Server Cert File!")
		return nil, nil, err
	}

	// Load the server's private key
	keyPem, err := ioutil.ReadFile(serverKeyFile)
	if err != nil {
		fmt.Println("There was an error loading the Server Key File!")
		return nil, nil, err
	}

	// Decode and Parse the server certificate
	block, _ := pem.Decode(certPem)
	if block == nil {
		return nil, nil, errors.New("There was an error decoding the Server Cert File!")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("There was an error parsing the Server Cert File!")
		return nil, nil, err
	}

	// Decode and Parse the server key
	block, _ = pem.Decode(keyPem)
	if block == nil {
		return nil, nil, errors.New("There was an error decoding the Server Key File!")
	}
	if block.Type == "PRIVATE KEY" {
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			fmt.Println("There was an error parsing the Server Key File!")
			return nil, nil, err
		}
		
		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("Failed to convert loaded key into a RSA Private Key!")
		}
		key = rsaKey
	} else {
		return nil, nil, errors.New("Unsupported Key Type Found!")
	}

	return cert, key, nil
}
