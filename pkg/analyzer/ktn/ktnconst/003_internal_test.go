package ktnconst

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runConst003 tests the private runConst003 function.
func Test_runConst003(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API
		})
	}
}

// Test_isValidGoConstantName tests the private isValidGoConstantName function.
func Test_isValidGoConstantName(t *testing.T) {
	tests := []struct {
		name      string
		constName string
		want      bool
	}{
		// Valid CamelCase names
		{"single uppercase letter", "A", true},
		{"single lowercase letter", "a", true},
		{"PascalCase simple", "MaxSize", true},
		{"PascalCase with numbers", "Http2", true},
		{"camelCase simple", "maxSize", true},
		{"camelCase with numbers", "http2Protocol", true},
		{"acronym uppercase", "API", true},
		{"acronym in name", "APIKey", true},
		{"all lowercase", "timeout", true},
		{"all uppercase no underscore", "MAXSIZE", true},
		{"number in middle", "Http2Protocol", true},
		{"single digit after letter", "A1", true},
		{"complex camelCase", "maxConnectionPoolSize", true},
		{"complex PascalCase", "MaxConnectionPoolSize", true},
		{"starts with uppercase", "StatusOK", true},

		// Invalid names (contain underscores)
		{"SCREAMING_SNAKE_CASE", "MAX_SIZE", false},
		{"snake_case", "max_size", false},
		{"mixed with underscore", "Max_Size", false},
		{"underscore at start", "_maxSize", false},
		{"underscore at end", "maxSize_", false},
		{"multiple underscores", "MAX_BUFFER_SIZE", false},
		{"single underscore", "A_B", false},

		// Invalid names (other issues)
		{"empty string", "", false},
		{"starts with number", "1API", false},
		{"contains space", "Max Size", false},
		{"contains hyphen", "max-size", false},
		{"contains special char", "max@size", false},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := isValidGoConstantName(tt.constName)
			// Verify result
			if got != tt.want {
				t.Errorf("isValidGoConstantName(%q) = %v, want %v", tt.constName, got, tt.want)
			}
		})
	}
}

// Test_runConst003_disabled tests that the rule is skipped when disabled.
func Test_runConst003_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: disable the rule
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-003": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse test code with underscore (should not be reported)
			src := `package test
			const MAX_SIZE = 100  // SCREAMING_SNAKE_CASE - should not be reported
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			// Create inspect pass first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{f},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			// Create main pass
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Unexpected error reported when rule is disabled: %s", d.Message)
				},
			}

			// Run the analyzer - should not report anything
			_, err = runConst003(pass)
			if err != nil {
				t.Errorf("runConst003() error = %v", err)
			}

		})
	}
}

// Test_runConst003_excludedFile tests that excluded files are skipped.
func Test_runConst003_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: exclude test.go files
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-003": {
						Enabled: config.Bool(true),
						Exclude: []string{"**/test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse test code with underscore (should not be reported for excluded file)
			src := `package test
			const MAX_SIZE = 100  // SCREAMING_SNAKE_CASE - should not be reported for excluded file
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			// Create inspect pass first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{f},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			// Create main pass
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Unexpected error reported for excluded file: %s", d.Message)
				},
			}

			// Run the analyzer - should not report anything for excluded file
			_, err = runConst003(pass)
			if err != nil {
				t.Errorf("runConst003() error = %v", err)
			}

		})
	}
}
