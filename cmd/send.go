package cmd

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/spf13/cobra"

	"smtp_tls_test/config"
	"smtp_tls_test/logger"
	smtpClient "smtp_tls_test/smtp"
	"smtp_tls_test/validator"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a test email",
	Long: `Send a test email using the configured SMTP server.
The email can be sent with or without TLS encryption.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up logger
		logger.SetDefaultVerbose(verbose)

		// Load configuration
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			// If config file doesn't exist, use hardcoded defaults
			if os.IsNotExist(err) {
				logger.Info("Config file not found, using default values")
				cfg = &config.Config{
					SMTP: config.SMTPConfig{
						Host:     "smtp.example.com",
						Port:     587,
						From:     "sender@example.com",
						Password: "your_password",
						To:       []string{"recipient@example.com"},
					},
				}
			} else {
				exitWithError("Failed to load config: %v", err)
			}
		}

		// Validate email addresses
		if err := validator.ValidateEmail(cfg.SMTP.From); err != nil {
			exitWithError("Invalid sender email: %v", err)
		}
		if err := validator.ValidateEmails(cfg.SMTP.To); err != nil {
			exitWithError("Invalid recipient email(s): %v", err)
		}

		// SMTP server configuration
		smtpHost := cfg.SMTP.Host
		smtpPort := cfg.SMTP.Port
		from := cfg.SMTP.From
		password := cfg.SMTP.Password
		to := cfg.SMTP.To

		// Message content
		subject := "Test Email from Go SMTP Program"
		body := "This is a test email sent from a Go program using SMTP"

		// Create message
		message := smtpClient.NewMessage(subject, body).
			SetFrom(from).
			AddTo(to...).
			AddHeader("X-Mailer", "SMTP TLS Test").
			Bytes()

		// Authentication
		auth := smtp.PlainAuth("", from, password, smtpHost)

		// Address in format host:port
		addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

		// Automatically detect connection type based on port
		// Port 465 typically uses direct TLS
		// Port 587 typically uses STARTTLS
		// User can override with -tls flag
		useDirectTLS := useTLS || smtpPort == 465

		if verbose {
			logger.LogSMTPConversation(smtpHost, smtpPort, from, to)
		}

		if useDirectTLS {
			logger.Info("Sending email with direct TLS connection (port %d)...", smtpPort)
			err = smtpClient.SendMailWithTLS(addr, auth, from, to, message, smtpHost, verbose)
		} else {
			logger.Info("Sending email with STARTTLS if supported (port %d)...", smtpPort)
			if verbose {
				err = smtpClient.SendMailWithSTARTTLS(addr, auth, from, to, message, smtpHost)
			} else {
				err = smtp.SendMail(addr, auth, from, to, message)
			}
		}

		if err != nil {
			exitWithError("Failed to send email: %v", err)
		}

		logger.Info("Email sent successfully!")
	},
}
