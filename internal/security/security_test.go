package security

import (
	"reflect"
	"testing"
)

func TestUnpad(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03, 0x03, 0x03, 0x03}
	expectedOutput := []byte{0x01, 0x02, 0x03}
	result := unpad(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 1 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	input = []byte{0x01, 0x02, 0x03, 0x03, 0x03, 0x03, 0x03}
	expectedOutput = []byte{0x01, 0x02, 0x03, 0x03}
	result = unpad(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 2 failed. Expected: %v, go: %v", expectedOutput, result)
	}

	input = []byte{0x01, 0x02, 0x03, 0x03}
	expectedOutput = []byte{0x01}
	result = unpad(input)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Test case 3 failed. Expected: %v, go: %v", expectedOutput, result)
	}
}