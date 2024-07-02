package main

import (
	"github.com/fsnotify/fsnotify"
)

// Monitor interactions with lures.
func monitorLures(config *Config) {
	logger := logMultiWriter(config.General.LogFile)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatalf("Error creating file watcher: %v", err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Println("event:", event)
				if event.Has(fsnotify.Write) {
					logger.Println("modified file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					logger.Println("created file:", event.Name)
				}
				if event.Has(fsnotify.Remove) {
					logger.Println("deleted file:", event.Name)
				}
				if event.Has(fsnotify.Rename) {
					logger.Println("renamed file:", event.Name)
				}
				if event.Has(fsnotify.Chmod) {
					logger.Println("chmod file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Println("error:", err)
			}
		}
	}()

	// Add paths to watch.
	for _, lure := range config.Lures {
		if lure.Type != "registry_key" && lure.Type != "ssh_key" {
			err = watcher.Add(lure.Location)
			if err != nil {
				logger.Fatalf("Error adding watcher to path %s: %v", lure.Location, err)
			}
		}
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
