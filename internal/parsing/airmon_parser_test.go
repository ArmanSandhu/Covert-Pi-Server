package parsing

import (
	"reflect"
	"testing"
	"os/exec"
	"sort"
	"github.com/ArmanSandhu/CovertPi/internal/models"
)

func TestGetAirmonInterfaces(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng").Output()
	if err != nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	airmonInterfaces := []models.Airmon_Interface{}
	wlan0 := models.Airmon_Interface{PHY: "phy0", Interface: "wlan0", Driver: "88XXau", Chipset: "Realtek Semiconductor Corp. Realtek 8812AU/8821AU 802.11ac WLAN Adapter [USB Wireless Dual-Band Adapter 2.4/5Ghz]", Mode: "Managed"}
	wlan1 := models.Airmon_Interface{PHY: "phy1", Interface: "wlan1", Driver: "brcmfmac", Chipset: "Broadcom 43430", Mode: "Managed"}
	airmonInterfaces = append(airmonInterfaces, wlan0)
	airmonInterfaces = append(airmonInterfaces, wlan1)
	airmonResult := models.Airmon_Result{Airmon_Interfaces: airmonInterfaces, Result: "success", Details: "n/a"}
	result := GetAirmonInterfaces(string(cmdOut))
	if !reflect.DeepEqual(result, airmonResult) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", airmonResult, result)
	}
}

func TestCheckAirmon(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng", "check").Output()
	if err != nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	expectedKeys := []string{"dhclient", "NetworkManager", "wpa_supplicant"}
	resultMap := CheckAirmon(string(cmdOut))
	result := make([]string, 0, len(resultMap))
	for key := range resultMap {
		result = append(result, key)
	}
	sort.Strings(expectedKeys)
	sort.Strings(result)
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expectedKeys, result)
	}
}

func TestStartAirmonOnValid(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng", "start", "wlan0").Output()
	if err != nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	airmonInterfaces := []models.Airmon_Interface{}
	wlan0 := models.Airmon_Interface{PHY: "phy0", Interface: "wlan0", Driver: "88XXau", Chipset: "Realtek Semiconductor Corp. Realtek 8812AU/8821AU 802.11ac WLAN Adapter [USB Wireless Dual-Band Adapter 2.4/5Ghz]", Mode: "Monitor"}
	wlan1 := models.Airmon_Interface{PHY: "phy1", Interface: "wlan1", Driver: "brcmfmac", Chipset: "Broadcom 43430", Mode: "Managed"}
	airmonInterfaces = append(airmonInterfaces, wlan0)
	airmonInterfaces = append(airmonInterfaces, wlan1)

	airmonResult := models.Airmon_Result{Airmon_Interfaces: airmonInterfaces, Details: "(monitor mode enabled)\n", Result: "success",}
	result := StartStopAirmon(string(cmdOut))
	if !reflect.DeepEqual(result, airmonResult) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", airmonResult, result)
	}
}

func TestStopAirmonOnValid(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng", "stop", "wlan0").Output()
	if err != nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	airmonInterfaces := []models.Airmon_Interface{}
	wlan0 := models.Airmon_Interface{PHY: "phy0", Interface: "wlan0", Driver: "88XXau", Chipset: "Realtek Semiconductor Corp. Realtek 8812AU/8821AU 802.11ac WLAN Adapter [USB Wireless Dual-Band Adapter 2.4/5Ghz]", Mode: "Managed"}
	wlan1 := models.Airmon_Interface{PHY: "phy1", Interface: "wlan1", Driver: "brcmfmac", Chipset: "Broadcom 43430", Mode: "Managed"}
	airmonInterfaces = append(airmonInterfaces, wlan0)
	airmonInterfaces = append(airmonInterfaces, wlan1)

	airmonResult := models.Airmon_Result{Airmon_Interfaces: airmonInterfaces, Details: "(monitor mode disabled)\n", Result: "success",}
	result := StartStopAirmon(string(cmdOut))
	if !reflect.DeepEqual(result, airmonResult) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", airmonResult, result)
	}
}

func TestStartAirmonOnInValid(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng", "start", "wlan1").Output()
	if err == nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	airmonInterfaces := []models.Airmon_Interface{}
	wlan0 := models.Airmon_Interface{PHY: "phy0", Interface: "wlan0", Driver: "88XXau", Chipset: "Realtek Semiconductor Corp. Realtek 8812AU/8821AU 802.11ac WLAN Adapter [USB Wireless Dual-Band Adapter 2.4/5Ghz]", Mode: "Managed"}
	wlan1 := models.Airmon_Interface{PHY: "phy1", Interface: "wlan1", Driver: "brcmfmac", Chipset: "Broadcom 43430", Mode: "Managed"}
	airmonInterfaces = append(airmonInterfaces, wlan0)
	airmonInterfaces = append(airmonInterfaces, wlan1)


	airmonResult := models.Airmon_Result{Airmon_Interfaces: airmonInterfaces, Result: "fail", Details: "ERROR adding monitor mode interface: command failed: Operation not supported (-95)\n"}
	result := StartStopAirmon(string(cmdOut))
	if !reflect.DeepEqual(result, airmonResult) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", airmonResult, result)
	}
}

func TestStopAirmonOnInValid(t *testing.T) {
	// Test Case 1
	cmdOut, err := exec.Command("airmon-ng", "stop", "wlan0").Output()
	if err == nil {
		t.Errorf("Test case 1 failed. Error: %s", err.Error())
	}
	airmonInterfaces := []models.Airmon_Interface{}
	wlan0 := models.Airmon_Interface{PHY: "phy0", Interface: "wlan0", Driver: "88XXau", Chipset: "Realtek Semiconductor Corp. Realtek 8812AU/8821AU 802.11ac WLAN Adapter [USB Wireless Dual-Band Adapter 2.4/5Ghz]", Mode: "Managed"}
	airmonInterfaces = append(airmonInterfaces, wlan0)

	airmonResult := models.Airmon_Result{Airmon_Interfaces: airmonInterfaces, Result: "fail", Details: "You are trying to stop a device that isn't in monitor mode."}
	result := StartStopAirmon(string(cmdOut))
	if !reflect.DeepEqual(result, airmonResult) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", airmonResult, result)
	}
}