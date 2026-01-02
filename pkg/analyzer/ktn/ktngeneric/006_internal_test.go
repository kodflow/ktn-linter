package ktngeneric

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestIsOrderedOp(t *testing.T) {
	tests := []struct {
		name     string
		op       token.Token
		expected bool
	}{
		// Comparison operators (ordered)
		{name: "LSS (<)", op: token.LSS, expected: true},
		{name: "LEQ (<=)", op: token.LEQ, expected: true},
		{name: "GTR (>)", op: token.GTR, expected: true},
		{name: "GEQ (>=)", op: token.GEQ, expected: true},
		// Arithmetic operators
		{name: "ADD (+)", op: token.ADD, expected: true},
		{name: "SUB (-)", op: token.SUB, expected: true},
		{name: "MUL (*)", op: token.MUL, expected: true},
		{name: "QUO (/)", op: token.QUO, expected: true},
		{name: "REM (%)", op: token.REM, expected: true},
		// Non-ordered operators
		{name: "EQL (==)", op: token.EQL, expected: false},
		{name: "NEQ (!=)", op: token.NEQ, expected: false},
		{name: "LAND (&&)", op: token.LAND, expected: false},
		{name: "LOR (||)", op: token.LOR, expected: false},
		{name: "AND (&)", op: token.AND, expected: false},
		{name: "OR (|)", op: token.OR, expected: false},
		{name: "XOR (^)", op: token.XOR, expected: false},
		{name: "SHL (<<)", op: token.SHL, expected: false},
		{name: "SHR (>>)", op: token.SHR, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOrderedOp(tt.op)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("isOrderedOp(%v) = %v, want %v", tt.op, result, tt.expected)
			}
		})
	}
}

func TestMergeStringMaps(t *testing.T) {
	tests := []struct {
		name     string
		m1       map[string]string
		m2       map[string]string
		expected map[string]string
	}{
		{
			name:     "both empty",
			m1:       map[string]string{},
			m2:       map[string]string{},
			expected: map[string]string{},
		},
		{
			name:     "first empty",
			m1:       map[string]string{},
			m2:       map[string]string{"a": "T"},
			expected: map[string]string{"a": "T"},
		},
		{
			name:     "second empty",
			m1:       map[string]string{"a": "T"},
			m2:       map[string]string{},
			expected: map[string]string{"a": "T"},
		},
		{
			name:     "both non-empty",
			m1:       map[string]string{"a": "T"},
			m2:       map[string]string{"b": "U"},
			expected: map[string]string{"a": "T", "b": "U"},
		},
		{
			name:     "overlapping keys (m2 wins)",
			m1:       map[string]string{"a": "T"},
			m2:       map[string]string{"a": "U"},
			expected: map[string]string{"a": "U"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeStringMaps(tt.m1, tt.m2)
			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			// Check values
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("result[%s] = %s, want %s", k, result[k], v)
				}
			}
		})
	}
}

func TestAnalyzeGenericFuncOrdered(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "nil type params",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: nil,
				},
			},
		},
		{
			name: "no any type params",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "T"}},
								Type:  &ast.Ident{Name: "comparable"},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			analyzeGenericFuncOrdered(nil, tt.funcDecl)
		})
	}
}

func TestCheckOrderedUsage(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "nil body",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "x"}},
								Type:  &ast.Ident{Name: "T"},
							},
						},
					},
				},
				Body: nil,
			},
		},
		{
			name: "empty params",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: nil,
				},
				Body: &ast.BlockStmt{List: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			checkOrderedUsage(nil, tt.funcDecl, anyTypeParams)
		})
	}
}

func TestCollectLocalVarsWithAnyType(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		expected map[string]string
	}{
		{
			name: "nil body",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: nil,
			},
			expected: map[string]string{},
		},
		{
			name: "empty body",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{List: nil},
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectLocalVarsWithAnyType(tt.funcDecl, anyTypeParams)
			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
		})
	}
}

func TestExtractVarDeclsFromStmt(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}
	result := make(map[string]string)

	tests := []struct {
		name string
		stmt ast.Stmt
	}{
		{
			name: "expr statement (not handled)",
			stmt: &ast.ExprStmt{
				X: &ast.Ident{Name: "x"},
			},
		},
		{
			name: "assign statement (not handled)",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			extractVarDeclsFromStmt(tt.stmt, anyTypeParams, result)
		})
	}
}

func TestExtractFromDeclStmt(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}
	result := make(map[string]string)

	tests := []struct {
		name     string
		declStmt *ast.DeclStmt
	}{
		{
			name: "not a GenDecl",
			declStmt: &ast.DeclStmt{
				Decl: &ast.FuncDecl{
					Name: &ast.Ident{Name: "foo"},
				},
			},
		},
		{
			name: "GenDecl with non-ValueSpec",
			declStmt: &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ImportSpec{
							Path: &ast.BasicLit{Value: `"fmt"`},
						},
					},
				},
			},
		},
		{
			name: "GenDecl with ValueSpec not matching any type",
			declStmt: &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{Name: "x"}},
							Type:  &ast.Ident{Name: "int"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			extractFromDeclStmt(tt.declStmt, anyTypeParams, result)
		})
	}
}

func TestExtractFromRangeStmt(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}
	result := make(map[string]string)

	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
	}{
		{
			name: "nil value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: nil,
				X:     &ast.Ident{Name: "slice"},
			},
		},
		{
			name: "value not an ident",
			rangeStmt: &ast.RangeStmt{
				Key: &ast.Ident{Name: "i"},
				Value: &ast.IndexExpr{
					X:     &ast.Ident{Name: "arr"},
					Index: &ast.Ident{Name: "0"},
				},
				X: &ast.Ident{Name: "slice"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			extractFromRangeStmt(tt.rangeStmt, anyTypeParams, result)
		})
	}
}

func TestExtractArrayElementType(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple ident",
			expr:     &ast.Ident{Name: "slice"},
			expected: "",
		},
		{
			name: "call expression",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractArrayElementType(tt.expr, anyTypeParams)
			// Verify result
			if result != tt.expected {
				t.Errorf("extractArrayElementType() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestCheckOperandUsesAnyType006(t *testing.T) {
	paramNames := map[string]string{
		"x": "T",
	}
	anyTypeParams := map[string]bool{
		"T": true,
	}

	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "identifier in paramNames",
			expr:     &ast.Ident{Name: "x"},
			expected: true,
		},
		{
			name:     "identifier not in paramNames",
			expr:     &ast.Ident{Name: "z"},
			expected: false,
		},
		{
			name: "index expression",
			expr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "x"},
				Index: &ast.Ident{Name: "i"},
			},
			expected: true,
		},
		{
			name: "call expression (not supported)",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "foo"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkOperandUsesAnyType(tt.expr, paramNames, anyTypeParams)
			// Verify result
			if result != tt.expected {
				t.Errorf("checkOperandUsesAnyType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_runGeneric006_disabled tests behavior when rule is disabled.
func Test_runGeneric006_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec regle desactivee
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-006": {Enabled: config.Bool(false)},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			// Creer un pass minimal
			result, err := runGeneric006(&analysis.Pass{})
			// Verification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification du resultat nil
			if result != nil {
				t.Errorf("Expected nil result when rule disabled, got %v", result)
			}
		})
	}
}

// Test_runGeneric006_excludedFile tests behavior with excluded files.
func Test_runGeneric006_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-006": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			code := `package test
func foo[T any](a, b T) T { return a + b }
`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Verification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Creer un inspector
			files := []*ast.File{file}
			inspectResult, _ := inspect.Analyzer.Run(&analysis.Pass{
				Fset:  fset,
				Files: files,
			})

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Expected no diagnostics for excluded file, got: %s", d.Message)
				},
			}

			// Executer l'analyseur
			result, err := runGeneric006(pass)
			// Verification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification du resultat nil
			if result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestReportIfUsesAnyTypeParamOrdered tests the reportIfUsesAnyTypeParamOrdered function.
func TestReportIfUsesAnyTypeParamOrdered(t *testing.T) {
	// Reset config
	config.Reset()
	defer config.Reset()

	paramNames := map[string]string{"x": "T"}
	anyTypeParams := map[string]bool{"T": true}

	tests := []struct {
		name       string
		binaryExpr *ast.BinaryExpr
		reported   map[string]bool
		expectCall bool
	}{
		{
			name: "neither operand uses any type",
			binaryExpr: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "z"},
				Y:  &ast.Ident{Name: "w"},
				Op: token.ADD,
			},
			reported:   make(map[string]bool),
			expectCall: false,
		},
		{
			name: "already reported",
			binaryExpr: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.Ident{Name: "y"},
				Op: token.LSS,
			},
			reported:   map[string]bool{"foo": true},
			expectCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
			}
			// Just verify no panic
			reportIfUsesAnyTypeParamOrdered(nil, funcDecl, tt.binaryExpr, paramNames, anyTypeParams, tt.reported)
		})
	}
}

// TestExtractFromRangeStmtWithValue tests extractFromRangeStmt with various scenarios.
func TestExtractFromRangeStmtWithValue(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}
	result := make(map[string]string)

	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
	}{
		{
			name: "value with array type",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "v"},
				X: &ast.ArrayType{
					Elt: &ast.Ident{Name: "T"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			extractFromRangeStmt(tt.rangeStmt, anyTypeParams, result)
		})
	}
}
