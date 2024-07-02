package main

import (
	"fmt"
	"os/exec"
)

// Opens media files.
func openMediaFiles(config *Config) {
	mediaFiles := selectRandomOrHardcoded(config.MediaFiles.Options, config.MediaFiles.SelectionMethod)
	for _, mediaFile := range mediaFiles {
		cmd := exec.Command("cmd", "/C", "start", mediaFile)
		err := cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to open media file %s: %v", mediaFile, err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Opened media file %s", mediaFile))
		}
	}
}
