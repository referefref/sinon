package main

import (
	"log"
	"io"
	"os"
)

func logToFile(logFile, message string) {
	if logFile == "" {
		return
	}
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
}

func logMultiWriter(logFile string) *log.Logger {
	logFileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFileWriter)
	logger := log.New(multiWriter, "", log.LstdFlags)

	return logger
}
