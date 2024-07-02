package main

import (
	"math/rand"
	"os/exec"
	"time"
)

func browseWebsites(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	websites := selectRandomOrHardcoded(config.Websites.Options, config.Websites.SelectionMethod)
	for _, website := range websites {
		cmd := exec.Command("cmd", "/C", "start", website)
		err := cmd.Run()
		if err != nil {
			logger.Printf("Failed to browse %s: %v", website, err)
		} else {
			logger.Printf("Browsing %s", website)
		}
		time.Sleep(time.Duration(config.General.ActionDelay+rand.Intn(2*config.General.RandomnessFactor+1)-config.General.RandomnessFactor) * time.Second)
	}
}
