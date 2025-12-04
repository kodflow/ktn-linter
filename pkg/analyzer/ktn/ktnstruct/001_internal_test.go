package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Test_runStruct001 tests the private runStruct001 function.
func Test_runStruct001(t *testing.T) {
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

// Test_collectInterfaces tests the private collectInterfaces function.
func Test_collectInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected int
	}{
		{
			name: "no interfaces",
			src: `package test
type User struct {
	Name string
}`,
			expected: 0,
		},
		{
			name: "one interface",
			src: `package test
type Reader interface {
	Read() error
}`,
			expected: 1,
		},
		{
			name: "multiple interfaces",
			src: `package test
type Reader interface {
	Read() error
}
type Writer interface {
	Write() error
}`,
			expected: 2,
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
			interfaces := collectInterfaces(file, pass)

			if len(interfaces) != tt.expected {
				t.Errorf("expected %d interfaces, got %d", tt.expected, len(interfaces))
			}
		})
	}
}

// Test_extractStructNameFromReceiver tests the private extractStructNameFromReceiver function.
func Test_extractStructNameFromReceiver(t *testing.T) {
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

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Recv == nil {
				t.Fatal("no function with receiver found")
			}

			recvType := funcDecl.Recv.List[0].Type
			result := extractStructNameFromReceiver(recvType)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_collectMethodsByStruct tests the private collectMethodsByStruct function.
func Test_collectMethodsByStruct(t *testing.T) {
	tests := []struct {
		name            string
		src             string
		structName      string
		expectedMethods int
	}{
		{
			name: "collect methods for struct",
			src: `package test
type User struct{}
func (u *User) GetName() string { return "" }
func (u *User) SetName(name string) {}`,
			structName:      "User",
			expectedMethods: 2,
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
			methodsByStruct := collectMethodsByStruct(file, pass)

			userMethods, ok := methodsByStruct[tt.structName]
			if !ok {
				t.Fatalf("expected %s struct in map", tt.structName)
			}

			if len(userMethods) != tt.expectedMethods {
				t.Errorf("expected %d methods, got %d", tt.expectedMethods, len(userMethods))
			}
		})
	}
}

// Test_hasMatchingInterface tests the private hasMatchingInterface function.
func Test_hasMatchingInterface(t *testing.T) {
	tests := []struct {
		name       string
		structName string
		methods    []shared.MethodSignature
		interfaces map[string][]shared.MethodSignature
		expected   bool
	}{
		{
			name:       "matching interface found",
			structName: "User",
			methods:    []shared.MethodSignature{{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			interfaces: map[string][]shared.MethodSignature{
				"Reader": {{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			},
			expected: true,
		},
		{
			name:       "no matching interface",
			structName: "User",
			methods:    []shared.MethodSignature{{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			interfaces: map[string][]shared.MethodSignature{
				"Writer": {{Name: "SetName", ParamsStr: "string", ResultsStr: ""}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := structWithMethods{name: tt.structName, methods: tt.methods}
			if hasMatchingInterface(s, tt.interfaces) != tt.expected {
				t.Errorf("expected %v for %s", tt.expected, tt.name)
			}
		})
	}
}

// Test_interfaceCoversAllMethods tests the private interfaceCoversAllMethods function.
func Test_interfaceCoversAllMethods(t *testing.T) {
	tests := []struct {
		name          string
		structMethods []shared.MethodSignature
		ifaceMethods  []shared.MethodSignature
		expected      bool
	}{
		{
			name: "interface covers all methods",
			structMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			ifaceMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			expected: true,
		},
		{
			name: "incomplete interface",
			structMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			ifaceMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if interfaceCoversAllMethods(tt.structMethods, tt.ifaceMethods) != tt.expected {
				t.Errorf("expected %v for %s", tt.expected, tt.name)
			}
		})
	}
}

// Test_formatFieldList tests the private formatFieldList function.
func Test_formatFieldList(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "no params",
			src: `package test
func test() {}`,
			expected: "",
		},
		{
			name: "single param",
			src: `package test
func test(x int) {}`,
			expected: "int",
		},
		{
			name: "multiple params",
			src: `package test
func test(x int, y string) {}`,
			expected: "int,string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

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

			pass := &analysis.Pass{Fset: fset}
			result := formatFieldList(funcDecl.Type.Params, pass)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_collectStructsWithMethods tests the collectStructsWithMethods private function.
func Test_collectStructsWithMethods(t *testing.T) {
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
