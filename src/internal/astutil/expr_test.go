package astutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
)

func TestExprToString(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "simple identifier",
			code:     "var x int",
			expected: "int",
		},
		{
			name:     "qualified identifier",
			code:     "var x time.Time",
			expected: "time.Time",
		},
		{
			name:     "slice type",
			code:     "var x []string",
			expected: "[]string",
		},
		{
			name:     "nested slice",
			code:     "var x [][]int",
			expected: "[][]int",
		},
		{
			name:     "map type",
			code:     "var x map[string]int",
			expected: "map[string]int",
		},
		{
			name:     "complex map",
			code:     "var x map[string][]int",
			expected: "map[string][]int",
		},
		{
			name:     "pointer type",
			code:     "var x *string",
			expected: "*string",
		},
		{
			name:     "double pointer",
			code:     "var x **int",
			expected: "**int",
		},
		{
			name:     "channel type",
			code:     "var x chan int",
			expected: "chan int",
		},
		{
			name:     "send channel",
			code:     "var x chan<- string",
			expected: "chan<- string",
		},
		{
			name:     "receive channel",
			code:     "var x <-chan bool",
			expected: "<-chan bool",
		},
		{
			name:     "complex nested type",
			code:     "var x *[]map[string]*int",
			expected: "*[]map[string]*int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// Extraire la déclaration de variable
			decl := file.Decls[0].(*ast.GenDecl)
			spec := decl.Specs[0].(*ast.ValueSpec)

			result := astutil.ExprToString(spec.Type)
			if result != tt.expected {
				t.Errorf("ExprToString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExprToString_Unknown(t *testing.T) {
	// Test avec un type non supporté (struct type)
	code := "var x struct { name string }"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", "package test\n"+code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	decl := file.Decls[0].(*ast.GenDecl)
	spec := decl.Specs[0].(*ast.ValueSpec)

	result := astutil.ExprToString(spec.Type)
	if result != "unknown" {
		t.Errorf("ExprToString() for unsupported type = %q, want %q", result, "unknown")
	}
}

func TestGetTypeString(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "with explicit type",
			code:     "var x int = 42",
			expected: "int",
		},
		{
			name:     "with pointer type",
			code:     "var x *string = nil",
			expected: "*string",
		},
		{
			name:     "with slice type",
			code:     "var x []int = nil",
			expected: "[]int",
		},
		{
			name:     "without type (inferred)",
			code:     "var x = 42",
			expected: "<type>",
		},
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

			result := astutil.GetTypeString(spec)
			if result != tt.expected {
				t.Errorf("GetTypeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetTypeString_Const(t *testing.T) {
	// Test avec une constante
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "const with explicit type",
			code:     "const x int = 42",
			expected: "int",
		},
		{
			name:     "const without type",
			code:     "const x = 42",
			expected: "<type>",
		},
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

			result := astutil.GetTypeString(spec)
			if result != tt.expected {
				t.Errorf("GetTypeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}
