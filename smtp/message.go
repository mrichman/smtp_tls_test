package smtp

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Message represents an email message
type Message struct {
	headers map[string]string
	body    string
}

// NewMessage creates a new email message with the given subject and body
func NewMessage(subject, body string) *Message {
	m := &Message{
		headers: make(map[string]string),
		body:    body,
	}
	m.headers["Subject"] = subject
	m.headers["MIME-Version"] = "1.0"
	m.headers["Content-Type"] = "text/plain; charset=UTF-8"
	m.headers["Date"] = time.Now().Format(time.RFC1123Z)
	return m
}

// SetFrom sets the sender of the message
func (m *Message) SetFrom(from string) *Message {
	m.headers["From"] = from
	return m
}

// AddTo adds recipients to the message
func (m *Message) AddTo(to ...string) *Message {
	m.headers["To"] = strings.Join(to, ", ")
	return m
}

// AddHeader adds a custom header to the message
func (m *Message) AddHeader(key, value string) *Message {
	m.headers[key] = value
	return m
}

// Bytes returns the message as a byte slice
func (m *Message) Bytes() []byte {
	var buf bytes.Buffer

	// Write headers
	for k, v := range m.headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
	}

	// Separator between headers and body
	buf.WriteString("\r\n")

	// Write body
	buf.WriteString(m.body)

	return buf.Bytes()
}
