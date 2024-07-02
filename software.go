package main

import (
	"fmt"
	"os/exec"
)

func installApplications(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	// Check and install Chocolatey if not installed
	err := checkAndInstallChocolatey(config)
	if err != nil {
		logger.Printf("Failed to install Chocolatey: %v", err)
		return
	}

	apps := selectRandomOrHardcoded(config.Applications.Options, config.Applications.SelectionMethod)
	for _, app := range apps {
		cmd := exec.Command("choco", "install", app, "-y")
		err := cmd.Run()
		if err != nil {
			logger.Printf("Failed to install %s: %v", app, err)
		} else {
			logger.Printf("Installed %s", app)
		}
	}
}

func manageSoftware(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	operations := selectRandomOrHardcoded(config.SoftwareManagement.Options, config.SoftwareManagement.SelectionMethod)
	for _, operation := range operations {
		cmd := exec.Command("choco", operation, "-y")
		err := cmd.Run()
		if err != nil {
			logger.Printf("Failed to perform software management operation %s: %v", operation, err)
		} else {
			logger.Printf("Performed software management operation %s", operation)
		}
	}
}

func addStartMenuItems(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	items := selectRandomOrHardcoded(config.StartMenuItems.Options, config.StartMenuItems.SelectionMethod)
	for _, item := range items {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`$s = New-Object -ComObject WScript.Shell; $shortcut = $s.CreateShortcut("C:\\Users\\Public\\Desktop\\%s.lnk"); $shortcut.TargetPath = "%s"; $shortcut.Save()`, item, item))
		err := cmd.Run()
		if err != nil {
			logger.Printf("Failed to add start menu item %s: %v", item, err)
		} else {
			logger.Printf("Added start menu item %s", item)
		}
	}
}

func checkAndInstallChocolatey(config *Config) error {
	logger := logMultiWriter(config.General.LogFile)
	cmd := exec.Command("choco", "-v")
	err := cmd.Run()
	if err != nil {
		logger.Println("Chocolatey not found, installing...")
		cmd = exec.Command("powershell", "-Command", "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install Chocolatey: %v", err)
		}
		logger.Println("Chocolatey installed successfully")
	}
	return nil
}

func checkAndInstallPSWindowsUpdate(config *Config) error {
	logger := logMultiWriter(config.General.LogFile)
	cmd := exec.Command("powershell", "-Command", "Get-Module -ListAvailable -Name PSWindowsUpdate")
	err := cmd.Run()
	if err != nil {
		logger.Println("PSWindowsUpdate module not found, installing...")
		cmd = exec.Command("powershell", "-Command", "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; Install-Module -Name PSWindowsUpdate -Force")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install PSWindowsUpdate module: %v", err)
		}
		logger.Println("PSWindowsUpdate module installed successfully")
	}
	return nil
}

func performSystemUpdates(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	err := checkAndInstallPSWindowsUpdate(config)
	if err != nil {
		logger.Printf("Failed to install PSWindowsUpdate module: %v", err)
		return
	}

	var cmd *exec.Cmd

	switch config.SystemUpdates.Method {
	case "install_all":
		cmd = exec.Command("powershell", "-Command", "Install-WindowsUpdate -MicrosoftUpdate -AcceptAll -AutoReboot")
	case "install_specific":
		for _, update := range config.SystemUpdates.SpecificUpdates {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Get-WindowsUpdate -Install -KBArticleID %s", update))
			err = cmd.Run()
			if err != nil {
				logger.Printf("Failed to install update %s: %v", update, err)
			} else {
				logger.Printf("Installed update %s", update)
			}
		}
		return
	case "install_random":
		updates := selectRandomOrHardcoded(config.SystemUpdates.SpecificUpdates, config.SystemUpdates.SelectionMethod)
		for _, update := range updates {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Get-WindowsUpdate -Install -KBArticleID %s", update))
			err = cmd.Run()
			if err != nil {
				logger.Printf("Failed to install update %s: %v", update, err)
			} else {
				logger.Printf("Installed update %s", update)
			}
		}
		return
	case "remove_random":
		updates := selectRandomOrHardcoded(config.SystemUpdates.HideUpdates, config.SystemUpdates.SelectionMethod)
		for _, update := range updates {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Remove-WindowsUpdate -KBArticleID %s -NoRestart", update))
			err = cmd.Run()
			if err != nil {
				logger.Printf("Failed to remove update %s: %v", update, err)
			} else {
				logger.Printf("Removed update %s", update)
			}
		}
		return
	case "hide_updates":
		for _, update := range config.SystemUpdates.HideUpdates {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("$HideList = '%s'; Get-WindowsUpdate -KBArticleID $HideList –Hide", update))
			err = cmd.Run()
			if err != nil {
				logger.Printf("Failed to hide update %s: %v", update, err)
			} else {
				logger.Printf("Hidden update %s", update)
			}
		}
		return
	default:
		logger.Printf("Unknown system update method: %s", config.SystemUpdates.Method)
		return
	}

	err = cmd.Run()
	if err != nil {
		logger.Printf("Failed to perform system update: %v", err)
	} else {
		logger.Printf("Performed system update with method %s", config.SystemUpdates.Method)
	}
}
