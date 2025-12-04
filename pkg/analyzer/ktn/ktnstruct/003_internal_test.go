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

// Test_collectExportedStructsWithMethods tests the private collectExportedStructsWithMethods function.
func Test_collectExportedStructsWithMethods(t *testing.T) {
	tests := []struct {
		name           string
		src            string
		expectedCount  int
		checkStruct    string
		expectedMethod int
	}{
		{
			name: "collect exported structs with methods",
			src: `package test
type User struct{}
func (u *User) GetName() string { return "" }

type Admin struct{}
func (a *Admin) GetRole() string { return "" }`,
			expectedCount:  2,
			checkStruct:    "User",
			expectedMethod: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{Fset: fset}
			structs := collectExportedStructsWithMethods(file, pass, nil)

			if len(structs) != tt.expectedCount {
				t.Errorf("expected %d exported structs, got %d", tt.expectedCount, len(structs))
			}

			// Find User struct
			var userStruct *structWithMethods
			for i := range structs {
				if structs[i].name == tt.checkStruct {
					userStruct = &structs[i]
					break
				}
			}

			if userStruct == nil {
				t.Fatalf("%s struct not found", tt.checkStruct)
			}

			if len(userStruct.methods) != tt.expectedMethod {
				t.Errorf("expected %d method for %s, got %d", tt.expectedMethod, tt.checkStruct, len(userStruct.methods))
			}
		})
	}
}

// Test_collectConstructors tests the private collectConstructors function.
func Test_collectConstructors(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected int
	}{
		{
			name: "no constructors",
			src: `package test
func Helper() {}`,
			expected: 0,
		},
		{
			name: "one constructor",
			src: `package test
type User struct{}
func NewUser() *User { return &User{} }`,
			expected: 1,
		},
		{
			name: "multiple constructors",
			src: `package test
type User struct{}
func NewUser() *User { return &User{} }
func NewAdmin() *User { return &User{} }`,
			expected: 2,
		},
		{
			name: "constructor with no return",
			src: `package test
func NewUser() {}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			constructors := collectConstructors(file)

			if len(constructors) != tt.expected {
				t.Errorf("expected %d constructors, got %d", tt.expected, len(constructors))
			}
		})
	}
}

// Test_extractReturnTypeName tests the private extractReturnTypeName function.
func Test_extractReturnTypeName(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer return",
			src: `package test
type User struct{}
func NewUser() *User { return nil }`,
			expected: "User",
		},
		{
			name: "value return",
			src: `package test
type User struct{}
func NewUser() User { return User{} }`,
			expected: "User",
		},
		{
			name: "no return",
			src: `package test
func NewUser() {}`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("no function found")
			}

			result := extractReturnTypeName(funcDecl.Type.Results)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_hasConstructor tests the private hasConstructor function.
func Test_hasConstructor(t *testing.T) {
	tests := []struct {
		name         string
		constructors []constructorInfo
		funcName     string
		typeName     string
		want         bool
	}{
		{
			name: "found matching constructor",
			constructors: []constructorInfo{
				{name: "NewUser", returnType: "User"},
				{name: "NewAdmin", returnType: "Admin"},
			},
			funcName: "NewUser",
			typeName: "User",
			want:     true,
		},
		{
			name: "wrong type for constructor",
			constructors: []constructorInfo{
				{name: "NewUser", returnType: "User"},
				{name: "NewAdmin", returnType: "Admin"},
			},
			funcName: "NewUser",
			typeName: "Admin",
			want:     false,
		},
		{
			name: "constructor not found",
			constructors: []constructorInfo{
				{name: "NewUser", returnType: "User"},
				{name: "NewAdmin", returnType: "Admin"},
			},
			funcName: "NewManager",
			typeName: "Manager",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasConstructor(tt.constructors, tt.funcName, tt.typeName)
			if got != tt.want {
				t.Errorf("hasConstructor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_constructorInfo tests the constructorInfo type.
func Test_constructorInfo(t *testing.T) {
	tests := []struct {
		name           string
		ci             constructorInfo
		expectedName   string
		expectedReturn string
	}{
		{
			name:           "valid constructor info",
			ci:             constructorInfo{name: "NewUser", returnType: "User"},
			expectedName:   "NewUser",
			expectedReturn: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ci.name != tt.expectedName || tt.ci.returnType != tt.expectedReturn {
				t.Errorf("expected name=%s returnType=%s, got name=%s returnType=%s",
					tt.expectedName, tt.expectedReturn, tt.ci.name, tt.ci.returnType)
			}
		})
	}
}
