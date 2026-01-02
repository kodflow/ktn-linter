package ktnvar

import (
	"go/ast"
	"testing"
)

// TestIsEmptyInterface tests the isEmptyInterface function.
func TestIsEmptyInterface(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		interfaceType *ast.InterfaceType
		expected      bool
	}{
		{
			name:          "nil methods",
			interfaceType: &ast.InterfaceType{Methods: nil},
			expected:      true,
		},
		{
			name: "empty methods list",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{List: nil},
			},
			expected: true,
		},
		{
			name: "empty methods slice",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{List: []*ast.Field{}},
			},
			expected: true,
		},
		{
			name: "non-empty interface",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Method"}},
						},
					},
				},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isEmptyInterface(tt.interfaceType)
			// Check result
			if result != tt.expected {
				t.Errorf("isEmptyInterface() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckEmptyInterface tests the checkEmptyInterface function.
func TestCheckEmptyInterface(t *testing.T) {
	// Test: non-empty interface
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: "Method"}},
				},
			},
		},
	}
	// Should not panic
	checkEmptyInterface(nil, interfaceType)
}
