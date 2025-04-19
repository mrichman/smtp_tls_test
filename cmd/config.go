package cmd

import (
	"github.com/spf13/cobra"

	"smtp_tls_test/config"
	"smtp_tls_test/logger"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Create, validate, or display configuration for the SMTP TLS test tool.`,
}

var createConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a default configuration file",
	Long:  `Create a default configuration file with example values.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.CreateDefaultConfig(configPath); err != nil {
			exitWithError("Failed to create config file: %v", err)
		}
		logger.Info("Default config file created successfully")
	},
}

var validateConfigCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Check if the configuration file is valid.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			exitWithError("Failed to load config: %v", err)
		}
		logger.Info("Configuration file is valid")
		logger.Info("SMTP Host: %s", cfg.SMTP.Host)
		logger.Info("SMTP Port: %d", cfg.SMTP.Port)
		logger.Info("From: %s", cfg.SMTP.From)
		logger.Info("To: %v", cfg.SMTP.To)
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  `Display the current configuration values.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			exitWithError("Failed to load config: %v", err)
		}
		logger.Info("Current Configuration:")
		logger.Info("SMTP Host: %s", cfg.SMTP.Host)
		logger.Info("SMTP Port: %d", cfg.SMTP.Port)
		logger.Info("From: %s", cfg.SMTP.From)
		logger.Info("To: %v", cfg.SMTP.To)
		logger.Info("Password: ********")
	},
}

func init() {
	configCmd.AddCommand(createConfigCmd)
	configCmd.AddCommand(validateConfigCmd)
	configCmd.AddCommand(showConfigCmd)
}
