// External tests for the package loader.
package orchestrator_test

import (
	"bytes"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// TestNewPackageLoader tests the NewPackageLoader function.
func TestNewPackageLoader(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create package loader",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			loader := orchestrator.NewPackageLoader(&buf)

			// Verify loader created
			if loader == nil {
				t.Error("expected non-nil loader")
			}
		})
	}
}

// TestPackageLoader_Load tests the Load method.
func TestPackageLoader_Load(t *testing.T) {
	tests := []struct {
		name        string
		patterns    []string
		expectError bool
		minPackages int
	}{
		{
			name:        "load valid package",
			patterns:    []string{"../../pkg/formatter"},
			expectError: false,
			minPackages: 1,
		},
		{
			name:        "load invalid pattern returns error",
			patterns:    []string{"./nonexistent/package"},
			expectError: true,
			minPackages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			loader := orchestrator.NewPackageLoader(&buf)

			pkgs, err := loader.Load(tt.patterns)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify package count
			if !tt.expectError && len(pkgs) < tt.minPackages {
				t.Errorf("expected at least %d packages, got %d", tt.minPackages, len(pkgs))
			}
		})
	}
}
