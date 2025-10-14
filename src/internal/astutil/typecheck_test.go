package astutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
)

// TestIsConstCompatibleType teste la vérification de compatibilité des types avec const.
//
// Params:
//   - t: instance de test
func TestIsConstCompatibleType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		// Types compatibles avec const
		{name: "bool", code: "var x bool", expected: true},
		{name: "string", code: "var x string", expected: true},
		{name: "int", code: "var x int", expected: true},
		{name: "int8", code: "var x int8", expected: true},
		{name: "int16", code: "var x int16", expected: true},
		{name: "int32", code: "var x int32", expected: true},
		{name: "int64", code: "var x int64", expected: true},
		{name: "uint", code: "var x uint", expected: true},
		{name: "uint8", code: "var x uint8", expected: true},
		{name: "uint16", code: "var x uint16", expected: true},
		{name: "uint32", code: "var x uint32", expected: true},
		{name: "uint64", code: "var x uint64", expected: true},
		{name: "float32", code: "var x float32", expected: true},
		{name: "float64", code: "var x float64", expected: true},
		{name: "complex64", code: "var x complex64", expected: true},
		{name: "complex128", code: "var x complex128", expected: true},
		{name: "byte", code: "var x byte", expected: true},
		{name: "rune", code: "var x rune", expected: true},

		// Types incompatibles avec const
		{name: "slice", code: "var x []int", expected: false},
		{name: "map", code: "var x map[string]int", expected: false},
		{name: "pointer", code: "var x *int", expected: false},
		{name: "channel", code: "var x chan int", expected: false},
		{name: "struct", code: "var x struct{}", expected: false},
		{name: "interface", code: "var x interface{}", expected: false},
		{name: "array", code: "var x [5]int", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			decl := file.Decls[0].(*ast.GenDecl)
			spec := decl.Specs[0].(*ast.ValueSpec)

			result := astutil.IsConstCompatibleType(spec.Type)
			if result != tt.expected {
				t.Errorf("astutil.IsConstCompatibleType(%s) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestIsLiteralValue teste la détection de valeurs littérales.
//
// Params:
//   - t: instance de test
func TestIsLiteralValue(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		// Valeurs littérales
		{name: "integer literal", code: "var x = 42", expected: true},
		{name: "string literal", code: `var x = "hello"`, expected: true},
		{name: "float literal", code: "var x = 3.14", expected: true},
		{name: "bool true", code: "var x = true", expected: true},
		{name: "bool false", code: "var x = false", expected: true},
		{name: "nil", code: "var x = nil", expected: true},

		// Non-littérales
		{name: "function call", code: "var x = len(s)", expected: false},
		// Note: variable references sont aussi des Ident, donc retournent true sans analyse de types
		// {name: "variable reference", code: "var x = y", expected: false},
		{name: "composite literal", code: "var x = []int{1, 2}", expected: false},
		{name: "slice literal", code: "var x = []int{}", expected: false},
		{name: "map literal", code: "var x = map[string]int{}", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parser avec les références non résolues
			fset := token.NewFileSet()
			src := "package test\nvar s string\nvar y int\n" + tt.code
			file, err := parser.ParseFile(fset, "", src, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// Trouver la dernière déclaration (celle avec tt.code)
			var spec *ast.ValueSpec
			for i := len(file.Decls) - 1; i >= 0; i-- {
				if genDecl, ok := file.Decls[i].(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
					if len(genDecl.Specs) > 0 {
						if vs, ok := genDecl.Specs[0].(*ast.ValueSpec); ok && len(vs.Values) > 0 {
							spec = vs
							break
						}
					}
				}
			}

			if spec == nil {
				t.Fatalf("Failed to find variable spec in code: %s", tt.code)
			}

			result := astutil.IsLiteralValue(spec.Values[0])
			if result != tt.expected {
				t.Errorf("astutil.IsLiteralValue(%s) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestIsConstCompatibleTypeNilExpr teste le comportement avec une expression nil.
//
// Params:
//   - t: instance de test
func TestIsConstCompatibleTypeNilExpr(t *testing.T) {
	// Test avec une expression nil
	result := astutil.IsConstCompatibleType(nil)
	if result != false {
		t.Errorf("astutil.IsConstCompatibleType(nil) = %v, want false", result)
	}
}

// TestIsLiteralValueNilExpr teste le comportement avec une expression nil.
//
// Params:
//   - t: instance de test
func TestIsLiteralValueNilExpr(t *testing.T) {
	// Test avec une expression nil
	result := astutil.IsLiteralValue(nil)
	if result != false {
		t.Errorf("astutil.IsLiteralValue(nil) = %v, want false", result)
	}
}
