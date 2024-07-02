package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Message structure for the OpenAI chat format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse structure for parsing the response from OpenAI API
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func generateContentUsingGPT(apiKey string, prompt string, config *Config) string {
	logToFile(config.General.LogFile, "This is where the prompt goes from the config file input: "+prompt)

	// Prepare the messages payload for the request
	messages := []Message{
		{"system", "You will respond exactly as requested with real information, do not wrap codeblocks like ```yaml, only return the raw yaml file."},
		{"user", prompt},
	}

	// Prepare the JSON payload for the request
	payload := map[string]interface{}{
		"model":       "gpt-4-turbo-preview",
		"messages":    messages,
		"max_tokens":  2048,
		"temperature": 0.3,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logToFile(config.General.LogFile, "Error marshalling payload: "+err.Error())
		return ""
	}

	// Create a new request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(payloadBytes))
	if err != nil {
		logToFile(config.General.LogFile, "Error creating request: "+err.Error())
		return ""
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logToFile(config.General.LogFile, "Error making request to OpenAI: "+err.Error())
		return ""
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logToFile(config.General.LogFile, "Error reading response body: "+err.Error())
		return ""
	}

	// Parse the response
	var chatResponse ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResponse); err != nil {
		logToFile(config.General.LogFile, "Error parsing response JSON: "+err.Error())
		return ""
	}

	// Extract generated content
	var generatedContent string
	for _, choice := range chatResponse.Choices {
		if choice.Message.Role == "assistant" {
			generatedContent = choice.Message.Content
			logToFile(config.General.LogFile, "Generated Content: "+generatedContent)
			break
		}
	}

	if generatedContent == "" {
		logToFile(config.General.LogFile, "No choices returned from OpenAI API")
	}

	return generatedContent
}

