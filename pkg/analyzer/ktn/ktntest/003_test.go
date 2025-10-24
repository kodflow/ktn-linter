package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TestIsPublicFunction tests isPublicFunction helper
func TestIsPublicFunction(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"public function", "func MyFunc() {}", true},
		{"private function", "func myFunc() {}", false},
		{"public with receiver", "func (r *Receiver) MyMethod() {}", true},
		{"private with receiver", "func (r *Receiver) myMethod() {}", false},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Get the function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					// Early return from iteration.
					return false
				}
				// Continue traversal
				return true
			})

			if funcDecl == nil {
				t.Fatal("No function declaration found")
			}

			got := isPublicFunction(funcDecl)
			if got != tt.want {
				t.Errorf("isPublicFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCollectTestedFunctions tests collectTestedFunctions helper
func TestCollectTestedFunctions(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		wantKey  string
	}{
		{"test function", "TestMyFunc", "MyFunc"},
		{"test with underscore", "TestMy_Func", "My_Func"},
		{"non-test function", "MyFunc", ""},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\nfunc " + tt.funcName + "(t *testing.T) {}"
			file, err := parser.ParseFile(fset, "", code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			testedFuncs := make(map[string]bool, 0)
			// Get the function declaration
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					collectTestedFunctions(fd, testedFuncs)
				}
				// Continue traversal
				return true
			})

			if tt.wantKey != "" {
				if !testedFuncs[tt.wantKey] {
					t.Errorf("Expected %q in testedFuncs, but not found", tt.wantKey)
				}
			} else {
				if len(testedFuncs) > 0 {
					t.Errorf("Expected empty testedFuncs, got %v", testedFuncs)
				}
			}
		})
	}
}

// TestIsExemptFunction tests isExemptFunction helper
func TestIsExemptFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		{"init function", "init", true},
		{"main function", "main", true},
		{"regular function", "MyFunc", false},
		{"Test function", "TestSomething", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptFunction(tt.funcName)
			if got != tt.want {
				t.Errorf("isExemptFunction(%q) = %v, want %v", tt.funcName, got, tt.want)
			}
		})
	}
}

// TestCollectFunctions tests collectFunctions helper
func TestCollectFunctions(t *testing.T) {
	code := `
package mypackage

// Public function
func PublicFunc() {}

// private function
func privateFunc() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create a minimal pass
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
	}

	var publicFuncs []publicFuncInfo
	testedFuncs := make(map[string]bool, 0)

	collectFunctions(pass, &publicFuncs, testedFuncs)

	// Should find only PublicFunc
	if len(publicFuncs) != 1 {
		t.Errorf("Expected 1 public function, got %d", len(publicFuncs))
	}

	if len(publicFuncs) > 0 && publicFuncs[0].name != "PublicFunc" {
		t.Errorf("Expected PublicFunc, got %s", publicFuncs[0].name)
	}
}
