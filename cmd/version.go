package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "1.0.0"
	BuildDate = "2025-04-19"
	GitCommit = "development"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version, build date, and git commit of the SMTP TLS test tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SMTP TLS Test v%s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
	},
}
