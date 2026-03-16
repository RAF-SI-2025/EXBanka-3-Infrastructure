package unit_test

import (
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
)

func TestDefaultPermissions_AllHaveNames(t *testing.T) {
	for i, p := range models.DefaultPermissions {
		if p.Name == "" {
			t.Errorf("DefaultPermissions[%d] has empty Name", i)
		}
	}
}

func TestDefaultPermissions_AllHaveSubjectType(t *testing.T) {
	for i, p := range models.DefaultPermissions {
		if p.SubjectType != models.PermissionSubjectEmployee && p.SubjectType != models.PermissionSubjectClient {
			t.Errorf("DefaultPermissions[%d] (%q) has unexpected SubjectType %q", i, p.Name, p.SubjectType)
		}
	}
}

func TestDefaultPermissions_NoDuplicateNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, p := range models.DefaultPermissions {
		if seen[p.Name] {
			t.Errorf("duplicate permission name: %q", p.Name)
		}
		seen[p.Name] = true
	}
}

func TestDefaultPermissions_EmployeeCount(t *testing.T) {
	count := 0
	for _, p := range models.DefaultPermissions {
		if p.SubjectType == models.PermissionSubjectEmployee {
			count++
		}
	}
	if count < 10 {
		t.Errorf("expected at least 10 employee permissions, got %d", count)
	}
}

func TestDefaultPermissions_ClientCount(t *testing.T) {
	count := 0
	for _, p := range models.DefaultPermissions {
		if p.SubjectType == models.PermissionSubjectClient {
			count++
		}
	}
	if count < 5 {
		t.Errorf("expected at least 5 client permissions, got %d", count)
	}
}

func TestPermissionConstants_SubjectTypes(t *testing.T) {
	if models.PermissionSubjectEmployee != "employee" {
		t.Errorf("PermissionSubjectEmployee = %q, want %q", models.PermissionSubjectEmployee, "employee")
	}
	if models.PermissionSubjectClient != "client" {
		t.Errorf("PermissionSubjectClient = %q, want %q", models.PermissionSubjectClient, "client")
	}
}
