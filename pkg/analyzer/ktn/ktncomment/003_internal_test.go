// Internal tests for 003.go - constant comment analyzer.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runComment003 tests the runComment003 function configuration.
// The actual analyzer is tested via analysistest in 003_external_test.go.
func Test_runComment003(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment003 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer003 is properly configured
			if Analyzer003 == nil {
				t.Error("Analyzer003 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer003.Name != "ktncomment003" {
				t.Errorf("Analyzer003.Name = %q, want %q", Analyzer003.Name, "ktncomment003")
			}
			// Check analyzer requires inspect
			if len(Analyzer003.Requires) == 0 {
				t.Error("Analyzer003 should require inspect.Analyzer")
			}
		})
	}
}

// Test_runComment003_ruleDisabled tests behavior when rule is disabled.
func Test_runComment003_ruleDisabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-003": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			const myConst = 1`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment003(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment003 failed: %v", err)
			}

			// Should report no errors when rule disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runComment003_fileExcluded tests behavior when file is excluded.
func Test_runComment003_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-003": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			const myConst = 1`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment003(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment003 failed: %v", err)
			}

			// Should report no errors when file excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}
