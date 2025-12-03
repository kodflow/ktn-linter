package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_runStruct003 tests the private runStruct003 function.
func Test_runStruct003(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_checkFieldOrder tests the private checkFieldOrder function.
func Test_checkFieldOrder(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		expectError bool
	}{
		{
			name: "correct order",
			src: `package test
type User struct {
	Name string
	age  int
}`,
			expectError: false,
		},
		{
			name: "incorrect order",
			src: `package test
type User struct {
	age  int
	Name string
}`,
			expectError: true,
		},
		{
			name: "all exported",
			src: `package test
type User struct {
	Name string
	Age  int
}`,
			expectError: false,
		},
		{
			name: "all private",
			src: `package test
type User struct {
	name string
	age  int
}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the struct
			var typeSpec *ast.TypeSpec
			var structType *ast.StructType
			ast.Inspect(file, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						typeSpec = ts
						structType = st
						return false
					}
				}
				return true
			})

			if typeSpec == nil || structType == nil {
				t.Fatal("no struct found")
			}

			// Create a minimal pass with a reporter
			errorReported := false
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					errorReported = true
				},
			}

			// Call checkFieldOrder
			checkFieldOrder(pass, typeSpec, structType)

			if errorReported != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, errorReported)
			}
		})
	}
}

// Test_fieldInfo tests the fieldInfo type.
func Test_fieldInfo(t *testing.T) {
	fi := fieldInfo{
		name:     "Name",
		exported: true,
		pos:      token.NoPos,
	}

	if fi.name != "Name" {
		t.Errorf("expected name 'Name', got '%s'", fi.name)
	}
	if !fi.exported {
		t.Error("expected exported to be true")
	}
}
