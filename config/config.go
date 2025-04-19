package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	SMTP SMTPConfig `json:"smtp"`
}

// SMTPConfig represents the SMTP server configuration
type SMTPConfig struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	From     string   `json:"from"`
	Password string   `json:"password"`
	To       []string `json:"to"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			From:     "sender@example.com",
			Password: "your_password",
			To:       []string{"recipient@example.com"},
		},
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// If no path is provided, use the default
	if path == "" {
		path = "config.json"
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
	}

	return &config, nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	// If no path is provided, use the default
	if path == "" {
		path = "config.json"
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Create the config
	config := DefaultConfig()

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
