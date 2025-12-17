// External tests for orchestrator types.
package orchestrator_test

import (
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"golang.org/x/tools/go/analysis"
)

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
