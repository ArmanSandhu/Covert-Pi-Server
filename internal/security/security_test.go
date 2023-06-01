package security

import (
	"reflect"
	"testing"
	"fmt"
	"crypto/aes"
	"encoding/json"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func TestUnpad(t *testing.T) {
	// Test Case 1
	input := []byte{84, 101, 115, 116, 85, 110, 112, 97, 100, 7, 7, 7, 7, 7, 7, 7}
	expectedOutput := []byte("TestUnpad")
	result := unpad(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 2
	input = []byte{84, 104, 105, 115, 32, 105, 115, 32, 97, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	expectedOutput = []byte("This is another test")
	result = unpad(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 2 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}

func TestPad(t *testing.T) {
	// Test Case 1
	input := []byte("TestPad")
	expectedOutput := []byte{84, 101, 115, 116, 80, 97, 100, 9, 9, 9, 9, 9, 9, 9 ,9 ,9}
	result := pad(input, aes.BlockSize)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	// Test Case 2
	input = []byte("This is another test")
	expectedOutput = []byte{84, 104, 105, 115, 32, 105, 115, 32, 97, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	result = pad(input, aes.BlockSize)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 2 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}

func TestDecrypt(t *testing.T) {
	// Test Case 1
	encCmdStr := "CVz+1ubD94k9u5F6zyI+C3uOubI8cZSaH6lGLRhCFVVU0UlnzP4Rmb0GLU6aSrm6"
	expectedOutput := "nmap 192.168.1.1"
	result := Decrypt(encCmdStr)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}

func TestEncrypt(t *testing.T) {
	// Test Case 1
	var hosts []models.Host
	var host *models.Host = new(models.Host)
	host.Name = "RT-AC5300-C3A0"
	host.Ip = "192.168.1.1"
	host.Status = "up"
	host.Latency = "0.0011s"
	host.Mac = "D0:17:C2:EA:C3:A0"
	host.Ports = []models.Port {
		models.Port {
			Num: "53",
			Porttype: "tcp",
			State: "open",
			Service: "domain",
		},
		models.Port  {
			Num: "80",
			Porttype: "tcp",
			State: "open",
			Service: "http",
		},
		models.Port  {
			Num: "515",
			Porttype: "tcp",
			State: "open",
			Service: "printer",
		},
		models.Port  {
			Num: "9100",
			Porttype: "tcp",
			State: "open",
			Service: "jetdirect",
		},
	}
	hosts = append(hosts, *host)
	jsonHosts, err := json.Marshal(hosts)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	result := Encrypt(jsonHosts)
	if string(jsonHosts) == result {
		t.Errorf("Test case 1 failed. JsonString: %v, Results: %v", string(jsonHosts), result)
	}
}