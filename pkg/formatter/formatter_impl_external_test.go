package formatter_test

import (
	"bytes"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
	"golang.org/x/tools/go/analysis"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name       string
		aiMode     bool
		noColor    bool
		simpleMode bool
	}{
		{"default mode", false, false, false},
		{"AI mode", true, false, false},
		{"no color", false, true, false},
		{"simple mode", false, false, true},
		{"all flags", true, true, true},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f := formatter.NewFormatter(buf, tt.aiMode, tt.noColor, tt.simpleMode, false)
			if f == nil {
				t.Error("NewFormatter returned nil")
			}
		})
	}
}

func TestFormatterImpl_Format(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []analysis.Diagnostic
		wantSuccess bool
	}{
		{
			name:        "empty diagnostics shows success",
			diagnostics: []analysis.Diagnostic{},
			wantSuccess: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := formatter.NewFormatter(&buf, false, true, false, false)
			fset := token.NewFileSet()

			// Execute format
			f.Format(fset, tt.diagnostics)
			output := buf.String()

			// Vérification résultat
			if tt.wantSuccess && !strings.Contains(output, "No issues found") {
				t.Errorf("Expected success message, got: %s", output)
			}
		})
	}
}
