package main

import (
	"log"
	"time"
)

func main() {
	configFile := "config.yaml"

	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	if config.General.LogFile != "" {
		logToFile(config.General.LogFile, "Starting interaction")
	}

	// Check if Chocolatey is installed
	if err := checkAndInstallChocolatey(); err != nil {
		logToFile(config.General.LogFile, "Error installing Chocolatey: "+err.Error())
		log.Fatalf("Error installing Chocolatey: %v", err)
	}

	// Perform initial setup
	createAndModifyFiles(config)
	sendEmails(config)
	performSystemUpdates(config)
	installApplications(config)
	manageSoftware(config)
	addStartMenuItems(config)
	manageLures(config)
	browseWebsites(config)
	changePreferences(config)
	performScheduledTasks(config)
	openMediaFiles(config)
	manageUserAccounts(config)
	manageNetworkSettings(config)
	printDocuments(config)
	downloadDecoyFiles(config)

	// Monitor for interactions with lures
	go monitorLures(config)

	duration := config.General.InteractionDuration
	if duration == 0 {
		duration = 60 // Default duration in minutes
	}

	// Keep the application running
	select {
	case <-time.After(time.Duration(duration) * time.Minute):
		if config.General.LogFile != "" {
			logToFile(config.General.LogFile, "Interaction completed")
		}
	}
}
