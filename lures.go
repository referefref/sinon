package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os/exec"
	"strconv"
)

func manageLures(config *Config) {
	for _, lure := range config.Lures {
		var content string
		var username string
		switch lure.GenerativeType {
		case "golang":
			switch lure.Type {
			case "credential_pair":
				if len(config.General.Usernames) > 0 {
					username = selectRandomOrHardcoded(config.General.Usernames, config.General.SelectionMethod)[0]
				} else {
					logToFile(config.General.LogFile, "Error: No usernames provided for credential_pair.")
					continue
				}
				lengthVal, ok := lure.GenerationParams["length"]
				if !ok {
					logToFile(config.General.LogFile, "Error: 'length' parameter is required for credential_pair.")
					continue
				}
				length, err := strconv.Atoi(fmt.Sprintf("%v", lengthVal))
				if err != nil {
					logToFile(config.General.LogFile, fmt.Sprintf("Error: Invalid 'length' parameter: %v", err))
					continue
				}
				password, err := generateRandomPassword(length, config)
				if err != nil {
					logToFile(config.General.LogFile, fmt.Sprintf("Error generating password: %v", err))
					continue
				}
				content = fmt.Sprintf("Username: %s\nPassword: %s", username, password)
			case "ssh_key":
				cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", lure.Location, "-N", "")
				err := cmd.Run()
				if err != nil {
					logToFile(config.General.LogFile, fmt.Sprintf("Error generating SSH key: %v", err))
					continue
				}
				content = fmt.Sprintf("Generated SSH Key at %s", lure.Location)
			case "website_url":
				baseURL, ok := lure.GenerationParams["base_url"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'base_url' parameter is required for website_url.")
					continue
				}
				content = baseURL
			case "registry_key":
				registryKeyType, ok := lure.GenerationParams["registry_key_type"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'registry_key_type' parameter is required for registry_key.")
					continue
				}
				registryKeyValue, ok := lure.GenerationParams["registry_key_value"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'registry_key_value' parameter is required for registry_key.")
					continue
				}
				err := createRegistryKey(lure.Location, registryKeyType, registryKeyValue, config)
				if err != nil {
					logToFile(config.General.LogFile, fmt.Sprintf("Error creating registry key: %v", err))
					continue
				}
				content = fmt.Sprintf("Created registry key at %s with type %s and value %s", lure.Location, registryKeyType, registryKeyValue)
			case "csv":
				documentContent, ok := lure.GenerationParams["document_content"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'document_content' parameter is required for csv.")
					continue
				}
				content = documentContent
			case "api_key":
				apiKeyFormat, ok := lure.GenerationParams["api_key_format"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'api_key_format' parameter is required for api_key.")
					continue
				}
				if apiKeyFormat == "uuid" {
					apiKey := make([]byte, 16)
					_, err := rand.Read(apiKey)
					if err != nil {
						logToFile(config.General.LogFile, fmt.Sprintf("Error generating API key: %v", err))
						continue
					}
					content = fmt.Sprintf("Generated API Key: %x", apiKey)
				} else {
					logToFile(config.General.LogFile, "Unsupported API key format")
					continue
				}
			case "lnk":
				targetPath, ok := lure.GenerationParams["target_path"].(string)
				if !ok {
					logToFile(config.General.LogFile, "Error: 'target_path' parameter is required for lnk.")
					continue
				}
				err := createLinkFile(lure.Location, targetPath, config)
				if err != nil {
					logToFile(config.General.LogFile, fmt.Sprintf("Error creating link file at location: %s, error: %v", lure.Location, err))
					continue
				}
				content = fmt.Sprintf("Created link file at %s pointing to %s", lure.Location, targetPath)
			default:
				logToFile(config.General.LogFile, "Unsupported lure type for Go random generation.")
				continue
			}
		case "openai":
			if config.General.OpenaiApiKey == "" {
				logToFile(config.General.LogFile, "Error: OpenAI API key is required for OpenAI generation.")
				continue
			}
			content = generateContentUsingGPT(config.General.OpenaiApiKey, lure.OpenaiPrompt, config)
		default:
			logToFile(config.General.LogFile, "Unsupported generative type.")
			continue
		}

		if lure.Type != "lnk" && lure.Type != "ssh_key" && lure.Type != "registry_key" {
			if err := createFileAtLocation(lure.Location, content, config); err != nil {
				logToFile(config.General.LogFile, fmt.Sprintf("Error creating file at location: %s, error: %v", lure.Location, err))
				continue
			}
		}

		sourceIP, err := getSourceIP()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Error getting source IP: %v", err))
			continue
		}

		if err := sendMetadataToRedis(config, username, content, sourceIP); err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Error sending metadata to Redis: %v", err))
			continue
		}

		logToFile(config.General.LogFile, fmt.Sprintf("Generated lure: %s", content))
		logToFile(config.General.LogFile, fmt.Sprintf("Lure created and metadata processed successfully."))
	}
}

func getSourceIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot find local IP address")
}
