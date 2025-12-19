// Package updater_test provides black-box tests for the updater package.
package updater_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/updater"
)

// TestNewUpdater tests updater creation with various versions.
func TestNewUpdater(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "with version",
			version: "v1.0.0",
		},
		{
			name:    "dev version",
			version: "dev",
		},
		{
			name:    "empty version",
			version: "",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := updater.NewUpdater(tt.version)
			// Check updater is not nil
			if u == nil {
				t.Error("NewUpdater() returned nil")
			}
		})
	}
}

// TestCheckForUpdateDevBuild tests that dev builds cannot check for updates.
func TestCheckForUpdateDevBuild(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "empty version",
			version: "",
		},
		{
			name:    "dev version",
			version: "dev",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := updater.NewUpdater(tt.version)
			info, err := u.CheckForUpdate()
			// Check error is returned for dev builds
			if err == nil {
				t.Error("CheckForUpdate() should return error for dev build")
			}
			// Check current version is set
			if info.CurrentVersion != tt.version {
				t.Errorf("CurrentVersion = %q, want %q", info.CurrentVersion, tt.version)
			}
		})
	}
}

// TestUpgradeDevBuild tests that dev builds cannot upgrade.
func TestUpgradeDevBuild(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "empty version",
			version: "",
		},
		{
			name:    "dev version",
			version: "dev",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := updater.NewUpdater(tt.version)
			info, err := u.Upgrade()
			// Check error is returned for dev builds
			if err == nil {
				t.Error("Upgrade() should return error for dev build")
			}
			// Check current version is set
			if info.CurrentVersion != tt.version {
				t.Errorf("CurrentVersion = %q, want %q", info.CurrentVersion, tt.version)
			}
		})
	}
}
