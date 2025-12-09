// Internal tests for 018.go private functions
package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestIsSnakeCase teste la fonction isSnakeCase.
//
// Params:
//   - t: contexte de test
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
		t.Run(tt.name, func(t *testing.T) {
			result := isSnakeCase(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSnakeCase(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSnakeToCamel teste la fonction snakeToCamel.
//
// Params:
//   - t: contexte de test
func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "conversion simple snake_case",
			input:    "my_variable",
			expected: "myVariable",
		},
		{
			name:     "conversion multiple underscores",
			input:    "my_long_variable_name",
			expected: "myLongVariableName",
		},
		{
			name:     "pas d'underscore retourne inchangé",
			input:    "simple",
			expected: "simple",
		},
		{
			name:     "underscore en fin ignoré",
			input:    "my_var_",
			expected: "myVar",
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := snakeToCamel(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("snakeToCamel(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckVar018Names teste la fonction checkVar018Names.
//
// Params:
//   - t: contexte de test
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
//
// Params:
//   - t: contexte de test
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
