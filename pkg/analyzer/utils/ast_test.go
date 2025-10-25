package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestGetExprAsString tests the functionality of the corresponding implementation.
func TestGetExprAsString(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"ident", "myVar", "myVar"},
		{"selector", "pkg.Type", "pkg.Type"},
		{"array", "[]int", "[]int"},
		{"map", "map[string]int", "map[string]int"},
		{"pointer", "*int", "*int"},
		{"chan bidirectional", "chan int", "chan int"},
		{"chan send", "chan<- int", "chan<- int"},
		{"chan recv", "<-chan int", "<-chan int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.code, err)
			}
			got := GetExprAsString(expr)
			if got != tt.expected {
				t.Errorf("GetExprAsString(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

// TestGetExprAsStringWithUnknownType tests the functionality of the corresponding implementation.
func TestGetExprAsStringWithUnknownType(t *testing.T) {
	// Test avec un type non supporté (FuncType)
	expr := &ast.FuncType{
		Params: &ast.FieldList{},
	}
	got := GetExprAsString(expr)
	if got != "unknown" {
		t.Errorf("GetExprAsString(unsupported) = %q, want \"unknown\"", got)
	}
}

// TestGetTypeString tests the functionality of the corresponding implementation.
func TestGetTypeString(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"with explicit type", "var x int", "int"},
		{"with slice type", "var s []string", "[]string"},
		{"with map type", "var m map[string]int", "map[string]int"},
		{"with pointer type", "var p *int", "*int"},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.code, err)
			}

			// Extract the ValueSpec from the GenDecl
			genDecl := file.Decls[0].(*ast.GenDecl)
			spec := genDecl.Specs[0].(*ast.ValueSpec)

			got := GetTypeString(spec)
			if got != tt.expected {
				t.Errorf("GetTypeString(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

// TestGetTypeStringWithNoType tests the functionality of the corresponding implementation.
func TestGetTypeStringWithNoType(t *testing.T) {
	// Test avec une ValueSpec sans type explicite
	fset := token.NewFileSet()
	code := "package test\nvar x = 5"
	file, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	genDecl := file.Decls[0].(*ast.GenDecl)
	spec := genDecl.Specs[0].(*ast.ValueSpec)

	got := GetTypeString(spec)
	if got != "<type>" {
		t.Errorf("GetTypeString(no type) = %q, want \"<type>\"", got)
	}
}

// TestGetExprAsStringNested tests the functionality of the corresponding implementation.
func TestGetExprAsStringNested(t *testing.T) {
	// Test avec des types imbriqués complexes
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"nested slice", "[][]int", "[][]int"},
		{"slice of pointers", "[]*int", "[]*int"},
		{"map of slices", "map[string][]int", "map[string][]int"},
		{"chan of pointers", "chan *int", "chan *int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.code, err)
			}
			got := GetExprAsString(expr)
			if got != tt.expected {
				t.Errorf("GetExprAsString(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}
