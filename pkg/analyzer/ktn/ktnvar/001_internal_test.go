package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar001 tests the private runVar001 function.
func Test_runVar001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_isScreamingSnakeCase tests the private isScreamingSnakeCase helper function.
func Test_isScreamingSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		expected bool
	}{
		{
			name:     "screaming snake case",
			varName:  "MAX_SIZE",
			expected: true,
		},
		{
			name:     "screaming snake case with digits",
			varName:  "HTTP_200_OK",
			expected: true,
		},
		{
			name:     "camelCase",
			varName:  "maxSize",
			expected: false,
		},
		{
			name:     "PascalCase",
			varName:  "MaxSize",
			expected: false,
		},
		{
			name:     "single letter",
			varName:  "X",
			expected: false,
		},
		{
			name:     "all uppercase no underscore",
			varName:  "HTTP",
			expected: false,
		},
		{
			name:     "blank identifier",
			varName:  "_",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isScreamingSnakeCase(tt.varName)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isScreamingSnakeCase(%q) = %v, expected %v", tt.varName, result, tt.expected)
			}
		})
	}
}

// Test_runVar001_disabled tests runVar001 with disabled rule.
func Test_runVar001_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-001": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse code with SCREAMING_SNAKE_CASE variable
			code := `package test
			var BAD_VARIABLE int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar001() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar001() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar001_fileExcluded tests runVar001 with excluded file.
func Test_runVar001_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-001": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse code with SCREAMING_SNAKE_CASE variable
			code := `package test
			var BAD_VARIABLE int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar001() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar001() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}
