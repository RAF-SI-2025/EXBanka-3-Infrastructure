package unit_test

import (
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
)

func TestClient_PermissionNames_Multiple(t *testing.T) {
	client := models.Client{
		Permissions: []models.Permission{
			{Name: "client.basic"},
			{Name: "client.trading"},
			{Name: "client.bank_operations"},
		},
	}
	names := client.PermissionNames()
	if len(names) != 3 {
		t.Fatalf("PermissionNames() returned %d names, want 3", len(names))
	}
	expected := []string{"client.basic", "client.trading", "client.bank_operations"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("PermissionNames()[%d] = %q, want %q", i, names[i], want)
		}
	}
}

func TestClient_PermissionNames_Empty(t *testing.T) {
	client := models.Client{
		Permissions: []models.Permission{},
	}
	names := client.PermissionNames()
	if names == nil {
		t.Error("PermissionNames() returned nil, want empty slice")
	}
	if len(names) != 0 {
		t.Errorf("PermissionNames() returned %d names, want 0", len(names))
	}
}

func TestClient_PermissionNames_Order(t *testing.T) {
	client := models.Client{
		Permissions: []models.Permission{
			{Name: "client.bank_operations"},
			{Name: "client.stock_trading"},
			{Name: "client.fund_investing"},
			{Name: "client.basic"},
		},
	}
	names := client.PermissionNames()
	expected := []string{"client.bank_operations", "client.stock_trading", "client.fund_investing", "client.basic"}
	if len(names) != len(expected) {
		t.Fatalf("PermissionNames() returned %d names, want %d", len(names), len(expected))
	}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("PermissionNames()[%d] = %q, want %q (order must match Permissions slice order)", i, names[i], want)
		}
	}
}
