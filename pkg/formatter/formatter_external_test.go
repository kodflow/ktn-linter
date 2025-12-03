// External tests for formatter.go (black-box testing).
package formatter_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
)

// TestDiagnosticGroupData tests the DiagnosticGroupData struct.
func TestDiagnosticGroupData(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"empty filename", ""},
		{"valid filename", "test.go"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			data := formatter.DiagnosticGroupData{
				Filename:    tt.filename,
				Diagnostics: nil,
			}
			// Vérification du nom de fichier
			if data.Filename != tt.filename {
				t.Errorf("Filename = %q, want %q", data.Filename, tt.filename)
			}
		})
	}
}

// TestFormatterConstants tests exported constants.
func TestFormatterConstants(t *testing.T) {
	// Test color codes are not empty
	if formatter.RED == "" {
		t.Error("RED should not be empty")
	}
	// Test INITIAL_FILE_MAP_CAP
	if formatter.INITIAL_FILE_MAP_CAP <= 0 {
		t.Errorf("INITIAL_FILE_MAP_CAP = %d, want > 0", formatter.INITIAL_FILE_MAP_CAP)
	}
}
