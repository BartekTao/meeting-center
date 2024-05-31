package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type MailHandler interface {
	Send(recipientEmails []string, subject, content string) error
}

type gmailHandler struct {
	smtpHost    string
	smtpPort    string
	senderEmail string
	auth        smtp.Auth
}

func NewGmailHandler() (MailHandler, error) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		return nil, fmt.Errorf("missing required environment variables: SENDER_EMAIL and SENDER_PASSWORD")
	}

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	return &gmailHandler{
		smtpHost:    smtpHost,
		smtpPort:    smtpPort,
		senderEmail: senderEmail,
		auth:        auth,
	}, nil
}

func (m *gmailHandler) Send(recipientEmails []string, subject, content string) error {
	headers := "From: " + m.senderEmail + "\r\n" +
		"To: " + recipientEmails[0] + "\r\n" +
		"Subject: " + subject + "\r\n\r\n"
	message := []byte(headers + content)

	log.Printf("Attempting to send email to %v", recipientEmails)

	err := smtp.SendMail(m.smtpHost+":"+m.smtpPort, m.auth, m.senderEmail, recipientEmails, message)
	if err != nil {
		return err
	}

	log.Printf("Email sent successfully to: %v", recipientEmails)
	return nil
}
