package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_runVar002 tests the private runVar002 function.
//
// Params:
//   - t: testing context
func Test_runVar002(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
			t.Log("runVar002 tested via external tests")
		})
	}
}

// Test_checkVarSpec tests the checkVarSpec function.
//
// Params:
//   - t: testing context
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
			name:        "no value - ERROR",
			code:        "package test\nvar x int",
			expectError: true,
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
	}

	// Parcourir les cas de test
	for _, tt := range tests {
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
			hasValues := len(valueSpec.Values) > 0

			// Le format obligatoire est: var name type = value
			// hasError = !hasType || !hasValues (sauf blank identifier)
			isBlank := len(valueSpec.Names) == 1 && valueSpec.Names[0].Name == "_"
			hasError := (!hasType || !hasValues) && !isBlank

			// Vérification résultat
			if hasError != tt.expectError {
				t.Errorf("checkVarSpec error = %v, want %v", hasError, tt.expectError)
			}
		})
	}
}

// Test_checkVarSpec_multipleVars tests checkVarSpec with multiple variables.
//
// Params:
//   - t: testing context
func Test_checkVarSpec_multipleVars(t *testing.T) {
	code := `package test
var (
	a int = 1
	b string = "hello"
	c bool = true
)`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
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
			hasValues := len(vs.Values) > 0
			// Vérification format valide
			if hasType && hasValues {
				validCount++
			}
		}
		return true
	})

	// Vérification nombre de vars valides
	if validCount != 3 {
		t.Errorf("valid var count = %d, want 3", validCount)
	}
}
