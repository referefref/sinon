package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

func performScheduledTasks(config *Config) {
	tasks := selectRandomOrHardcodedScheduledTasks(config.ScheduledTasks.Options, config.ScheduledTasks.SelectionMethod)
	for _, task := range tasks {
		cmd := exec.Command("schtasks", "/create", "/tn", task.Name, "/tr", task.Path, "/sc", task.Schedule, "/st", task.StartTime)
		err := cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to create scheduled task %s: %v", task.Name, err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Created scheduled task %s", task.Name))
		}
		time.Sleep(time.Duration(config.General.ActionDelay+rand.Intn(2*config.General.RandomnessFactor+1)-config.General.RandomnessFactor) * time.Second)
	}
}

func selectRandomOrHardcodedScheduledTasks(options []ScheduledTask, method string) []ScheduledTask {
	if method == "hardcoded" {
		return options
	}
	rand.Seed(time.Now().UnixNano())
	return []ScheduledTask{options[rand.Intn(len(options))]}
}
