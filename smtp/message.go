package smtp

import (
	"gopkg.in/gomail.v2"
)

// CreateMessage creates a new email message with gomail
func CreateMessage(from string, to []string, subject string, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.SetHeader("X-Mailer", "SMTP TLS Test")
	return m
}
