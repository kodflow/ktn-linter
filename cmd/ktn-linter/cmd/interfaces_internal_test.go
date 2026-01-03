// Internal tests for interfaces.
package cmd

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/updater"
)

// TestFlagGetterInterface tests that flagGetter interface is implemented correctly.
func TestFlagGetterInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "mock implements flagGetter"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mock := newMockFlagGetter()
			mock.boolValues["test"] = true
			mock.stringValues["test"] = "value"

			// Test interface compliance
			var _ flagGetter = mock

			// Test methods
			boolVal, err := mock.GetBool("test")
			// Check error
			if err != nil {
				t.Errorf("GetBool error = %v", err)
			}
			// Check value
			if !boolVal {
				t.Error("GetBool = false, want true")
			}

			strVal, err := mock.GetString("test")
			// Check error
			if err != nil {
				t.Errorf("GetString error = %v", err)
			}
			// Check value
			if strVal != "value" {
				t.Errorf("GetString = %q, want %q", strVal, "value")
			}
		})
	}
}

// TestUpdaterServiceInterface tests that updaterService interface is implemented.
func TestUpdaterServiceInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "mock implements updaterService"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mock := &mockUpdaterService{
				checkInfo: updater.UpdateInfo{
					Available:      true,
					CurrentVersion: "1.0.0",
					LatestVersion:  "2.0.0",
				},
				upgradeInfo: updater.UpdateInfo{
					Available:      true,
					CurrentVersion: "1.0.0",
					LatestVersion:  "2.0.0",
				},
			}

			// Test interface compliance
			var _ updaterService = mock

			// Test CheckForUpdate
			info, err := mock.CheckForUpdate()
			// Check error
			if err != nil {
				t.Errorf("CheckForUpdate error = %v", err)
			}
			// Check info
			if !info.Available {
				t.Error("CheckForUpdate Available = false, want true")
			}

			// Test Upgrade
			info, err = mock.Upgrade()
			// Check error
			if err != nil {
				t.Errorf("Upgrade error = %v", err)
			}
			// Check info
			if !info.Available {
				t.Error("Upgrade Available = false, want true")
			}
		})
	}
}
