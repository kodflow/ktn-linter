// External tests for the diagnostics processor.
package orchestrator_test

import (
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"golang.org/x/tools/go/analysis"
)

// TestNewDiagnosticsProcessor tests the NewDiagnosticsProcessor function.
func TestNewDiagnosticsProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create diagnostics processor",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			processor := orchestrator.NewDiagnosticsProcessor()

			// Verify processor created
			if processor == nil {
				t.Error("expected non-nil processor")
			}
		})
	}
}

// TestDiagnosticsProcessor_Filter tests the Filter method.
func TestDiagnosticsProcessor_Filter(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() []orchestrator.DiagnosticResult
		wantLen int
	}{
		{
			name: "filter nil diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				return nil
			},
			wantLen: 0,
		},
		{
			name: "filter empty diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				return []orchestrator.DiagnosticResult{}
			},
			wantLen: 0,
		},
		{
			name: "filter cache file diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				cacheFile := fset.AddFile("/.cache/go-build/test.go", -1, 100)
				return []orchestrator.DiagnosticResult{
					{
						Diag:         analysis.Diagnostic{Pos: cacheFile.Pos(10), Message: "issue"},
						Fset:         fset,
						AnalyzerName: "test",
					},
				}
			},
			wantLen: 0,
		},
		{
			name: "keep normal file diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				normalFile := fset.AddFile("/project/main.go", -1, 100)
				return []orchestrator.DiagnosticResult{
					{
						Diag:         analysis.Diagnostic{Pos: normalFile.Pos(10), Message: "issue"},
						Fset:         fset,
						AnalyzerName: "test",
					},
				}
			},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			processor := orchestrator.NewDiagnosticsProcessor()
			diagnostics := tt.setup()

			filtered := processor.Filter(diagnostics)

			// Verify result length
			if len(filtered) != tt.wantLen {
				t.Errorf("expected %d diagnostics, got %d", tt.wantLen, len(filtered))
			}
		})
	}
}

// TestDiagnosticsProcessor_Extract tests the Extract method.
func TestDiagnosticsProcessor_Extract(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() []orchestrator.DiagnosticResult
		wantLen     int
		checkPrefix string
	}{
		{
			name: "extract nil diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				return nil
			},
			wantLen:     0,
			checkPrefix: "",
		},
		{
			name: "extract empty diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				return []orchestrator.DiagnosticResult{}
			},
			wantLen:     0,
			checkPrefix: "",
		},
		{
			name: "extract single diagnostic",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				file := fset.AddFile("/project/main.go", -1, 100)
				return []orchestrator.DiagnosticResult{
					{
						Diag:         analysis.Diagnostic{Pos: file.Pos(10), Message: "test issue"},
						Fset:         fset,
						AnalyzerName: "ktnfunc001",
					},
				}
			},
			wantLen:     1,
			checkPrefix: "",
		},
		{
			name: "deduplicate identical diagnostics",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				file := fset.AddFile("/project/main.go", -1, 100)
				diag := orchestrator.DiagnosticResult{
					Diag:         analysis.Diagnostic{Pos: file.Pos(10), Message: "test issue"},
					Fset:         fset,
					AnalyzerName: "ktnfunc001",
				}
				return []orchestrator.DiagnosticResult{diag, diag}
			},
			wantLen:     1,
			checkPrefix: "",
		},
		{
			name: "add modernize prefix to analyzer",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				file := fset.AddFile("/project/main.go", -1, 100)
				return []orchestrator.DiagnosticResult{
					{
						Diag:         analysis.Diagnostic{Pos: file.Pos(10), Message: "use any instead"},
						Fset:         fset,
						AnalyzerName: "any",
					},
				}
			},
			wantLen:     1,
			checkPrefix: "KTN-MDRNZ-ANY",
		},
		{
			name: "skip double prefix for KTN messages",
			setup: func() []orchestrator.DiagnosticResult {
				fset := token.NewFileSet()
				file := fset.AddFile("/project/main.go", -1, 100)
				return []orchestrator.DiagnosticResult{
					{
						Diag:         analysis.Diagnostic{Pos: file.Pos(10), Message: "KTN-FUNC-001: function too long"},
						Fset:         fset,
						AnalyzerName: "ktnfunc001",
					},
				}
			},
			wantLen:     1,
			checkPrefix: "KTN-FUNC-001",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			processor := orchestrator.NewDiagnosticsProcessor()
			diagnostics := tt.setup()

			extracted := processor.Extract(diagnostics)

			// Verify result length
			if len(extracted) != tt.wantLen {
				t.Errorf("expected %d diagnostics, got %d", tt.wantLen, len(extracted))
			}

			// Verify prefix if expected
			if tt.checkPrefix != "" && len(extracted) > 0 {
				if !strings.HasPrefix(extracted[0].Message, tt.checkPrefix) {
					t.Errorf("expected prefix %q, got message %q", tt.checkPrefix, extracted[0].Message)
				}
			}
		})
	}
}
