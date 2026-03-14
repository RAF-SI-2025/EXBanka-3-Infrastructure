package unit_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/util"
)

const testSecret = "test-jwt-secret-key"

func TestGenerateAccessToken(t *testing.T) {
	token, err := util.GenerateAccessToken(1, "test@bank.com", "testuser", []string{"employee.read"}, testSecret, 15)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateAccessToken() returned empty token")
	}
}

func TestParseToken_ValidAccessToken(t *testing.T) {
	perms := []string{"employee.read", "employee.create"}
	token, err := util.GenerateAccessToken(42, "emp@bank.com", "emp01", perms, testSecret, 15)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error: %v", err)
	}

	claims, err := util.ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("ParseToken() error: %v", err)
	}

	if claims.EmployeeID != 42 {
		t.Errorf("claims.EmployeeID = %d, want 42", claims.EmployeeID)
	}
	if claims.Email != "emp@bank.com" {
		t.Errorf("claims.Email = %q, want %q", claims.Email, "emp@bank.com")
	}
	if claims.Username != "emp01" {
		t.Errorf("claims.Username = %q, want %q", claims.Username, "emp01")
	}
	if claims.TokenType != "access" {
		t.Errorf("claims.TokenType = %q, want %q", claims.TokenType, "access")
	}
	if len(claims.Permissions) != 2 {
		t.Errorf("len(claims.Permissions) = %d, want 2", len(claims.Permissions))
	}
}

func TestParseToken_ValidRefreshToken(t *testing.T) {
	token, err := util.GenerateRefreshToken(7, "ref@bank.com", "refuser", testSecret, 24)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error: %v", err)
	}

	claims, err := util.ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("ParseToken() error: %v", err)
	}

	if claims.EmployeeID != 7 {
		t.Errorf("claims.EmployeeID = %d, want 7", claims.EmployeeID)
	}
	if claims.TokenType != "refresh" {
		t.Errorf("claims.TokenType = %q, want %q", claims.TokenType, "refresh")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, _ := util.GenerateAccessToken(1, "a@b.com", "user", nil, testSecret, 15)
	_, err := util.ParseToken(token, "wrong-secret")
	if err == nil {
		t.Error("ParseToken() expected error for wrong secret, got nil")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// Build an already-expired token manually
	claims := util.Claims{
		EmployeeID: 1,
		Email:      "a@b.com",
		Username:   "user",
		TokenType:  "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = util.ParseToken(token, testSecret)
	if err == nil {
		t.Error("ParseToken() expected error for expired token, got nil")
	}
}

func TestParseToken_MalformedToken(t *testing.T) {
	_, err := util.ParseToken("this.is.not.a.jwt", testSecret)
	if err == nil {
		t.Error("ParseToken() expected error for malformed token, got nil")
	}
}

func TestHasPermission(t *testing.T) {
	claims := &util.Claims{
		Permissions: []string{"employee.read", "employee.create"},
	}

	if !util.HasPermission(claims, "employee.read") {
		t.Error("HasPermission() returned false for present permission")
	}
	if !util.HasPermission(claims, "employee.create") {
		t.Error("HasPermission() returned false for present permission")
	}
	if util.HasPermission(claims, "employee.update") {
		t.Error("HasPermission() returned true for absent permission")
	}
	if util.HasPermission(claims, "admin") {
		t.Error("HasPermission() returned true for absent admin permission")
	}
}

func TestHasPermission_AdminBypass(t *testing.T) {
	claims := &util.Claims{
		Permissions: []string{"admin"},
	}
	// admin is in permissions — HasPermission("admin") is true
	if !util.HasPermission(claims, "admin") {
		t.Error("HasPermission() returned false for admin")
	}
}
