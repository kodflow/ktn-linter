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
		{
			name:    "prerelease version",
			version: "v1.0.0-beta.1",
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

// TestUpdater_CheckForUpdate tests that dev builds cannot check for updates.
func TestUpdater_CheckForUpdate(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantErr     bool
		wantCurrent string
	}{
		{
			name:        "empty version returns error",
			version:     "",
			wantErr:     true,
			wantCurrent: "",
		},
		{
			name:        "dev version returns error",
			version:     "dev",
			wantErr:     true,
			wantCurrent: "dev",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := updater.NewUpdater(tt.version)
			info, err := u.CheckForUpdate()
			// Check error expectation
			if tt.wantErr {
				// Expect error for dev builds
				if err == nil {
					t.Error("CheckForUpdate() should return error for dev build")
				}
				// Skip remaining checks on expected error
				return
			}
			// Expect no error
			if err != nil {
				t.Errorf("CheckForUpdate() unexpected error = %v", err)
				return
			}
			// Check current version is set
			if info.CurrentVersion != tt.wantCurrent {
				t.Errorf("CurrentVersion = %q, want %q", info.CurrentVersion, tt.wantCurrent)
			}
		})
	}
}

// TestUpdater_Upgrade tests that dev builds cannot upgrade.
func TestUpdater_Upgrade(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantErr     bool
		wantCurrent string
	}{
		{
			name:        "empty version returns error",
			version:     "",
			wantErr:     true,
			wantCurrent: "",
		},
		{
			name:        "dev version returns error",
			version:     "dev",
			wantErr:     true,
			wantCurrent: "dev",
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			u := updater.NewUpdater(tt.version)
			info, err := u.Upgrade()

			// Check error expectation
			if tt.wantErr {
				// Expect error for dev builds
				if err == nil {
					t.Error("Upgrade() should return error for dev build")
				}
				// Avoid nil deref on error paths
				return
			}

			// Expect no error
			if err != nil {
				t.Errorf("Upgrade() unexpected error = %v", err)
				return
			}

			// Check current version is set
			if info.CurrentVersion != tt.wantCurrent {
				t.Errorf("CurrentVersion = %q, want %q", info.CurrentVersion, tt.wantCurrent)
			}
		})
	}
}
