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
	tests := []struct {
		name  string
		check func() bool
		msg   string
	}{
		{name: "RED not empty", check: func() bool { return formatter.RED != "" }, msg: "RED should not be empty"},
		{name: "INITIAL_FILE_MAP_CAP positive", check: func() bool { return formatter.INITIAL_FILE_MAP_CAP > 0 }, msg: "INITIAL_FILE_MAP_CAP should be > 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.msg)
			}
		})
	}
}
