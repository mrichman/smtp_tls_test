package validator

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "test@sub.example.com", false},
		{"Valid email with plus", "test+tag@example.com", false},
		{"Valid email with dots", "test.name@example.com", false},
		{"Empty email", "", true},
		{"Missing @", "testexample.com", true},
		{"Missing domain", "test@", true},
		{"Missing username", "@example.com", true},
		{"Invalid TLD", "test@example.c", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmails(t *testing.T) {
	tests := []struct {
		name    string
		emails  []string
		wantErr bool
	}{
		{"Valid emails", []string{"test1@example.com", "test2@example.com"}, false},
		{"Empty list", []string{}, true},
		{"One invalid email", []string{"test1@example.com", "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmails(tt.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
