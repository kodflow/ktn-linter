// Internal tests for the diagnostics processor.
package orchestrator

import (
	"testing"
)

// TestDiagnosticsProcessor_isModernize tests the isModernize method.
func TestDiagnosticsProcessor_isModernize(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
		want     bool
	}{
		{
			name:     "modernize analyzer any",
			analyzer: "any",
			want:     true,
		},
		{
			name:     "modernize analyzer minmax",
			analyzer: "minmax",
			want:     true,
		},
		{
			name:     "modernize analyzer bloop",
			analyzer: "bloop",
			want:     true,
		},
		{
			name:     "non-modernize analyzer ktnfunc",
			analyzer: "ktnfunc001",
			want:     false,
		},
		{
			name:     "empty analyzer name",
			analyzer: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			p := NewDiagnosticsProcessor()
			got := p.isModernize(tt.analyzer)

			// Verify result
			if got != tt.want {
				t.Errorf("isModernize(%q) = %v, want %v", tt.analyzer, got, tt.want)
			}
		})
	}
}

// TestDiagnosticsProcessor_formatModernizeCode tests the formatModernizeCode method.
func TestDiagnosticsProcessor_formatModernizeCode(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
		want     string
	}{
		{
			name:     "format any",
			analyzer: "any",
			want:     "KTN-MDRNZ-ANY",
		},
		{
			name:     "format minmax",
			analyzer: "minmax",
			want:     "KTN-MDRNZ-MINMAX",
		},
		{
			name:     "format slicescontains",
			analyzer: "slicescontains",
			want:     "KTN-MDRNZ-SLICESCONTAINS",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			p := NewDiagnosticsProcessor()
			got := p.formatModernizeCode(tt.analyzer)

			// Verify result
			if got != tt.want {
				t.Errorf("formatModernizeCode(%q) = %q, want %q", tt.analyzer, got, tt.want)
			}
		})
	}
}
