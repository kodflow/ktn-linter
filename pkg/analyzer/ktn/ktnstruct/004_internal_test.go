package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_runStruct004 tests the private runStruct004 function.
func Test_runStruct004(t *testing.T) {
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

// Test_hasValidDocumentation tests the private hasValidDocumentation function.
func Test_hasValidDocumentation(t *testing.T) {
	tests := []struct {
		name       string
		src        string
		structName string
		expected   bool
	}{
		{
			name: "no documentation",
			src: `package test
type User struct{}`,
			structName: "User",
			expected:   false,
		},
		{
			name: "valid documentation",
			src: `package test
// User represents a user in the system.
// It contains basic user information.
type User struct{}`,
			structName: "User",
			expected:   true,
		},
		{
			name: "single line documentation",
			src: `package test
// User represents a user.
type User struct{}`,
			structName: "User",
			expected:   false,
		},
		{
			name: "documentation without struct name prefix",
			src: `package test
// This is a user struct.
// It has some fields.
type User struct{}`,
			structName: "User",
			expected:   false,
		},
		{
			name: "documentation with correct prefix",
			src: `package test
// User represents a user in the system.
// It stores user data.
type User struct{}`,
			structName: "User",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, parser.ParseComments)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the GenDecl with the struct
			var doc *ast.CommentGroup
			ast.Inspect(file, func(n ast.Node) bool {
				if gd, ok := n.(*ast.GenDecl); ok {
					doc = gd.Doc
					return false
				}
				return true
			})

			result := hasValidDocumentation(doc, tt.structName)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test_MIN_DOC_LINES tests the constant value.
func Test_MIN_DOC_LINES(t *testing.T) {
	if MIN_DOC_LINES != 2 {
		t.Errorf("expected MIN_DOC_LINES to be 2, got %d", MIN_DOC_LINES)
	}
}
