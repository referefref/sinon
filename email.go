package main

import (
	"fmt"
	"net/smtp"
)

func sendEmails(config *Config) {
	for _, emailOp := range config.EmailOperations.SendReceive {
		var body string
		if emailOp.UseGPT {
			body = generateContentUsingGPT(config.General.OpenaiApiKey, emailOp.GPTPrompt, config)
		} else {
			body = emailOp.Body
		}
		sendEmail(config.EmailOperations.GoogleAccount, emailOp.SendTo, emailOp.Subject, body, config)
	}
}

func sendEmail(account EmailAccount, to, subject, body string, config *Config) {
	from := account.Email
	password := account.Password

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		logToFile(config.General.LogFile, fmt.Sprintf("Failed to send email: %v", err))
	} else {
		logToFile(config.General.LogFile, fmt.Sprintf("Email sent to %s", to))
	}
}
