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

func TestIsAnyConstraint(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "any constraint",
			expr:     &ast.Ident{Name: "any"},
			expected: true,
		},
		{
			name:     "comparable constraint",
			expr:     &ast.Ident{Name: "comparable"},
			expected: false,
		},
		{
			name:     "other identifier",
			expr:     &ast.Ident{Name: "Reader"},
			expected: false,
		},
		{
			name: "selector expression (not ident)",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAnyConstraint(tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isAnyConstraint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractTypeName(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple identifier",
			expr:     &ast.Ident{Name: "T"},
			expected: "T",
		},
		{
			name: "array type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "T"},
			},
			expected: "T",
		},
		{
			name: "nested array type",
			expr: &ast.ArrayType{
				Elt: &ast.ArrayType{
					Elt: &ast.Ident{Name: "T"},
				},
			},
			expected: "T",
		},
		{
			name:     "star expression (unsupported)",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "T"}},
			expected: "",
		},
		{
			name: "map type (unsupported)",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTypeName(tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("extractTypeName() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestCheckOperandUsesAnyType(t *testing.T) {
	paramNames := map[string]string{
		"x": "T",
		"y": "T",
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
			name: "index expression with any type",
			expr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "x"},
				Index: &ast.Ident{Name: "i"},
			},
			expected: true,
		},
		{
			name: "index expression without any type",
			expr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "z"},
				Index: &ast.Ident{Name: "i"},
			},
			expected: false,
		},
		{
			name: "call expression (unsupported)",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "foo"},
			},
			expected: false,
		},
		{
			name: "binary expression (unsupported)",
			expr: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
				Op: 0,
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

func TestCollectAnyTypeParams(t *testing.T) {
	tests := []struct {
		name       string
		typeParams *ast.FieldList
		expected   map[string]bool
	}{
		{
			name: "single any param",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "T"}},
						Type:  &ast.Ident{Name: "any"},
					},
				},
			},
			expected: map[string]bool{"T": true},
		},
		{
			name: "multiple any params",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "T"}, {Name: "U"}},
						Type:  &ast.Ident{Name: "any"},
					},
				},
			},
			expected: map[string]bool{"T": true, "U": true},
		},
		{
			name: "comparable param (not any)",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "T"}},
						Type:  &ast.Ident{Name: "comparable"},
					},
				},
			},
			expected: map[string]bool{},
		},
		{
			name: "mixed params",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "T"}},
						Type:  &ast.Ident{Name: "any"},
					},
					{
						Names: []*ast.Ident{{Name: "K"}},
						Type:  &ast.Ident{Name: "comparable"},
					},
				},
			},
			expected: map[string]bool{"T": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectAnyTypeParams(tt.typeParams)
			// Verify length
			if len(result) != len(tt.expected) {
				t.Errorf("collectAnyTypeParams() length = %d, want %d", len(result), len(tt.expected))
			}
			// Verify values
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("collectAnyTypeParams()[%s] = %v, want %v", k, result[k], v)
				}
			}
		})
	}
}

func TestCollectParamNamesWithAnyType(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		expected map[string]string
	}{
		{
			name: "nil params",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: nil,
				},
			},
			expected: map[string]string{},
		},
		{
			name: "empty params list",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: nil},
				},
			},
			expected: map[string]string{},
		},
		{
			name: "param with any type",
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
			},
			expected: map[string]string{"x": "T"},
		},
		{
			name: "param without any type",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "x"}},
								Type:  &ast.Ident{Name: "int"},
							},
						},
					},
				},
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectParamNamesWithAnyType(tt.funcDecl, anyTypeParams)
			// Verify length
			if len(result) != len(tt.expected) {
				t.Errorf("collectParamNamesWithAnyType() length = %d, want %d", len(result), len(tt.expected))
			}
			// Verify values
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("collectParamNamesWithAnyType()[%s] = %s, want %s", k, result[k], v)
				}
			}
		})
	}
}

func TestAnalyzeGenericFunc(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "non-generic function (nil TypeParams)",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: nil,
				},
			},
		},
		{
			name: "generic function with comparable (not any)",
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
			// Just verify no panic - function returns nothing useful
			analyzeGenericFunc(nil, tt.funcDecl)
		})
	}
}

func TestCheckEqualityUsage(t *testing.T) {
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
			name: "empty params (no any type usage)",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "x"}},
								Type:  &ast.Ident{Name: "int"},
							},
						},
					},
				},
				Body: &ast.BlockStmt{List: nil},
			},
		},
		{
			name: "nil params list",
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
			checkEqualityUsage(nil, tt.funcDecl, anyTypeParams)
		})
	}
}

// TestCheckEqualityUsageEmptyParamNames tests when paramNames map is empty.
func TestCheckEqualityUsageEmptyParamNames(t *testing.T) {
	emptyAnyTypeParams := map[string]bool{}

	funcDecl := &ast.FuncDecl{
		Name: &ast.Ident{Name: "foo"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "x"}},
						Type:  &ast.Ident{Name: "int"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{List: nil},
	}

	// Should return early when paramNames is empty
	checkEqualityUsage(nil, funcDecl, emptyAnyTypeParams)
}

// TestCheckEqualityUsageWithBinaryExprNotAnyType tests checkEqualityUsage with binary expressions
// that don't use any type (to avoid nil pass panic).
func TestCheckEqualityUsageWithBinaryExprNotAnyType(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	// Function with equality comparison but using non-T type variable
	funcDecl := &ast.FuncDecl{
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "z"},
						Op: token.EQL,
						Y:  &ast.Ident{Name: "w"},
					},
				},
			},
		},
	}

	// Test with nil pass - should not panic (no any type used)
	checkEqualityUsage(nil, funcDecl, anyTypeParams)
}

// TestCheckEqualityUsageNonEqualityOp tests checkEqualityUsage with non-equality operators.
func TestCheckEqualityUsageNonEqualityOp(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	// Function with non-equality comparison (ADD)
	funcDecl := &ast.FuncDecl{
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "x"},
						Op: token.ADD,
						Y:  &ast.Ident{Name: "x"},
					},
				},
			},
		},
	}

	// Test should not report for non-equality operator
	checkEqualityUsage(nil, funcDecl, anyTypeParams)
}

// TestCheckEqualityUsageNonBinaryNode tests checkEqualityUsage with non-binary AST nodes.
func TestCheckEqualityUsageNonBinaryNode(t *testing.T) {
	anyTypeParams := map[string]bool{"T": true}

	// Function with non-binary statements
	funcDecl := &ast.FuncDecl{
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: "println"},
					},
				},
			},
		},
	}

	// Test with nil pass - should not panic (non-binary node)
	checkEqualityUsage(nil, funcDecl, anyTypeParams)
}

// Test_runGeneric001_disabled tests behavior when rule is disabled.
func Test_runGeneric001_disabled(t *testing.T) {
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
					"KTN-GENERIC-001": {Enabled: config.Bool(false)},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			// Creer un pass minimal
			result, err := runGeneric001(&analysis.Pass{})
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

// Test_runGeneric001_excludedFile tests behavior with excluded files.
func Test_runGeneric001_excludedFile(t *testing.T) {
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
					"KTN-GENERIC-001": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			code := `package test
func foo[T any](a, b T) bool { return a == b }
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
			result, err := runGeneric001(pass)
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

// TestReportIfUsesAnyTypeParam tests the reportIfUsesAnyTypeParam function.
func TestReportIfUsesAnyTypeParam(t *testing.T) {
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
				X: &ast.Ident{Name: "z"},
				Y: &ast.Ident{Name: "w"},
			},
			reported:   make(map[string]bool),
			expectCall: false,
		},
		{
			name: "already reported",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "y"},
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
			reportIfUsesAnyTypeParam(nil, funcDecl, tt.binaryExpr, paramNames, anyTypeParams, tt.reported)
		})
	}
}
