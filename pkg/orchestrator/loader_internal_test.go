// Internal tests for the package loader.
package orchestrator

import (
	"bytes"
	"testing"

	"golang.org/x/tools/go/packages"
)

// TestPackageLoader_checkErrors tests the checkErrors method.
func TestPackageLoader_checkErrors(t *testing.T) {
	tests := []struct {
		name      string
		pkgs      []*packages.Package
		expectErr bool
	}{
		{
			name:      "no packages returns nil",
			pkgs:      []*packages.Package{},
			expectErr: false,
		},
		{
			name: "package without errors returns nil",
			pkgs: []*packages.Package{
				{PkgPath: "test/pkg", Errors: nil},
			},
			expectErr: false,
		},
		{
			name: "package with VCS error only returns nil",
			pkgs: []*packages.Package{
				{
					PkgPath: "test/pkg",
					Errors: []packages.Error{
						{Msg: "VCS status error"},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "package with real error returns error",
			pkgs: []*packages.Package{
				{
					PkgPath: "test/pkg",
					Errors: []packages.Error{
						{Msg: "syntax error"},
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			loader := NewPackageLoader(&buf)

			err := loader.checkErrors(tt.pkgs)

			// Verify error expectation
			if tt.expectErr && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
