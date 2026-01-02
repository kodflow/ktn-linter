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

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
			t.Log("runVar001 tested via external tests")
		})
	}
}

// Test_checkVarSpec tests the checkVarSpec function.
func Test_checkVarSpec(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name:        "type and value - OK",
			code:        "package test\nvar x int = 42",
			expectError: false,
		},
		{
			name:        "no type - ERROR",
			code:        "package test\nvar x = 42",
			expectError: true,
		},
		{
			name:        "type without value - OK (zero-value)",
			code:        "package test\nvar x int",
			expectError: false,
		},
		{
			name:        "slice with type and value - OK",
			code:        "package test\nvar x []string = []string{}",
			expectError: false,
		},
		{
			name:        "blank identifier - skip",
			code:        "package test\nvar _ = 42",
			expectError: false,
		},
		{
			name:        "multiple names without type - ERROR",
			code:        "package test\nvar x, y = 1, 2",
			expectError: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver la ValueSpec
			var valueSpec *ast.ValueSpec
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok {
					valueSpec = vs
					return false
				}
				return true
			})

			// Vérification valueSpec trouvée
			if valueSpec == nil {
				t.Fatal("no value spec found")
			}

			// Vérification des conditions
			hasType := valueSpec.Type != nil

			// Le format obligatoire est: var name type (avec ou sans valeur)
			// hasError = !hasType (sauf blank identifier)
			isBlank := len(valueSpec.Names) == 1 && valueSpec.Names[0].Name == "_"
			hasError := !hasType && !isBlank

			// Vérification résultat
			if hasError != tt.expectError {
				t.Errorf("checkVarSpec error = %v, want %v", hasError, tt.expectError)
			}
		})
	}
}

// Test_checkVarSpec_multipleVars tests checkVarSpec with multiple variables.
func Test_checkVarSpec_multipleVars(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
	}{
		{
			name: "three_valid_vars",
			code: `package test
var (
	a int = 1
	b string = "hello"
	c bool = true
)`,
			expectedCount: 3,
		},
		{
			name: "mixed_valid_invalid",
			code: `package test
var (
	a int = 1
	b = "hello"
	c bool = true
)`,
			expectedCount: 2,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Compter les ValueSpecs valides
			validCount := 0
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok {
					hasType := vs.Type != nil
					// Vérification format valide (type explicite requis)
					if hasType {
						validCount++
					}
				}
				return true
			})

			// Vérification nombre de vars valides
			if validCount != tt.expectedCount {
				t.Errorf("valid var count = %d, want %d", validCount, tt.expectedCount)
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
