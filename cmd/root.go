package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	useTLS     bool
	verbose    bool
	configPath string

	rootCmd = &cobra.Command{
		Use:   "smtp_tls_test",
		Short: "A simple SMTP TLS testing tool",
		Long: `SMTP TLS Test is a command line tool for testing SMTP servers with different TLS configurations.
It supports both direct TLS connections and STARTTLS, and can be used to debug SMTP communication.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand is provided, run the send command
			if len(args) == 0 {
				sendCmd.Run(cmd, args)
			}
		},
	}
)

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&useTLS, "tls", "t", false, "Force TLS for SMTP connection")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging of SMTP conversation")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to config file (default: config.json)")

	// Add subcommands
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

// exitWithError prints an error message and exits with code 1
func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
