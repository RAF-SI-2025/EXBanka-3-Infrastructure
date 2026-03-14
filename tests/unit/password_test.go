package unit_test

import (
	"strings"
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/util"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := util.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}
	if salt1 == "" {
		t.Fatal("GenerateSalt() returned empty string")
	}

	salt2, err := util.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error on second call: %v", err)
	}
	if salt1 == salt2 {
		t.Error("GenerateSalt() returned identical salts on consecutive calls")
	}
}

func TestHashPassword(t *testing.T) {
	salt, err := util.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}

	hash, err := util.HashPassword("TestPass12", salt)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty string")
	}

	// Same password + salt always produces the same hash
	hash2, err := util.HashPassword("TestPass12", salt)
	if err != nil {
		t.Fatalf("HashPassword() error on second call: %v", err)
	}
	if hash != hash2 {
		t.Error("HashPassword() is not deterministic")
	}

	// Different password produces a different hash
	hashOther, err := util.HashPassword("OtherPass12", salt)
	if err != nil {
		t.Fatalf("HashPassword() error with different password: %v", err)
	}
	if hash == hashOther {
		t.Error("HashPassword() produced same hash for different passwords")
	}
}

func TestHashPassword_InvalidSalt(t *testing.T) {
	_, err := util.HashPassword("TestPass12", "not-valid-base64!!!")
	if err == nil {
		t.Error("HashPassword() expected error for invalid salt, got nil")
	}
}

func TestVerifyPassword(t *testing.T) {
	salt, _ := util.GenerateSalt()
	hash, _ := util.HashPassword("TestPass12", salt)

	ok, err := util.VerifyPassword("TestPass12", salt, hash)
	if err != nil {
		t.Fatalf("VerifyPassword() error: %v", err)
	}
	if !ok {
		t.Error("VerifyPassword() returned false for correct password")
	}

	ok, err = util.VerifyPassword("WrongPass12", salt, hash)
	if err != nil {
		t.Fatalf("VerifyPassword() error with wrong password: %v", err)
	}
	if ok {
		t.Error("VerifyPassword() returned true for incorrect password")
	}
}

func TestValidatePasswordPolicy(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{"valid password", "Abcdef12", false, ""},
		{"too short", "Ab1!", true, "8 and 32"},
		{"too long", strings.Repeat("a", 29) + "AA12", true, "8 and 32"},
		{"no digits", "Abcdefgh", true, "2 digits"},
		{"one digit only", "Abcdefg1", true, "2 digits"},
		{"no uppercase", "abcdef12", true, "1 uppercase"},
		{"no lowercase", "ABCDEF12", true, "1 lowercase"},
		{"exactly 8 chars valid", "Abcde123", false, ""},
		{"exactly 32 chars valid", "Abcde123Abcde123Abcde123Abcde123", false, ""},
		{"33 chars too long", "Abcde123Abcde123Abcde123Abcde1234", true, "8 and 32"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidatePasswordPolicy(tt.password)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePasswordPolicy(%q) expected error, got nil", tt.password)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePasswordPolicy(%q) error = %q, want contains %q", tt.password, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePasswordPolicy(%q) unexpected error: %v", tt.password, err)
				}
			}
		})
	}
}
