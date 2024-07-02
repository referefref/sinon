package main

import (
	"fmt"
	"os/exec"
)

// Prints documents.
func printDocuments(config *Config) {
	documents := selectRandomOrHardcoded(config.Printing.Options, config.Printing.SelectionMethod)
	for _, document := range documents {
		cmd := exec.Command("powershell", "-Command", "Start-Process", "notepad.exe", document, "-Verb", "Print")
		err := cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to print document %s: %v", document, err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Printed document %s", document))
		}
	}
}
