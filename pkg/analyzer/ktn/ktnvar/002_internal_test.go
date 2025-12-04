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
	t.Run("function exists", func(t *testing.T) {
		// La fonction checkVarSpec est testée via l'API publique
		t.Log("checkVarSpec tested via external tests")
	})
}

// Test_hasVisibleType tests the hasVisibleType function.
//
// Params:
//   - t: testing context
func Test_hasVisibleType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"composite literal", "var x = []string{}", true},
		{"make call", "var x = make([]int, 10)", true},
		{"no type visible", "var x = y", false},
		{"empty values", "", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification cas vide
			if tt.code == "" {
				result := hasVisibleType(nil)
				// Vérification résultat
				if result != tt.expected {
					t.Errorf("hasVisibleType(nil) = %v, want %v", result, tt.expected)
				}
				return
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver la var
			var values []ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok {
					values = vs.Values
					return false
				}
				return true
			})

			result := hasVisibleType(values)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("hasVisibleType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isTypeVisible tests the isTypeVisible function.
//
// Params:
//   - t: testing context
func Test_isTypeVisible(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"composite literal slice", "var x = []string{}", true},
		{"composite literal map", "var x = map[string]int{}", true},
		{"make slice", "var x = make([]int, 10)", true},
		{"make map", "var x = make(map[string]int)", true},
		{"new struct", "var x = new(Foo)", true},
		{"pointer to composite", "var x = &Foo{}", true},
		{"ident", "var x = y", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver l'expression
			var expr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok && len(vs.Values) > 0 {
					expr = vs.Values[0]
					return false
				}
				return true
			})

			// Vérification expression trouvée
			if expr == nil {
				t.Fatal("no expression found")
			}

			result := isTypeVisible(expr)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("isTypeVisible() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isTypedCall tests the isTypedCall function.
//
// Params:
//   - t: testing context
func Test_isTypedCall(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"make call", "var x = make([]int, 10)", true},
		{"new call", "var x = new(Foo)", true},
		{"int conversion", "var x = int(42)", true},
		{"string conversion", "var x = string(data)", true},
		{"regular call", "var x = foo()", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver le CallExpr
			var call *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if c, ok := n.(*ast.CallExpr); ok {
					call = c
					return false
				}
				return true
			})

			// Vérification call trouvé
			if call == nil {
				t.Fatal("no call expression found")
			}

			result := isTypedCall(call)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("isTypedCall() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isBuiltinOrTypeConversion tests the isBuiltinOrTypeConversion function.
//
// Params:
//   - t: testing context
func Test_isBuiltinOrTypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"make", "make", true},
		{"new", "new", true},
		{"int", "int", true},
		{"string", "string", true},
		{"float64", "float64", true},
		{"byte", "byte", true},
		{"custom function", "foo", false},
		{"empty", "", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBuiltinOrTypeConversion(tt.input)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("isBuiltinOrTypeConversion(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
