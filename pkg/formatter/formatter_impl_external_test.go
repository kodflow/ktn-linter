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
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f := formatter.NewFormatter(buf, tt.aiMode, tt.noColor, tt.simpleMode)
			if f == nil {
				t.Error("NewFormatter returned nil")
			}
		})
	}
}

func TestFormat(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false)
	fset := token.NewFileSet()

	// Test avec diagnostics vides
	f.Format(fset, []analysis.Diagnostic{})
	output := buf.String()

	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected success message for empty diagnostics, got: %s", output)
	}
}
