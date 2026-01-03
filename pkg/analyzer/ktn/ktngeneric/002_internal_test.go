package ktngeneric

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestExtractConstraintName(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple identifier",
			expr:     &ast.Ident{Name: "Reader"},
			expected: "Reader",
		},
		{
			name: "selector expression",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
			expected: "io.Reader",
		},
		{
			name:     "unknown type returns interface",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "Reader"}},
			expected: "interface",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractConstraintName(tt.expr)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("extractConstraintName() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestExtractSelectorName(t *testing.T) {
	tests := []struct {
		name     string
		sel      *ast.SelectorExpr
		expected string
	}{
		{
			name: "package.Type",
			sel: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
			expected: "io.Reader",
		},
		{
			name: "nested selector (returns only selector)",
			sel: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "pkg"},
					Sel: &ast.Ident{Name: "sub"},
				},
				Sel: &ast.Ident{Name: "Type"},
			},
			expected: "Type",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractSelectorName(tt.sel)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("extractSelectorName() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestIsGenericBuiltinConstraint(t *testing.T) {
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
			expected: true,
		},
		{
			name:     "other identifier",
			expr:     &ast.Ident{Name: "Reader"},
			expected: false,
		},
		{
			name: "selector expression",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isGenericBuiltinConstraint(tt.expr)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("isGenericBuiltinConstraint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainsTypeParam(t *testing.T) {
	tests := []struct {
		name          string
		expr          ast.Expr
		typeParamName string
		expected      bool
	}{
		{
			name:          "simple match",
			expr:          &ast.Ident{Name: "T"},
			typeParamName: "T",
			expected:      true,
		},
		{
			name:          "no match",
			expr:          &ast.Ident{Name: "U"},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "nested match in array type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "T"},
			},
			typeParamName: "T",
			expected:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := containsTypeParam(tt.expr, tt.typeParamName)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("containsTypeParam() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckTypeAssertionUsage(t *testing.T) {
	tests := []struct {
		name          string
		node          ast.Node
		typeParamName string
		expected      bool
	}{
		{
			name: "type assertion with type param",
			node: &ast.TypeAssertExpr{
				X:    &ast.Ident{Name: "x"},
				Type: &ast.Ident{Name: "T"},
			},
			typeParamName: "T",
			expected:      true,
		},
		{
			name: "type assertion without type param",
			node: &ast.TypeAssertExpr{
				X:    &ast.Ident{Name: "x"},
				Type: &ast.Ident{Name: "int"},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name:          "not a type assertion",
			node:          &ast.Ident{Name: "x"},
			typeParamName: "T",
			expected:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := checkTypeAssertionUsage(tt.node, tt.typeParamName)
			// Verify result
			if result != tt.expected {
				t.Errorf("checkTypeAssertionUsage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckTypeConversionUsage(t *testing.T) {
	tests := []struct {
		name          string
		node          ast.Node
		typeParamName string
		expected      bool
	}{
		{
			name: "call with type param as function",
			node: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "T"},
				Args: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			typeParamName: "T",
			expected:      true,
		},
		{
			name: "call with different function name",
			node: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "int"},
				Args: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "call with selector (not simple ident)",
			node: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "pkg"},
					Sel: &ast.Ident{Name: "Func"},
				},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name:          "not a call expression",
			node:          &ast.Ident{Name: "x"},
			typeParamName: "T",
			expected:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := checkTypeConversionUsage(tt.node, tt.typeParamName)
			// Verify result
			if result != tt.expected {
				t.Errorf("checkTypeConversionUsage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsTypeInReturnType(t *testing.T) {
	tests := []struct {
		name          string
		funcDecl      *ast.FuncDecl
		typeParamName string
		expected      bool
	}{
		{
			name: "no results",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: nil,
				},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "type param in return",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "T"}},
						},
					},
				},
			},
			typeParamName: "T",
			expected:      true,
		},
		{
			name: "type param not in return",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "error"}},
						},
					},
				},
			},
			typeParamName: "T",
			expected:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isTypeInReturnType(tt.funcDecl, tt.typeParamName)
			// Verify result
			if result != tt.expected {
				t.Errorf("isTypeInReturnType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsTypeUsedInBody(t *testing.T) {
	tests := []struct {
		name          string
		funcDecl      *ast.FuncDecl
		typeParamName string
		expected      bool
	}{
		{
			name: "nil body",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: nil,
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "empty body",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{List: nil},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "body with type assertion",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.TypeAssertExpr{
								X:    &ast.Ident{Name: "x"},
								Type: &ast.Ident{Name: "T"},
							},
						},
					},
				},
			},
			typeParamName: "T",
			expected:      true,
		},
		{
			name: "body with type conversion",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  &ast.Ident{Name: "T"},
								Args: []ast.Expr{&ast.Ident{Name: "x"}},
							},
						},
					},
				},
			},
			typeParamName: "T",
			expected:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isTypeUsedInBody(tt.funcDecl, tt.typeParamName)
			// Verify result
			if result != tt.expected {
				t.Errorf("isTypeUsedInBody() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsTypeParamOnlyForSignature(t *testing.T) {
	tests := []struct {
		name          string
		funcDecl      *ast.FuncDecl
		typeParamName string
		expected      bool
	}{
		{
			name: "used in return and body - justified",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "T"}},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  &ast.Ident{Name: "T"},
								Args: []ast.Expr{&ast.Ident{Name: "x"}},
							},
						},
					},
				},
			},
			typeParamName: "T",
			expected:      false,
		},
		{
			name: "not used in return - potentially unnecessary",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "error"}},
						},
					},
				},
				Body: &ast.BlockStmt{List: nil},
			},
			typeParamName: "T",
			expected:      true,
		},
		{
			name: "used in return only - justified",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "T"}},
						},
					},
				},
				Body: &ast.BlockStmt{List: nil},
			},
			typeParamName: "T",
			expected:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isTypeParamOnlyForSignature(tt.funcDecl, tt.typeParamName)
			// Verify result
			if result != tt.expected {
				t.Errorf("isTypeParamOnlyForSignature() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAnalyzeUnnecessaryGeneric(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "non-generic function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: nil,
				},
			},
		},
		{
			name: "generic with any constraint",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "T"}},
								Type:  &ast.Ident{Name: "any"},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			analyzeUnnecessaryGeneric(nil, tt.funcDecl)
		})
	}
}

// Test_runGeneric002 tests the main runGeneric002 function.
func Test_runGeneric002(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "generic function with interface constraint",
			code: `package test
import "io"
func foo[T io.Reader](r T) { r.Read(nil) }
`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configure rule as enabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-002": {Enabled: config.Bool(true)},
				},
			})
			defer config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Verification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Create inspector
			files := []*ast.File{file}
			inspectResult, inspErr := inspect.Analyzer.Run(&analysis.Pass{
				Fset:  fset,
				Files: files,
			})
			// Vérifier l'erreur d'inspect
			if inspErr != nil || inspectResult == nil {
				t.Fatalf("failed to run inspect analyzer: %v", inspErr)
			}

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					// Expected to report
				},
			}

			// Execute analyzer
			result, err := runGeneric002(pass)
			// Verification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification resultat nil
			if result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// Test_isInterfaceType tests the isInterfaceType function.
func Test_isInterfaceType(t *testing.T) {
	tests := []struct {
		name     string
		typ      types.Type
		expected bool
	}{
		{
			name:     "nil type",
			typ:      nil,
			expected: false,
		},
		{
			name:     "basic int type",
			typ:      types.Typ[types.Int],
			expected: false,
		},
		{
			name:     "basic string type",
			typ:      types.Typ[types.String],
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isInterfaceType(tt.typ)
			// Verification resultat
			if result != tt.expected {
				t.Errorf("isInterfaceType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_reportUnnecessaryGeneric tests the reportUnnecessaryGeneric function.
func Test_reportUnnecessaryGeneric(t *testing.T) {
	tests := []struct {
		name       string
		funcDecl   *ast.FuncDecl
		typeParam  string
		constraint ast.Expr
	}{
		{
			name: "simple case with nil pass",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Process"},
			},
			typeParam:  "T",
			constraint: &ast.Ident{Name: "Reader"},
		},
		{
			name: "selector constraint",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Process"},
			},
			typeParam: "T",
			constraint: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with nil pass
			reportUnnecessaryGeneric(nil, tt.funcDecl, tt.typeParam, tt.constraint)
		})
	}
}

// Test_runGeneric002_disabled tests behavior when rule is disabled.
func Test_runGeneric002_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec regle desactivee
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-002": {Enabled: config.Bool(false)},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			// Creer un pass minimal
			result, err := runGeneric002(&analysis.Pass{})
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

// Test_runGeneric002_excludedFile tests behavior with excluded files.
func Test_runGeneric002_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-002": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			code := `package test
import "io"
func foo[T io.Reader](r T) { r.Read(nil) }
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
			result, err := runGeneric002(pass)
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

// TestIsSingleInterfaceConstraint tests the isSingleInterfaceConstraint function.
func TestIsSingleInterfaceConstraint(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "any constraint - builtin",
			expr:     &ast.Ident{Name: "any"},
			expected: false,
		},
		{
			name:     "comparable constraint - builtin",
			expr:     &ast.Ident{Name: "comparable"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Pass nil for pass - function will return false for nil typeInfo
			result := isSingleInterfaceConstraint(nil, tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isSingleInterfaceConstraint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsSingleInterfaceConstraintEdgeCases tests edge cases for isSingleInterfaceConstraint.
func TestIsSingleInterfaceConstraintEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		pass     *analysis.Pass
		expr     ast.Expr
		expected bool
	}{
		{
			name: "nil TypesInfo",
			pass: &analysis.Pass{
				TypesInfo: nil,
			},
			expr:     &ast.Ident{Name: "Reader"},
			expected: false,
		},
		{
			name: "TypesInfo with empty Types map",
			pass: &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			},
			expr:     &ast.Ident{Name: "UnknownType"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isSingleInterfaceConstraint(tt.pass, tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isSingleInterfaceConstraint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestReportUnnecessaryGenericEdgeCases tests edge cases for reportUnnecessaryGeneric.
func TestReportUnnecessaryGenericEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		funcDecl   *ast.FuncDecl
		typeParam  string
		constraint ast.Expr
	}{
		{
			name: "nil pass",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
			},
			typeParam:  "T",
			constraint: &ast.Ident{Name: "Reader"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with nil pass
			reportUnnecessaryGeneric(nil, tt.funcDecl, tt.typeParam, tt.constraint)
		})
	}
}
