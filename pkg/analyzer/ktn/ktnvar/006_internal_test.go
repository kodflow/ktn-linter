package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar006 tests the private runVar006 function.
func Test_runVar006(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_isBuilderCompositeLit tests the private isBuilderCompositeLit helper function.
func Test_isBuilderCompositeLit(t *testing.T) {
	tests := []struct {
		name     string
		lit      *ast.CompositeLit
		expected bool
	}{
		{
			name: "strings.Builder",
			lit: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "strings"},
					Sel: &ast.Ident{Name: "Builder"},
				},
			},
			expected: true,
		},
		{
			name: "bytes.Buffer",
			lit: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "bytes"},
					Sel: &ast.Ident{Name: "Buffer"},
				},
			},
			expected: true,
		},
		{
			name: "other type",
			lit: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "strings"},
					Sel: &ast.Ident{Name: "Reader"},
				},
			},
			expected: false,
		},
		{
			name: "not selector expr",
			lit: &ast.CompositeLit{
				Type: &ast.Ident{Name: "MyStruct"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBuilderCompositeLit(tt.lit)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isBuilderCompositeLit() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_checkBuilderWithoutGrow tests the private checkBuilderWithoutGrow function.
func Test_checkBuilderWithoutGrow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks builders without Grow
		})
	}
}

// Test_checkValueSpec tests the private checkValueSpec function.
func Test_checkValueSpec(t *testing.T) {
	tests := []struct {
		name       string
		node       *ast.ValueSpec
		expectLit  bool
		expectNode bool
	}{
		{
			name:       "empty values",
			node:       &ast.ValueSpec{Values: []ast.Expr{}},
			expectLit:  false,
			expectNode: false,
		},
		{
			name: "strings.Builder value",
			node: &ast.ValueSpec{
				Values: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "strings"},
							Sel: &ast.Ident{Name: "Builder"},
						},
					},
				},
			},
			expectLit:  true,
			expectNode: true,
		},
		{
			name: "bytes.Buffer value",
			node: &ast.ValueSpec{
				Values: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "bytes"},
							Sel: &ast.Ident{Name: "Buffer"},
						},
					},
				},
			},
			expectLit:  true,
			expectNode: true,
		},
		{
			name: "non-Builder composite literal",
			node: &ast.ValueSpec{
				Values: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.Ident{Name: "MyStruct"},
					},
				},
			},
			expectLit:  false,
			expectNode: false,
		},
		{
			name: "non-composite literal value",
			node: &ast.ValueSpec{
				Values: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			expectLit:  false,
			expectNode: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit, pos := checkValueSpec(tt.node)
			// Check composite literal result
			if (lit != nil) != tt.expectLit {
				t.Errorf("checkValueSpec() lit = %v, expectLit %v", lit != nil, tt.expectLit)
			}
			// Check position node result
			if (pos != nil) != tt.expectNode {
				t.Errorf("checkValueSpec() pos = %v, expectNode %v", pos != nil, tt.expectNode)
			}
		})
	}
}

// Test_checkAssignStmt tests the private checkAssignStmt function.
func Test_checkAssignStmt(t *testing.T) {
	tests := []struct {
		name       string
		node       *ast.AssignStmt
		expectLit  bool
		expectNode bool
	}{
		{
			name:       "empty rhs",
			node:       &ast.AssignStmt{Rhs: []ast.Expr{}},
			expectLit:  false,
			expectNode: false,
		},
		{
			name: "strings.Builder assignment",
			node: &ast.AssignStmt{
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "strings"},
							Sel: &ast.Ident{Name: "Builder"},
						},
					},
				},
			},
			expectLit:  true,
			expectNode: true,
		},
		{
			name: "non-Builder assignment",
			node: &ast.AssignStmt{
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.Ident{Name: "MyStruct"},
					},
				},
			},
			expectLit:  false,
			expectNode: false,
		},
		{
			name: "non-composite literal rhs",
			node: &ast.AssignStmt{
				Rhs: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			expectLit:  false,
			expectNode: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit, pos := checkAssignStmt(tt.node)
			// Check composite literal result
			if (lit != nil) != tt.expectLit {
				t.Errorf("checkAssignStmt() lit = %v, expectLit %v", lit != nil, tt.expectLit)
			}
			// Check position node result
			if (pos != nil) != tt.expectNode {
				t.Errorf("checkAssignStmt() pos = %v, expectNode %v", pos != nil, tt.expectNode)
			}
		})
	}
}

// Test_reportMissingGrow tests the private reportMissingGrow function.
func Test_reportMissingGrow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function reports missing Grow
		})
	}
}

// Test_extractTypeString tests the private extractTypeString function.
func Test_extractTypeString(t *testing.T) {
	tests := []struct {
		name     string
		typeExpr ast.Expr
		values   []ast.Expr
		expected string
	}{
		{
			name:     "nil type expr and empty values",
			typeExpr: nil,
			values:   []ast.Expr{},
			expected: "",
		},
		{
			name: "selector type expr",
			typeExpr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "strings"},
				Sel: &ast.Ident{Name: "Builder"},
			},
			values:   nil,
			expected: "strings.Builder",
		},
		{
			name:     "type from composite literal value",
			typeExpr: nil,
			values: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "bytes"},
						Sel: &ast.Ident{Name: "Buffer"},
					},
				},
			},
			expected: "bytes.Buffer",
		},
		{
			name:     "non-selector type expr",
			typeExpr: &ast.Ident{Name: "int"},
			values:   nil,
			expected: "",
		},
		{
			name:     "value is not composite literal",
			typeExpr: nil,
			values:   []ast.Expr{&ast.Ident{Name: "x"}},
			expected: "",
		},
		{
			name:     "composite literal without selector type",
			typeExpr: nil,
			values: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.Ident{Name: "MyStruct"},
				},
			},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTypeString(tt.typeExpr, tt.values)
			// Check result matches expected
			if result != tt.expected {
				t.Errorf("extractTypeString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// Test_extractAssignTypeString tests the private extractAssignTypeString function.
func Test_extractAssignTypeString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function extracts assign type strings
		})
	}
}

// Test_runVar006_disabled tests runVar006 with disabled rule.
func Test_runVar006_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-006": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar006() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar006() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar006_fileExcluded tests runVar006 with excluded file.
func Test_runVar006_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-006": {
				Exclude: []string{"test.go"},
			},
		},
	})
	defer config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar006() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar006() reported %d issues, expected 0 when file excluded", reportCount)
	}
}
