// Internal tests for 018.go private functions
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

// TestIsSnakeCase teste la fonction isSnakeCase.
func TestIsSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "snake_case avec underscore",
			input:    "my_variable",
			expected: true,
		},
		{
			name:     "camelCase sans underscore",
			input:    "myVariable",
			expected: false,
		},
		{
			name:     "SCREAMING_SNAKE_CASE tout en majuscules",
			input:    "MY_CONSTANT",
			expected: false,
		},
		{
			name:     "nom simple sans underscore",
			input:    "simple",
			expected: false,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isSnakeCase(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSnakeCase(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckVar018Names teste la fonction checkVar018Names.
func TestCheckVar018Names(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name:        "snake_case détecté",
			code:        "var my_variable int",
			expectError: true,
		},
		{
			name:        "camelCase accepté",
			code:        "var myVariable int",
			expectError: false,
		},
		{
			name:        "blank identifier ignoré",
			code:        "var _ int",
			expectError: false,
		},
		{
			name:        "nom simple sans underscore accepté",
			code:        "var simple int",
			expectError: false,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\n"+tt.code, 0)
			// Vérification parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			errorFound := false
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					errorFound = true
				},
			}

			// Parcours des déclarations
			for _, decl := range file.Decls {
				// Vérification type GenDecl
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					// Parcours des specs
					for _, spec := range genDecl.Specs {
						// Vérification type ValueSpec
						if valueSpec, ok := spec.(*ast.ValueSpec); ok {
							checkVar018Names(pass, valueSpec)
						}
					}
				}
			}

			// Vérification résultat
			if errorFound != tt.expectError {
				t.Errorf("checkVar018Names() error = %v, expectError %v", errorFound, tt.expectError)
			}
		})
	}
}

// TestRunVar018 teste la fonction runVar018.
func TestRunVar018(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name: "snake_case dans déclaration var",
			code: `package test
var my_variable int`,
			expectError: true,
		},
		{
			name: "camelCase valide",
			code: `package test
var myVariable int`,
			expectError: false,
		},
		{
			name: "déclaration const ignorée",
			code: `package test
const MY_CONSTANT = 42`,
			expectError: false,
		},
		{
			name: "plusieurs variables dont snake_case",
			code: `package test
var (
	valid int
	my_invalid int
)`,
			expectError: true,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			errorFound := false

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(d analysis.Diagnostic) {
					errorFound = true
				},
			}

			_, err = runVar018(pass)
			// Vérification absence d'erreur d'exécution
			if err != nil {
				t.Fatalf("runVar018() returned error: %v", err)
			}

			// Vérification résultat
			if errorFound != tt.expectError {
				t.Errorf("runVar018() error = %v, expectError %v", errorFound, tt.expectError)
			}
		})
	}
}

// Test_runVar018_disabled tests runVar018 with disabled rule.
func Test_runVar018_disabled(t *testing.T) {
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
					"KTN-VAR-018": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
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

			_, err = runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar018() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar018_fileExcluded tests runVar018 with excluded file.
func Test_runVar018_fileExcluded(t *testing.T) {
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
					"KTN-VAR-018": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
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

			_, err = runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar018() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}
