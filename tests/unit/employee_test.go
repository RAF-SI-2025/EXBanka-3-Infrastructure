package unit_test

import (
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
)

func TestEmployee_IsAdmin_True(t *testing.T) {
	emp := models.Employee{
		Permissions: []models.Permission{
			{Name: "employee.read"},
			{Name: "admin"},
		},
	}
	if !emp.IsAdmin() {
		t.Error("IsAdmin() returned false for employee with admin permission")
	}
}

func TestEmployee_IsAdmin_False(t *testing.T) {
	emp := models.Employee{
		Permissions: []models.Permission{
			{Name: "employee.read"},
			{Name: "employee.create"},
		},
	}
	if emp.IsAdmin() {
		t.Error("IsAdmin() returned true for employee without admin permission")
	}
}

func TestEmployee_IsAdmin_Empty(t *testing.T) {
	emp := models.Employee{
		Permissions: []models.Permission{},
	}
	if emp.IsAdmin() {
		t.Error("IsAdmin() returned true for employee with no permissions")
	}
}

func TestEmployee_PermissionNames_Multiple(t *testing.T) {
	emp := models.Employee{
		Permissions: []models.Permission{
			{Name: "employee.read"},
			{Name: "employee.create"},
			{Name: "admin"},
		},
	}
	names := emp.PermissionNames()
	if len(names) != 3 {
		t.Fatalf("PermissionNames() returned %d names, want 3", len(names))
	}
	expected := []string{"employee.read", "employee.create", "admin"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("PermissionNames()[%d] = %q, want %q", i, names[i], want)
		}
	}
}

func TestEmployee_PermissionNames_Empty(t *testing.T) {
	emp := models.Employee{
		Permissions: []models.Permission{},
	}
	names := emp.PermissionNames()
	if names == nil {
		t.Error("PermissionNames() returned nil, want empty slice")
	}
	if len(names) != 0 {
		t.Errorf("PermissionNames() returned %d names, want 0", len(names))
	}
}
