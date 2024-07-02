package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func generateRandomPassword(length int, config *Config) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		logToFile(config.General.LogFile, fmt.Sprintf("Error generating random password: %v", err))
		return "", err
	}
	password := base64.URLEncoding.EncodeToString(bytes)[:length]
	logToFile(config.General.LogFile, fmt.Sprintf("Generated random password of length %d", length))
	return password, nil
}

func createFileAtLocation(location, content string, config *Config) error {
	dir := filepath.Dir(location)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to create directory %s: %v", dir, err))
			return err
		}
		logToFile(config.General.LogFile, fmt.Sprintf("Created directory %s", dir))
	}
	err := ioutil.WriteFile(location, []byte(content), 0644)
	if err != nil {
		logToFile(config.General.LogFile, fmt.Sprintf("Failed to create file %s: %v", location, err))
	} else {
		logToFile(config.General.LogFile, fmt.Sprintf("Created file at %s", location))
	}
	return err
}

func createLinkFile(location, targetPath string, config *Config) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-Item -Path '%s' -ItemType SymbolicLink -Value '%s'", location, targetPath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		logToFile(config.General.LogFile, fmt.Sprintf("Failed to create link file at %s pointing to %s: %v, output: %s", location, targetPath, err, string(output)))
		return err
	}
	logToFile(config.General.LogFile, fmt.Sprintf("Created link file at %s pointing to %s", location, targetPath))
	return nil
}

func createRegistryKey(location, keyType, value string, config *Config) error {
	cmd := exec.Command("reg", "add", location, "/v", "ExampleValue", "/t", keyType, "/d", value, "/f")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logToFile(config.General.LogFile, fmt.Sprintf("Failed to create registry key at %s with type %s and value %s: %v, output: %s", location, keyType, value, err, string(output)))
		return err
	}
	logToFile(config.General.LogFile, fmt.Sprintf("Created registry key at %s with type %s and value %s", location, keyType, value))
	return nil
}

func createAndModifyFiles(config *Config) {
	for _, fileOp := range config.FileOperations.CreateModifyFiles {
		var content string
		var err error
		if fileOp.UseGPT {
			content = generateContentUsingGPT(config.General.OpenaiApiKey, fileOp.GPTPrompt, config)
		} else {
			content = fileOp.Content
		}
		err = createFileAtLocation(fileOp.Path, content, config)
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to create/modify file %s: %v", fileOp.Path, err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Created/modified file %s", fileOp.Path))
		}
	}
}
