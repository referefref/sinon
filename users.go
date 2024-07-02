package main

import (
	"fmt"

	wapi "github.com/iamacarpet/go-win64api"
)

func manageUserAccounts(config *Config) {
	for _, account := range config.UserAccounts {
		ok, err := wapi.UserAdd(account.Name, account.FullName, account.Password)
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to manage user account %s: %v", account.Name, err))
		} else if !ok {
			logToFile(config.General.LogFile, fmt.Sprintf("User account %s was not added successfully.", account.Name))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Managed user account %s", account.Name))
		}
	}
}
