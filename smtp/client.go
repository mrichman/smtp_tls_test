package smtp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"

	"gopkg.in/gomail.v2"
)

// DebugConn is a connection wrapper that logs all data sent and received
type DebugConn struct {
	net.Conn
	r io.Reader
	w io.Writer
}

// Read reads data from the connection and logs it
func (c *DebugConn) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	if n > 0 {
		fmt.Fprintf(os.Stderr, "< %s", string(b[:n]))
	}
	return n, err
}

// Write writes data to the connection and logs it
func (c *DebugConn) Write(b []byte) (int, error) {
	fmt.Fprintf(os.Stderr, "> %s", string(b))
	return c.w.Write(b)
}

// NewDebugConn creates a new debug connection
func NewDebugConn(conn net.Conn) *DebugConn {
	return &DebugConn{
		Conn: conn,
		r:    conn,
		w:    conn,
	}
}

// SendMailWithGomail sends an email using gomail with appropriate TLS settings
func SendMailWithGomail(host string, port int, username string, password string, from string, to []string,
	subject string, body string, forceTLS bool, verbose bool) error {

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.SetHeader("X-Mailer", "SMTP TLS Test")

	// Create dialer with authentication
	d := gomail.NewDialer(host, port, username, password)

	// Configure TLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // Note: In production, set this to false
		ServerName:         host,
	}

	// Set SSL/TLS mode based on port and forceTLS flag
	if forceTLS || port == 465 {
		// Use SSL/TLS (implicit TLS)
		d.SSL = true
		if verbose {
			fmt.Println("Using direct TLS connection (implicit TLS)")
		}
	} else {
		// Use STARTTLS (explicit TLS) if available
		d.SSL = false
		if verbose {
			fmt.Println("Using STARTTLS if supported (explicit TLS)")
		}
	}

	if verbose {
		fmt.Printf("Connecting to %s:%d...\n", host, port)
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("From: %s\n", from)
		fmt.Printf("To: %v\n", to)
		fmt.Printf("Subject: %s\n", subject)
	}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendMailWithTLSDebug sends an email with TLS and debug logging
// This is a wrapper around gomail that provides detailed SMTP conversation logging
func SendMailWithTLSDebug(host string, port int, username string, password string, from string, to []string,
	subject string, body string, forceTLS bool) error {

	// Log the SMTP conversation
	fmt.Printf("SMTP Server: %s:%d\n", host, port)
	fmt.Printf("From: %s\n", from)
	fmt.Printf("To: %v\n", to)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("TLS Mode: %s\n", getTLSModeName(port, forceTLS))

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.SetHeader("X-Mailer", "SMTP TLS Test")

	// Create dialer with authentication
	d := gomail.NewDialer(host, port, username, password)

	// Configure TLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // Note: In production, set this to false
		ServerName:         host,
	}

	// Set SSL/TLS mode based on port and forceTLS flag
	if forceTLS || port == 465 {
		d.SSL = true
	} else {
		d.SSL = false
	}

	fmt.Println("\n--- SMTP Conversation Log ---")
	fmt.Println("Note: gomail doesn't provide direct access to SMTP conversation")
	fmt.Println("For full SMTP conversation logging, consider using a network proxy tool")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// Helper function to get TLS mode name
func getTLSModeName(port int, forceTLS bool) string {
	if forceTLS || port == 465 {
		return "Direct TLS (Implicit TLS)"
	}
	return "STARTTLS (Explicit TLS) if supported"
}
