package logger

import (
	"fmt"
	"os"
	"time"
)

var defaultVerbose bool

// SetDefaultVerbose sets the default verbose flag
func SetDefaultVerbose(verbose bool) {
	defaultVerbose = verbose
}

// Info logs an informational message
func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Debug logs a debug message if verbose mode is enabled
func Debug(format string, args ...interface{}) {
	if defaultVerbose {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
}

// LogSMTPConversation logs the SMTP conversation details
func LogSMTPConversation(host string, port int, from string, to []string) {
	fmt.Println("=== SMTP Conversation Log ===")
	fmt.Printf("Time: %s\n", time.Now().Format(time.RFC1123))
	fmt.Printf("Server: %s:%d\n", host, port)
	fmt.Printf("From: %s\n", from)
	for i, recipient := range to {
		fmt.Printf("To[%d]: %s\n", i, recipient)
	}
	fmt.Println("===========================")
}
