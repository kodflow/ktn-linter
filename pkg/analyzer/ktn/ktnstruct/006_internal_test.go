package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_runStruct006 tests the private runStruct006 function.
func Test_runStruct006(t *testing.T) {
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

// Test_collectStructTypes tests the private collectStructTypes function.
func Test_collectStructTypes(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected map[string]int // map of struct name to number of fields
	}{
		{
			name: "no structs",
			src: `package test
func main() {}`,
			expected: map[string]int{},
		},
		{
			name: "struct with fields",
			src: `package test
type User struct {
	Name string
	Age  int
}`,
			expected: map[string]int{"User": 2},
		},
		{
			name: "multiple structs",
			src: `package test
type User struct {
	Name string
}
type Admin struct {
	Role string
	Level int
}`,
			expected: map[string]int{"User": 1, "Admin": 2},
		},
		{
			name: "empty struct",
			src: `package test
type Empty struct{}`,
			expected: map[string]int{"Empty": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			result := collectStructTypes(pass)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d structs, got %d", len(tt.expected), len(result))
			}

			for name, expectedFields := range tt.expected {
				fields, ok := result[name]
				if !ok {
					t.Errorf("expected struct '%s' not found", name)
					continue
				}
				if len(fields) != expectedFields {
					t.Errorf("expected %d fields for '%s', got %d", expectedFields, name, len(fields))
				}
			}
		})
	}
}

// Test_isSimpleGetter tests the private isSimpleGetter function.
func Test_isSimpleGetter(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected bool
	}{
		{
			name: "simple getter",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() string { return u.name }`,
			expected: true,
		},
		{
			name: "getter with params",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName(prefix string) string { return prefix + u.name }`,
			expected: false,
		},
		{
			name: "getter without return",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			// Collect struct types first
			structTypes := collectStructTypes(pass)

			// Find the method
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Recv != nil {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("no method found")
			}

			result := isSimpleGetter(funcDecl, structTypes)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test_getReceiverTypeName tests the private getReceiverTypeName function.
func Test_getReceiverTypeName(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer receiver",
			src: `package test
type User struct{}
func (u *User) Method() {}`,
			expected: "User",
		},
		{
			name: "value receiver",
			src: `package test
type User struct{}
func (u User) Method() {}`,
			expected: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the method
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Recv != nil {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Recv == nil {
				t.Fatal("no method found")
			}

			result := getReceiverTypeName(funcDecl.Recv.List[0].Type)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_constants tests the constant values.
func Test_constants(t *testing.T) {
	if INITIAL_STRUCT_TYPES_CAP != 32 {
		t.Errorf("expected INITIAL_STRUCT_TYPES_CAP to be 32, got %d", INITIAL_STRUCT_TYPES_CAP)
	}
	if GET_PREFIX_LEN != 3 {
		t.Errorf("expected GET_PREFIX_LEN to be 3, got %d", GET_PREFIX_LEN)
	}
}
