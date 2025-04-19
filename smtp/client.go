package smtp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
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

// SendMailWithTLS sends an email using a direct TLS connection
func SendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, host string, verbose bool) error {
	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Note: In production, set this to false
		ServerName:         host,
	}

	// Connect to the server with TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create client, with debug wrapper if verbose mode is enabled
	var client *smtp.Client
	if verbose {
		// Wrap the TLS connection with our debug connection
		debugConn := NewDebugConn(conn)
		client, err = smtp.NewClient(debugConn, host)
	} else {
		client, err = smtp.NewClient(conn, host)
	}

	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Set the sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient: %v", err)
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data connection: %v", err)
	}

	// Send the QUIT command and close the connection
	return client.Quit()
}

// SendMailWithSTARTTLS sends an email using STARTTLS if available
func SendMailWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, host string) error {
	// Connect to the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Wrap the connection with our debug connection
	debugConn := NewDebugConn(conn)

	// Create client
	client, err := smtp.NewClient(debugConn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Check if server supports STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // Note: In production, set this to false
			ServerName:         host,
		}

		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %v", err)
		}

		fmt.Println("STARTTLS negotiation successful")
	} else {
		fmt.Println("Server does not support STARTTLS, continuing with unencrypted connection")
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Set the sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient: %v", err)
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data connection: %v", err)
	}

	// Send the QUIT command and close the connection
	return client.Quit()
}
