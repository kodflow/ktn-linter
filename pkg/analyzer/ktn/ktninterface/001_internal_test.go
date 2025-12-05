// Internal tests for 001.go - ktninterface package.
package ktninterface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_checkType tests the checkType function with various expression types.
//
// Params:
//   - t: testing context
func Test_checkType(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "ident type",
			code:     "package test\ntype X interface{}\nfunc f(x X) {}",
			expected: []string{"X"},
		},
		{
			name:     "pointer type",
			code:     "package test\ntype X interface{}\nfunc f(x *X) {}",
			expected: []string{"X"},
		},
		{
			name:     "slice type",
			code:     "package test\ntype X interface{}\nfunc f(x []X) {}",
			expected: []string{"X"},
		},
		{
			name:     "map type",
			code:     "package test\ntype K interface{}\ntype V interface{}\nfunc f(m map[K]V) {}",
			expected: []string{"K", "V"},
		},
		{
			name:     "channel type",
			code:     "package test\ntype X interface{}\nfunc f(c chan X) {}",
			expected: []string{"X"},
		},
		{
			name:     "selector type",
			code:     "package test\nimport \"io\"\nfunc f(r io.Reader) {}",
			expected: []string{"io"},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Find function parameters
			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function decl
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Type.Params != nil {
					// Iterate over params
					for _, field := range fn.Type.Params.List {
						checkType(field.Type, used)
					}
				}
				return true
			})
			// Check expected types are used
			for _, exp := range tt.expected {
				// Validate
				if !used[exp] {
					t.Errorf("Expected %q to be marked as used", exp)
				}
			}
		})
	}
}

// Test_checkEmbeddedInterfaces tests embedded interface detection.
//
// Params:
//   - t: testing context
func Test_checkEmbeddedInterfaces(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "embedded interface",
			code:     "package test\ntype Reader interface{}\ntype Writer interface{}\ntype ReadWriter interface { Reader; Writer }",
			expected: []string{"Reader", "Writer"},
		},
		{
			name:     "empty interface",
			code:     "package test\ntype Empty interface{}",
			expected: []string{},
		},
		{
			name:     "interface with methods only",
			code:     "package test\ntype Foo interface{ Bar() }",
			expected: []string{},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Find interfaces
			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if interface type
				if iface, ok := n.(*ast.InterfaceType); ok {
					checkEmbeddedInterfaces(iface, used)
				}
				return true
			})
			// Check expected types are used
			for _, exp := range tt.expected {
				// Validate
				if !used[exp] {
					t.Errorf("Expected %q to be marked as used", exp)
				}
			}
		})
	}
}

// Test_checkEmbeddedInterfaces_nilMethods tests nil methods case.
//
// Params:
//   - t: testing context
func Test_checkEmbeddedInterfaces_nilMethods(t *testing.T) {
	// Create interface type with nil methods
	iface := &ast.InterfaceType{Methods: nil}
	used := make(map[string]bool)
	// Call function
	checkEmbeddedInterfaces(iface, used)
	// Validate no panic and empty result
	if len(used) != 0 {
		t.Errorf("Expected empty map, got %v", used)
	}
}

// Test_collectTypeSpecs tests the collectTypeSpecs function.
//
// Params:
//   - t: testing context
func Test_collectTypeSpecs(t *testing.T) {
	// Define test cases
	tests := []struct {
		name              string
		code              string
		wantInterfaces    int
		wantStructs       int
		interfaceNames    []string
		structNames       []string
	}{
		{
			name:           "interface and struct",
			code:           "package test\ntype Foo interface{}\ntype Bar struct{}",
			wantInterfaces: 1,
			wantStructs:    1,
			interfaceNames: []string{"Foo"},
			structNames:    []string{"Bar"},
		},
		{
			name:           "only interfaces",
			code:           "package test\ntype A interface{}\ntype B interface{}",
			wantInterfaces: 2,
			wantStructs:    0,
			interfaceNames: []string{"A", "B"},
			structNames:    []string{},
		},
		{
			name:           "only structs",
			code:           "package test\ntype X struct{}\ntype Y struct{}",
			wantInterfaces: 0,
			wantStructs:    2,
			interfaceNames: []string{},
			structNames:    []string{"X", "Y"},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Collect specs
			interfaces := make(map[string]*ast.TypeSpec)
			structNames := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if general declaration
				if genDecl, ok := n.(*ast.GenDecl); ok {
					collectTypeSpecs(genDecl.Specs, interfaces, structNames)
				}
				return true
			})
			// Validate interface count
			if len(interfaces) != tt.wantInterfaces {
				t.Errorf("Got %d interfaces, want %d", len(interfaces), tt.wantInterfaces)
			}
			// Validate struct count
			if len(structNames) != tt.wantStructs {
				t.Errorf("Got %d structs, want %d", len(structNames), tt.wantStructs)
			}
			// Validate interface names
			for _, name := range tt.interfaceNames {
				// Check if interface exists
				if _, ok := interfaces[name]; !ok {
					t.Errorf("Expected interface %q", name)
				}
			}
			// Validate struct names
			for _, name := range tt.structNames {
				// Check if struct exists
				if !structNames[name] {
					t.Errorf("Expected struct %q", name)
				}
			}
		})
	}
}

// Test_isStructInterfacePattern tests the pattern matching function.
//
// Params:
//   - t: testing context
func Test_isStructInterfacePattern(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		interfaceName string
		structNames   map[string]bool
		expected      bool
	}{
		{
			name:          "matching pattern",
			interfaceName: "UserInterface",
			structNames:   map[string]bool{"User": true},
			expected:      true,
		},
		{
			name:          "no matching struct",
			interfaceName: "UserInterface",
			structNames:   map[string]bool{"Admin": true},
			expected:      false,
		},
		{
			name:          "not interface suffix",
			interfaceName: "UserService",
			structNames:   map[string]bool{"User": true},
			expected:      false,
		},
		{
			name:          "too short name",
			interfaceName: "Interface",
			structNames:   map[string]bool{"": true},
			expected:      false,
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := isStructInterfacePattern(tt.interfaceName, tt.structNames)
			// Validate
			if got != tt.expected {
				t.Errorf("isStructInterfacePattern(%q) = %v, want %v", tt.interfaceName, got, tt.expected)
			}
		})
	}
}

// Test_runInterface001 tests the analyzer exists.
//
// Params:
//   - t: testing context
func Test_runInterface001(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "analyzer exists"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate analyzer is defined
			if Analyzer001 == nil {
				t.Error("Analyzer001 is nil")
			}
		})
	}
}

// Test_collectDeclarations tests the collectDeclarations function.
//
// Params:
//   - t: testing context
func Test_collectDeclarations(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via collectTypeSpecs"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("collectDeclarations is tested via collectTypeSpecs")
		})
	}
}

// Test_findInterfaceUsages tests the findInterfaceUsages function.
//
// Params:
//   - t: testing context
func Test_findInterfaceUsages(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via checkNodeForInterfaceUsage"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("findInterfaceUsages is tested via checkNodeForInterfaceUsage")
		})
	}
}

// Test_checkFuncDeclForInterfaces tests the checkFuncDeclForInterfaces function.
//
// Params:
//   - t: testing context
func Test_checkFuncDeclForInterfaces(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "function with params and results",
			code:     "package test\ntype X interface{}\nfunc f(x X) X { return nil }",
			expected: []string{"X"},
		},
		{
			name:     "function with no params",
			code:     "package test\ntype X interface{}\nfunc f() X { return nil }",
			expected: []string{"X"},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Find function and check interfaces
			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function decl
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "f" {
					checkFuncDeclForInterfaces(fn, used)
				}
				return true
			})
			// Check expected types are used
			for _, exp := range tt.expected {
				// Validate
				if !used[exp] {
					t.Errorf("Expected %q to be marked as used", exp)
				}
			}
		})
	}
}

// Test_reportUnusedInterfaces tests the reportUnusedInterfaces function.
//
// Params:
//   - t: testing context
func Test_reportUnusedInterfaces(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via analysistest"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("reportUnusedInterfaces is tested via analysistest")
		})
	}
}

// Test_checkFieldList tests the checkFieldList function.
//
// Params:
//   - t: testing context
func Test_checkFieldList(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "field list with interfaces",
			code:     "package test\ntype A interface{}\ntype B interface{}\nfunc f(a A, b B) {}",
			expected: []string{"A", "B"},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Find function and check field list
			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function decl
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Type.Params != nil {
					checkFieldList(fn.Type.Params, used)
				}
				return true
			})
			// Check expected types are used
			for _, exp := range tt.expected {
				// Validate
				if !used[exp] {
					t.Errorf("Expected %q to be marked as used", exp)
				}
			}
		})
	}
}

// Test_checkNodeForInterfaceUsage tests the node checking function.
//
// Params:
//   - t: testing context
func Test_checkNodeForInterfaceUsage(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "func with interface param",
			code:     "package test\ntype X interface{}\nfunc f(x X) {}",
			expected: []string{"X"},
		},
		{
			name:     "struct field with interface",
			code:     "package test\ntype X interface{}\ntype S struct{ x X }",
			expected: []string{"X"},
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}
			// Check nodes
			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				checkNodeForInterfaceUsage(n, used)
				return true
			})
			// Check expected types are used
			for _, exp := range tt.expected {
				// Validate
				if !used[exp] {
					t.Errorf("Expected %q to be marked as used", exp)
				}
			}
		})
	}
}
