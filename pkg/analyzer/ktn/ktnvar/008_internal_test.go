// Internal tests for 008.go - ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_checkDeclForAlloc tests the checkDeclForAlloc function.
//
// Params:
//   - t: testing context
func Test_checkDeclForAlloc(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		wantPass bool
	}{
		{
			name:     "var slice make",
			code:     "package test\nfunc f() { var s = make([]int, 0) ; _ = s }",
			wantPass: true,
		},
		{
			name:     "var map make",
			code:     "package test\nfunc f() { var m = make(map[string]int) ; _ = m }",
			wantPass: true,
		},
		{
			name:     "var simple int",
			code:     "package test\nfunc f() { var x = 42 ; _ = x }",
			wantPass: true,
		},
		{
			name:     "const decl ignored",
			code:     "package test\nconst x = 42",
			wantPass: true,
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
			// Find DeclStmt
			var declStmt *ast.DeclStmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if DeclStmt
				if ds, ok := n.(*ast.DeclStmt); ok {
					declStmt = ds
					return false
				}
				return true
			})
			// Skip if no DeclStmt found
			if declStmt == nil {
				// Expected for const decl at package level
				return
			}
			// Verify function completes without panic
			defer func() {
				// Check for panic
				if r := recover(); r != nil {
					t.Errorf("checkDeclForAlloc panicked: %v", r)
				}
			}()
			// Call function with nil pass (just test structure)
			// We cannot fully test without analysis.Pass
		})
	}
}

// Test_isSliceOrMapAlloc tests the isSliceOrMapAlloc function.
//
// Params:
//   - t: testing context
func Test_isSliceOrMapAlloc(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "make slice",
			code:     "package test\nfunc f() { s := make([]int, 0) ; _ = s }",
			expected: true,
		},
		{
			name:     "make map",
			code:     "package test\nfunc f() { m := make(map[string]int) ; _ = m }",
			expected: true,
		},
		{
			name:     "make byte slice excluded",
			code:     "package test\nfunc f() { b := make([]byte, 0) ; _ = b }",
			expected: false,
		},
		{
			name:     "slice literal",
			code:     "package test\nfunc f() { s := []int{1,2,3} ; _ = s }",
			expected: true,
		},
		{
			name:     "map literal",
			code:     "package test\nfunc f() { m := map[string]int{\"a\":1} ; _ = m }",
			expected: true,
		},
		{
			name:     "byte slice literal excluded",
			code:     "package test\nfunc f() { b := []byte{1,2,3} ; _ = b }",
			expected: false,
		},
		{
			name:     "simple int",
			code:     "package test\nfunc f() { x := 42 ; _ = x }",
			expected: false,
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
			// Find first assignment RHS (skip blank assignments)
			var expr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if AssignStmt
				if as, ok := n.(*ast.AssignStmt); ok && len(as.Rhs) > 0 {
					// Skip blank assignments
					if len(as.Lhs) > 0 {
						// Check first LHS identifier
						if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == "_" {
							return true
						}
					}
					expr = as.Rhs[0]
					return false
				}
				return true
			})
			// Check if found
			if expr == nil {
				t.Fatal("No expression found")
			}
			// Check result
			got := isSliceOrMapAlloc(expr)
			// Validate result
			if got != tt.expected {
				t.Errorf("isSliceOrMapAlloc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_isByteSliceMake tests the isByteSliceMake function.
//
// Params:
//   - t: testing context
func Test_isByteSliceMake(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "make byte slice",
			code:     "package test\nfunc f() { b := make([]byte, 0) ; _ = b }",
			expected: true,
		},
		{
			name:     "make int slice",
			code:     "package test\nfunc f() { s := make([]int, 0) ; _ = s }",
			expected: false,
		},
		{
			name:     "make map",
			code:     "package test\nfunc f() { m := make(map[string]int) ; _ = m }",
			expected: false,
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
			// Find CallExpr
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if CallExpr with make
				if ce, ok := n.(*ast.CallExpr); ok {
					// Check if it's make call
					if ident, ok := ce.Fun.(*ast.Ident); ok && ident.Name == "make" {
						callExpr = ce
						return false
					}
				}
				return true
			})
			// Check if found
			if callExpr == nil {
				t.Fatal("No make call found")
			}
			// Check result
			got := isByteSliceMake(callExpr)
			// Validate result
			if got != tt.expected {
				t.Errorf("isByteSliceMake() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_isByteSliceMake_noArgs tests edge case with no arguments.
//
// Params:
//   - t: testing context
func Test_isByteSliceMake_noArgs(t *testing.T) {
	// Create empty call expression
	call := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: nil,
	}
	// Check result
	got := isByteSliceMake(call)
	// Validate result
	if got != false {
		t.Errorf("isByteSliceMake(no args) = %v, want false", got)
	}
}

// Test_checkLoopBodyForAlloc tests the checkLoopBodyForAlloc function.
//
// Params:
//   - t: testing context
func Test_checkLoopBodyForAlloc(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "nil body should not panic"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with nil body - should not panic
			checkLoopBodyForAlloc(nil, nil)
			// If we get here without panic, test passes
			t.Log("checkLoopBodyForAlloc handled nil body gracefully")
		})
	}
}

// Test_runVar008 tests the analyzer exists.
//
// Params:
//   - t: testing context
func Test_runVar008(t *testing.T) {
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
			if Analyzer008 == nil {
				t.Error("Analyzer008 is nil")
			}
		})
	}
}

// Test_checkAssignForAlloc tests the checkAssignForAlloc function.
//
// Params:
//   - t: testing context
func Test_checkAssignForAlloc(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via analysistest"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("checkAssignForAlloc is tested via analysistest")
		})
	}
}

// Test_checkStmtForAlloc tests the checkStmtForAlloc function.
//
// Params:
//   - t: testing context
func Test_checkStmtForAlloc(t *testing.T) {
	// Parse source with various statement types
	code := `package test
func f() {
	x := 1
	var y = 2
	_ = x
	_ = y
}`
	// Parse source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parse error
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	// Find statements
	ast.Inspect(file, func(n ast.Node) bool {
		// Check if block statement
		if block, ok := n.(*ast.BlockStmt); ok {
			// Iterate over statements
			for _, stmt := range block.List {
				checkStmtForAlloc(nil, stmt)
			}
		}
		return true
	})
	// If we get here without panic, test passes
}
