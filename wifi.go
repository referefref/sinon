package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// Manages network settings.
func manageNetworkSettings(config *Config) {
	for _, network := range config.NetworkSettings {
		xmlContent := fmt.Sprintf(`
<?xml version="1.0"?>
<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
    <name>%s</name>
    <SSIDConfig>
        <SSID>
            <name>%s</name>
        </SSID>
    </SSIDConfig>
    <connectionType>ESS</connectionType>
    <connectionMode>auto</connectionMode>
    <MSM>
        <security>
            <authEncryption>
                <authentication>WPA2PSK</authentication>
                <encryption>AES</encryption>
                <useOneX>false</useOneX>
            </authEncryption>
            <sharedKey>
                <keyType>passPhrase</keyType>
                <protected>false</protected>
                <keyMaterial>%s</keyMaterial>
            </sharedKey>
        </security>
    </MSM>
    <MacRandomization xmlns="http://www.microsoft.com/networking/WLAN/profile/v3">
        <enableRandomization>false</enableRandomization>
    </MacRandomization>
</WLANProfile>`, network.SSID, network.SSID, network.Password)

		xmlFilePath := fmt.Sprintf("%s_profile.xml", network.SSID)
		err := ioutil.WriteFile(xmlFilePath, []byte(xmlContent), 0644)
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to write XML file for network %s: %v", network.SSID, err))
			continue
		}

		cmd := exec.Command("netsh", "wlan", "add", "profile", fmt.Sprintf("filename=%s", xmlFilePath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to add profile for network %s: %v, output: %s", network.SSID, err, string(output)))
			continue
		}

		cmd = exec.Command("netsh", "wlan", "connect", fmt.Sprintf("name=%s", network.SSID))
		output, err = cmd.CombinedOutput()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to connect to network %s: %v, output: %s", network.SSID, err, string(output)))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Connected to network %s", network.SSID))
		}
	}
}
