package cmd

import (
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
						Username: "username@example.com",
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
		username := cfg.SMTP.Username
		from := cfg.SMTP.From
		password := cfg.SMTP.Password
		to := cfg.SMTP.To

		// If username is not set, use the from address as fallback
		if username == "" {
			username = from
			logger.Info("SMTP username not specified, using From address as username")
		}

		// Message content
		subject := "Test Email from Go SMTP Program"
		body := "This is a test email sent from a Go program using SMTP with gomail"

		// Automatically detect connection type based on port
		// Port 465 typically uses direct TLS
		// Port 587 typically uses STARTTLS
		// User can override with -tls flag
		useDirectTLS := useTLS || smtpPort == 465

		if verbose {
			logger.Info("Sending email with verbose logging...")
			err = smtpClient.SendMailWithTLSDebug(
				smtpHost,
				smtpPort,
				username,
				password,
				from,
				to,
				subject,
				body,
				useDirectTLS,
			)
		} else {
			if useDirectTLS {
				logger.Info("Sending email with direct TLS connection (port %d)...", smtpPort)
			} else {
				logger.Info("Sending email with STARTTLS if supported (port %d)...", smtpPort)
			}

			err = smtpClient.SendMailWithGomail(
				smtpHost,
				smtpPort,
				username,
				password,
				from,
				to,
				subject,
				body,
				useDirectTLS,
				false,
			)
		}

		if err != nil {
			exitWithError("Failed to send email: %v", err)
		}

		logger.Info("Email sent successfully!")
	},
}
