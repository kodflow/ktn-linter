// External tests for orchestrator types.
package orchestrator_test

import (
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"golang.org/x/tools/go/analysis"
)

// TestNewDiagnosticResult tests the NewDiagnosticResult function.
func TestNewDiagnosticResult(t *testing.T) {
	tests := []struct {
		name         string
		analyzerName string
		message      string
	}{
		{
			name:         "create diagnostic result",
			analyzerName: "testanalyzer",
			message:      "test message",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f := fset.AddFile("test.go", -1, 100)

			diag := analysis.Diagnostic{
				Pos:     f.Pos(10),
				Message: tt.message,
			}

			result := orchestrator.NewDiagnosticResult(diag, fset, tt.analyzerName)

			// Verify fields
			if result.AnalyzerName != tt.analyzerName {
				t.Errorf("expected analyzer name %s, got %s", tt.analyzerName, result.AnalyzerName)
			}
			// Verify diagnostic
			if result.Diag.Message != tt.message {
				t.Errorf("expected message %s, got %s", tt.message, result.Diag.Message)
			}
			// Verify fset
			if result.Fset != fset {
				t.Error("expected fset to be set")
			}
		})
	}
}

// TestDiagnosticResult_Position tests the Position method.
func TestDiagnosticResult_Position(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns position from diagnostic",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f := fset.AddFile("test.go", -1, 100)

			diag := analysis.Diagnostic{
				Pos:     f.Pos(10),
				Message: "test",
			}

			result := orchestrator.DiagnosticResult{
				Diag:         diag,
				Fset:         fset,
				AnalyzerName: "test",
			}

			pos := result.Position()

			// Verify position is returned
			if pos.Filename != "test.go" {
				t.Errorf("expected filename test.go, got %s", pos.Filename)
			}
			if pos.Offset != 10 {
				t.Errorf("expected offset 10, got %d", pos.Offset)
			}

			// Verify caching works
			pos2 := result.Position()
			if pos.Offset != pos2.Offset {
				t.Error("expected consistent position from cache")
			}
		})
	}
}
