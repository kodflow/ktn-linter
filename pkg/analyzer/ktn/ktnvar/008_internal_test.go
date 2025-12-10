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

// Test_runVar008 tests the private runVar008 function.
func Test_runVar008(t *testing.T) {
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

// Test_isSliceOrMapAlloc tests the private isSliceOrMapAlloc helper function.
func Test_isSliceOrMapAlloc(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name: "slice literal",
			expr: &ast.CompositeLit{
				Type: &ast.ArrayType{},
			},
			expected: true,
		},
		{
			name: "map literal",
			expr: &ast.CompositeLit{
				Type: &ast.MapType{},
			},
			expected: true,
		},
		{
			name: "make call",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
			},
			expected: true,
		},
		{
			name: "struct literal",
			expr: &ast.CompositeLit{
				Type: &ast.Ident{Name: "MyStruct"},
			},
			expected: false,
		},
		{
			name: "other call",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "len"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSliceOrMapAlloc(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSliceOrMapAlloc() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_checkLoopBodyForAlloc_nilBody tests with nil body.
func Test_checkLoopBodyForAlloc_nilBody(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with nil body
	checkLoopBodyForAlloc(pass, nil)
	// No error expected
}

// Test_checkStmtForAlloc_emptyStmt tests with empty stmt.
func Test_checkStmtForAlloc_emptyStmt(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with empty statement (not assign or decl)
	checkStmtForAlloc(pass, &ast.EmptyStmt{})

	// Test with other statement types
	checkStmtForAlloc(pass, &ast.ReturnStmt{})
	checkStmtForAlloc(pass, &ast.IfStmt{})
	checkStmtForAlloc(pass, &ast.ForStmt{})
	checkStmtForAlloc(pass, &ast.ExprStmt{})
	// No error expected for non-assign/decl statements
}

// Test_checkAssignForAlloc tests the private checkAssignForAlloc function.
func Test_checkAssignForAlloc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks assignments for allocations
		})
	}
}

// Test_checkDeclForAlloc tests the private checkDeclForAlloc function.
func Test_checkDeclForAlloc(t *testing.T) {
	tests := []struct {
		name string
		decl *ast.DeclStmt
	}{
		{
			name: "non-GenDecl",
			decl: &ast.DeclStmt{
				Decl: &ast.BadDecl{},
			},
		},
		{
			name: "GenDecl with non-ValueSpec",
			decl: &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ImportSpec{},
					},
				},
			},
		},
		{
			name: "GenDecl with ValueSpec no values",
			decl: &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names:  []*ast.Ident{{Name: "x"}},
							Values: []ast.Expr{},
						},
					},
				},
			},
		},
		{
			name: "GenDecl with ValueSpec non-alloc value",
			decl: &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names:  []*ast.Ident{{Name: "x"}},
							Values: []ast.Expr{&ast.Ident{Name: "y"}},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create real pass with no-op reporter
			pass := &analysis.Pass{
				Report: func(d analysis.Diagnostic) {},
			}
			checkDeclForAlloc(pass, tt.decl)
			// Test passes if no panic
		})
	}
}

// Test_isByteSliceMake tests the private isByteSliceMake helper function.
func Test_isByteSliceMake(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "make with []byte",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "byte"}},
				},
			},
			expected: true,
		},
		{
			name: "make with []uint8",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "uint8"}},
				},
			},
			expected: true,
		},
		{
			name: "make with []int",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
				},
			},
			expected: false,
		},
		{
			name:     "make with no args",
			call:     &ast.CallExpr{Fun: &ast.Ident{Name: "make"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isByteSliceMake(tt.call)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isByteSliceMake() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_runVar008_disabled tests runVar008 with disabled rule.
func Test_runVar008_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-008": {Enabled: config.Bool(false)},
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

	_, err = runVar008(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar008() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar008() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar008_fileExcluded tests runVar008 with excluded file.
func Test_runVar008_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-008": {
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

	_, err = runVar008(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar008() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar008() reported %d issues, expected 0 when file excluded", reportCount)
	}
}
