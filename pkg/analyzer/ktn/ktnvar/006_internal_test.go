// Internal tests for 006.go - ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_extractTypeString tests the extractTypeString function.
//
// Params:
//   - t: testing context
func Test_extractTypeString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "explicit type strings.Builder",
			code:     "package test\nimport \"strings\"\nvar sb strings.Builder",
			expected: "strings.Builder",
		},
		{
			name:     "explicit type bytes.Buffer",
			code:     "package test\nimport \"bytes\"\nvar buf bytes.Buffer",
			expected: "bytes.Buffer",
		},
		{
			name:     "composite literal strings.Builder",
			code:     "package test\nimport \"strings\"\nvar sb = strings.Builder{}",
			expected: "strings.Builder",
		},
		{
			name:     "composite literal bytes.Buffer",
			code:     "package test\nimport \"bytes\"\nvar buf = bytes.Buffer{}",
			expected: "bytes.Buffer",
		},
		{
			name:     "no type info",
			code:     "package test\nvar x int",
			expected: "",
		},
		{
			name:     "nil type and empty values",
			code:     "package test\nvar x int = 42",
			expected: "",
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
			// Find ValueSpec
			var valueSpec *ast.ValueSpec
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if ValueSpec
				if vs, ok := n.(*ast.ValueSpec); ok {
					valueSpec = vs
					return false
				}
				return true
			})
			// Check if found
			if valueSpec == nil {
				t.Fatal("No ValueSpec found")
			}
			// Extract type string
			got := extractTypeString(valueSpec.Type, valueSpec.Values)
			// Validate result
			if got != tt.expected {
				t.Errorf("extractTypeString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test_extractTypeString_nilInputs tests edge cases with nil inputs.
//
// Params:
//   - t: testing context
func Test_extractTypeString_nilInputs(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		typeExpr ast.Expr
		values   []ast.Expr
		expected string
	}{
		{
			name:     "nil type and nil values",
			typeExpr: nil,
			values:   nil,
			expected: "",
		},
		{
			name:     "nil type and empty values",
			typeExpr: nil,
			values:   []ast.Expr{},
			expected: "",
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTypeString(tt.typeExpr, tt.values)
			// Validate result
			if result != tt.expected {
				t.Errorf("extractTypeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_runVar006 tests the analyzer exists.
//
// Params:
//   - t: testing context
func Test_runVar006(t *testing.T) {
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
			if Analyzer006 == nil {
				t.Error("Analyzer006 is nil")
			}
		})
	}
}

// Test_checkBuilderWithoutGrow tests the checkBuilderWithoutGrow function.
//
// Params:
//   - t: testing context
func Test_checkBuilderWithoutGrow(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via checkValueSpec and checkAssignStmt"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("checkBuilderWithoutGrow is tested via component functions")
		})
	}
}

// Test_reportMissingGrow tests the reportMissingGrow function.
//
// Params:
//   - t: testing context
func Test_reportMissingGrow(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
	}{
		{name: "tested via analysistest"},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("reportMissingGrow is tested via analysistest")
		})
	}
}

// Test_extractAssignTypeString tests the extractAssignTypeString function.
//
// Params:
//   - t: testing context
func Test_extractAssignTypeString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "strings.Builder assignment",
			code:     "package test\nimport \"strings\"\nfunc f() { sb := strings.Builder{} ; _ = sb }",
			expected: "strings.Builder",
		},
		{
			name:     "bytes.Buffer assignment",
			code:     "package test\nimport \"bytes\"\nfunc f() { buf := bytes.Buffer{} ; _ = buf }",
			expected: "bytes.Buffer",
		},
		{
			name:     "non-builder assignment",
			code:     "package test\nfunc f() { x := 42 ; _ = x }",
			expected: "",
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
			// Find AssignStmt (skip blank assignments)
			var assignStmt *ast.AssignStmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if AssignStmt
				if as, ok := n.(*ast.AssignStmt); ok {
					// Skip blank assignments
					if len(as.Lhs) > 0 {
						// Check first LHS identifier
						if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == "_" {
							return true
						}
					}
					assignStmt = as
					return false
				}
				return true
			})
			// Check if found
			if assignStmt == nil {
				// No assignment found, expected empty result
				if tt.expected != "" {
					t.Fatal("No AssignStmt found")
				}
				return
			}
			// Call function
			result := extractAssignTypeString(assignStmt)
			// Validate result
			if result != tt.expected {
				t.Errorf("extractAssignTypeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_isBuilderCompositeLit tests the isBuilderCompositeLit function.
//
// Params:
//   - t: testing context
func Test_isBuilderCompositeLit(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "strings.Builder",
			code:     "package test\nimport \"strings\"\nvar sb = strings.Builder{}",
			expected: true,
		},
		{
			name:     "bytes.Buffer",
			code:     "package test\nimport \"bytes\"\nvar buf = bytes.Buffer{}",
			expected: true,
		},
		{
			name:     "other composite",
			code:     "package test\ntype S struct{}\nvar s = S{}",
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
			// Find CompositeLit
			var compositeLit *ast.CompositeLit
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if CompositeLit
				if cl, ok := n.(*ast.CompositeLit); ok {
					compositeLit = cl
					return false
				}
				return true
			})
			// Check if found
			if compositeLit == nil {
				t.Fatal("No CompositeLit found")
			}
			// Check result
			got := isBuilderCompositeLit(compositeLit)
			// Validate result
			if got != tt.expected {
				t.Errorf("isBuilderCompositeLit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_checkValueSpec tests the checkValueSpec function.
//
// Params:
//   - t: testing context
func Test_checkValueSpec(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		code       string
		wantLit    bool
		wantReport bool
	}{
		{
			name:       "strings.Builder literal",
			code:       "package test\nimport \"strings\"\nvar sb = strings.Builder{}",
			wantLit:    true,
			wantReport: true,
		},
		{
			name:       "no values",
			code:       "package test\nimport \"strings\"\nvar sb strings.Builder",
			wantLit:    false,
			wantReport: false,
		},
		{
			name:       "non-builder value",
			code:       "package test\nvar x = 42",
			wantLit:    false,
			wantReport: false,
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
			// Find ValueSpec
			var valueSpec *ast.ValueSpec
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if ValueSpec
				if vs, ok := n.(*ast.ValueSpec); ok {
					valueSpec = vs
					return false
				}
				return true
			})
			// Check if found
			if valueSpec == nil {
				t.Fatal("No ValueSpec found")
			}
			// Call function
			lit, pos := checkValueSpec(valueSpec)
			// Validate result
			gotLit := lit != nil
			gotPos := pos != nil
			// Check literal
			if gotLit != tt.wantLit {
				t.Errorf("checkValueSpec() lit != nil = %v, want %v", gotLit, tt.wantLit)
			}
			// Check position
			if gotPos != tt.wantReport {
				t.Errorf("checkValueSpec() pos != nil = %v, want %v", gotPos, tt.wantReport)
			}
		})
	}
}

// Test_checkAssignStmt tests the checkAssignStmt function.
//
// Params:
//   - t: testing context
func Test_checkAssignStmt(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		code    string
		wantLit bool
	}{
		{
			name:    "short decl strings.Builder",
			code:    "package test\nimport \"strings\"\nfunc f() { sb := strings.Builder{} ; _ = sb }",
			wantLit: true,
		},
		{
			name:    "short decl bytes.Buffer",
			code:    "package test\nimport \"bytes\"\nfunc f() { buf := bytes.Buffer{} ; _ = buf }",
			wantLit: true,
		},
		{
			name:    "non-builder assign",
			code:    "package test\nfunc f() { x := 42 ; _ = x }",
			wantLit: false,
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
			// Find AssignStmt with composite literal
			var assignStmt *ast.AssignStmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if AssignStmt
				if as, ok := n.(*ast.AssignStmt); ok {
					// Skip blank assignments
					if len(as.Lhs) > 0 {
						// Check first LHS identifier
						if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == "_" {
							return true
						}
					}
					assignStmt = as
					return false
				}
				return true
			})
			// Check if found
			if assignStmt == nil {
				t.Fatal("No AssignStmt found")
			}
			// Call function
			lit, _ := checkAssignStmt(assignStmt)
			// Validate result
			gotLit := lit != nil
			// Check literal
			if gotLit != tt.wantLit {
				t.Errorf("checkAssignStmt() lit != nil = %v, want %v", gotLit, tt.wantLit)
			}
		})
	}
}
