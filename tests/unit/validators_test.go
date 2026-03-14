package unit_test

import (
	"strings"
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/util"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid simple", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "first.last@bank.org", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
		{"missing TLD", "user@example", true},
		{"double @", "user@@example.com", true},
		{"spaces inside", "us er@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidateEmail(tt.email)
			if tt.wantErr && err == nil {
				t.Errorf("ValidateEmail(%q) expected error, got nil", tt.email)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidateEmail(%q) unexpected error: %v", tt.email, err)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{"non-empty value", "John", "ime", false},
		{"empty string", "", "ime", true},
		{"whitespace only", "   ", "prezime", true},
		{"tab only", "\t", "adresa", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidateRequired(tt.value, tt.fieldName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateRequired(%q, %q) expected error, got nil", tt.value, tt.fieldName)
					return
				}
				if !strings.Contains(err.Error(), tt.fieldName) {
					t.Errorf("ValidateRequired(%q, %q) error = %q, want to contain field name", tt.value, tt.fieldName, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("ValidateRequired(%q, %q) unexpected error: %v", tt.value, tt.fieldName, err)
				}
			}
		})
	}
}
