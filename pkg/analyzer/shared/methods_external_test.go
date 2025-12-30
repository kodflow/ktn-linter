// External tests for methods.go functions.
package shared_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
)

// TestCollectMethodsByStruct tests the CollectMethodsByStruct function.
func TestCollectMethodsByStruct(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
		structName    string
		methodName    string
	}{
		{
			name: "struct with one method",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() string {
	return u.name
}`,
			expectedCount: 1,
			structName:    "User",
			methodName:    "Name",
		},
		{
			name: "struct with multiple methods",
			code: `package test
type User struct {
	name string
	age int
}
func (u *User) Name() string {
	return u.name
}
func (u *User) Age() int {
	return u.age
}`,
			expectedCount: 2,
			structName:    "User",
			methodName:    "Name",
		},
		{
			name: "function without receiver",
			code: `package test
func DoSomething() {
}`,
			expectedCount: 0,
		},
		{
			name: "private method not collected",
			code: `package test
type User struct {
	name string
}
func (u *User) getName() string {
	return u.name
}`,
			expectedCount: 0,
		},
		{
			name: "value receiver",
			code: `package test
type User struct {
	name string
}
func (u User) Name() string {
	return u.name
}`,
			expectedCount: 1,
			structName:    "User",
			methodName:    "Name",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			methods := shared.CollectMethodsByStruct(f, nil)

			// Count total methods
			totalMethods := 0
			// Iterate over methods map
			for _, methodList := range methods {
				totalMethods += len(methodList)
			}

			// Verify count
			if totalMethods != tt.expectedCount {
				t.Errorf("CollectMethodsByStruct() returned %d methods, want %d", totalMethods, tt.expectedCount)
			}

			// Verify specific method if expected
			if tt.structName != "" && tt.methodName != "" {
				methodList, exists := methods[tt.structName]
				// Check struct exists
				if !exists {
					t.Errorf("Expected struct %q not found", tt.structName)
					return
				}
				// Check method exists
				found := false
				// Iterate over methods
				for _, method := range methodList {
					// Check method name
					if method == tt.methodName {
						found = true
						break
					}
				}
				// Verify method found
				if !found {
					t.Errorf("Expected method %q not found for struct %q", tt.methodName, tt.structName)
				}
			}
		})
	}
}

// TestExtractReceiverName tests the ExtractReceiverName function.
func TestExtractReceiverName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "pointer receiver",
			code: `package test
type User struct{}
func (u *User) Method() {}`,
			expected: "User",
		},
		{
			name: "value receiver",
			code: `package test
type User struct{}
func (u User) Method() {}`,
			expected: "User",
		},
		{
			name: "empty receiver list",
			code: `package test
func Function() {}`,
			expected: "",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			// Inspect AST
			ast.Inspect(f, func(n ast.Node) bool {
				// Check function declaration
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				// Continue traversal
				return true
			})

			// Check function found
			if funcDecl == nil {
				t.Fatal("No function declaration found")
			}

			// Test ExtractReceiverName
			result := shared.ExtractReceiverName(funcDecl.Recv)
			// Verify result
			if result != tt.expected {
				t.Errorf("ExtractReceiverName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestExtractReceiverNameEdgeCases tests ExtractReceiverName with edge cases.
func TestExtractReceiverNameEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		recv     *ast.FieldList
		expected string
	}{
		{
			name:     "nil_receiver",
			recv:     nil,
			expected: "",
		},
		{
			name:     "empty_list",
			recv:     &ast.FieldList{List: []*ast.Field{}},
			expected: "",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := shared.ExtractReceiverName(tt.recv)
			// Verify result
			if result != tt.expected {
				t.Errorf("ExtractReceiverName() = %q, want %q", result, tt.expected)
			}
		})
	}
}
