package models

import (
	"testing"
)

// Test HasPermission
func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		user       *User
		permission string
		expected   bool
	}{
		{
			name: "Admin has all permissions",
			user: &User{
				Role:        "admin",
				Permissions: []string{},
			},
			permission: "any.permission",
			expected:   true,
		},
		{
			name: "Manager has specific permission",
			user: &User{
				Role:        "manager",
				Permissions: []string{"properties.view", "properties.edit"},
			},
			permission: "properties.edit",
			expected:   true,
		},
		{
			name: "Manager does not have permission",
			user: &User{
				Role:        "manager",
				Permissions: []string{"properties.view"},
			},
			permission: "properties.delete",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.HasPermission(tt.permission)
			if result != tt.expected {
				t.Errorf("HasPermission() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test AddPermission
func TestUser_AddPermission(t *testing.T) {
	user := &User{
		Role:        "manager",
		Permissions: []string{"properties.view"},
	}

	user.AddPermission("properties.edit")

	if len(user.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(user.Permissions))
	}

	if !user.HasPermission("properties.edit") {
		t.Error("Expected user to have properties.edit permission")
	}

	// Try adding duplicate
	user.AddPermission("properties.edit")
	if len(user.Permissions) != 2 {
		t.Errorf("Expected 2 permissions after adding duplicate, got %d", len(user.Permissions))
	}
}

// Test RemovePermission
func TestUser_RemovePermission(t *testing.T) {
	user := &User{
		Role:        "manager",
		Permissions: []string{"properties.view", "properties.edit", "properties.delete"},
	}

	user.RemovePermission("properties.edit")

	if len(user.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(user.Permissions))
	}

	if user.HasPermission("properties.edit") {
		t.Error("Expected user to NOT have properties.edit permission")
	}

	if !user.HasPermission("properties.view") {
		t.Error("Expected user to still have properties.view permission")
	}
}

// Test IsValidUserRole
func TestIsValidUserRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"admin", true},
		{"manager", true},
		{"broker", false},
		{"broker_admin", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			result := IsValidUserRole(tt.role)
			if result != tt.expected {
				t.Errorf("IsValidUserRole(%s) = %v, expected %v", tt.role, result, tt.expected)
			}
		})
	}
}

// Test ValidUserRoles
func TestValidUserRoles(t *testing.T) {
	roles := ValidUserRoles()

	if len(roles) != 2 {
		t.Errorf("Expected 2 valid roles, got %d", len(roles))
	}

	expectedRoles := map[string]bool{
		"admin":   true,
		"manager": true,
	}

	for _, role := range roles {
		if !expectedRoles[role] {
			t.Errorf("Unexpected role in ValidUserRoles: %s", role)
		}
	}
}
