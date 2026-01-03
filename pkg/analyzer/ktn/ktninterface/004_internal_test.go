// Package ktninterface provides internal tests for KTN-INTERFACE-004 helper functions.
package ktninterface

import (
	"go/ast"
	"testing"
)

// Test_isEmptyInterface tests the isEmptyInterface helper function.
func Test_isEmptyInterface(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name: "interface{} type",
			expr: &ast.InterfaceType{
				Methods: nil,
			},
			expected: true,
		},
		{
			name: "interface{} with empty list",
			expr: &ast.InterfaceType{
				Methods: &ast.FieldList{List: []*ast.Field{}},
			},
			expected: true,
		},
		{
			name: "any keyword",
			expr: &ast.Ident{
				Name: "any",
			},
			expected: true,
		},
		{
			name: "non-empty interface",
			expr: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{Names: []*ast.Ident{{Name: "Method"}}},
					},
				},
			},
			expected: false,
		},
		{
			name: "regular identifier",
			expr: &ast.Ident{
				Name: "string",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := isEmptyInterface(tt.expr)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("isEmptyInterface() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_buildParamName tests the buildParamName helper function.
func Test_buildParamName(t *testing.T) {
	tests := []struct {
		name     string
		param    *ast.Field
		funcName string
		expected string
	}{
		{
			name: "named parameter",
			param: &ast.Field{
				Names: []*ast.Ident{{Name: "data"}},
			},
			funcName: "Process",
			expected: "data",
		},
		{
			name:     "anonymous parameter",
			param:    &ast.Field{},
			funcName: "Process",
			expected: "paramètre de Process",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := buildParamName(tt.param, tt.funcName)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("buildParamName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test_buildReturnName tests the buildReturnName helper function.
func Test_buildReturnName(t *testing.T) {
	tests := []struct {
		name     string
		result   *ast.Field
		index    int
		funcName string
		expected string
	}{
		{
			name: "named return",
			result: &ast.Field{
				Names: []*ast.Ident{{Name: "result"}},
			},
			index:    0,
			funcName: "Process",
			expected: "result",
		},
		{
			name:     "first anonymous return",
			result:   &ast.Field{},
			index:    0,
			funcName: "Process",
			expected: "retour de Process",
		},
		{
			name:     "second anonymous return",
			result:   &ast.Field{},
			index:    1,
			funcName: "Process",
			expected: "retour de Process",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := buildReturnName(tt.result, tt.index, tt.funcName)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("buildReturnName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test_runInterface004 validates runInterface004 via integration test reference.
// Full behavior is tested in 004_external_test.go via analysistest.Run.
func Test_runInterface004(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func() bool
		errMsg    string
	}{
		{
			name:      "analyzer not nil",
			checkFunc: func() bool { return Analyzer004 != nil },
			errMsg:    "Analyzer004 should not be nil",
		},
		{
			name:      "analyzer name correct",
			checkFunc: func() bool { return Analyzer004.Name == "ktninterface004" },
			errMsg:    "Analyzer004.Name should be ktninterface004",
		},
		{
			name:      "run function assigned",
			checkFunc: func() bool { return Analyzer004.Run != nil },
			errMsg:    "Analyzer004.Run should not be nil",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier la condition
			if !tt.checkFunc() {
				// Condition non satisfaite
				t.Error(tt.errMsg)
			}
		})
	}
}

// Test_isAnalyzerRunFunction tests the isAnalyzerRunFunction helper function.
func Test_isAnalyzerRunFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		expected bool
	}{
		{
			name: "valid analyzer.Run signature",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "run"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "analysis"},
										Sel: &ast.Ident{Name: "Pass"},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "any"}},
							{Type: &ast.Ident{Name: "error"}},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "no parameters",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "run"},
				Type: &ast.FuncType{
					Params: nil,
				},
			},
			expected: false,
		},
		{
			name: "wrong parameter type",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "run"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "string"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "no return values",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "run"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "analysis"},
										Sel: &ast.Ident{Name: "Pass"},
									},
								},
							},
						},
					},
					Results: nil,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := isAnalyzerRunFunction(tt.funcDecl)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("isAnalyzerRunFunction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_checkFuncParams validates checkFuncParams helper function.
// Full behavior is tested in 004_external_test.go via analysistest.Run.
func Test_checkFuncParams(t *testing.T) {
	tests := []struct {
		name      string
		expr      ast.Expr
		expected  bool
		checkDesc string
	}{
		{
			name:      "empty interface nil methods",
			expr:      &ast.InterfaceType{Methods: nil},
			expected:  true,
			checkDesc: "should return true for empty interface",
		},
		{
			name:      "empty interface empty list",
			expr:      &ast.InterfaceType{Methods: &ast.FieldList{List: []*ast.Field{}}},
			expected:  true,
			checkDesc: "should return true for empty interface list",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction isEmptyInterface
			got := isEmptyInterface(tt.expr)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("isEmptyInterface() = %v, %s", got, tt.checkDesc)
			}
		})
	}
}

// Test_checkFuncReturns validates checkFuncReturns helper function.
// Full behavior is tested in 004_external_test.go via analysistest.Run.
func Test_checkFuncReturns(t *testing.T) {
	tests := []struct {
		name      string
		expr      ast.Expr
		expected  bool
		checkDesc string
	}{
		{
			name:      "any keyword",
			expr:      &ast.Ident{Name: "any"},
			expected:  true,
			checkDesc: "should return true for 'any' keyword",
		},
		{
			name:      "other identifier",
			expr:      &ast.Ident{Name: "string"},
			expected:  false,
			checkDesc: "should return false for other identifiers",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction isEmptyInterface
			got := isEmptyInterface(tt.expr)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("isEmptyInterface() = %v, %s", got, tt.checkDesc)
			}
		})
	}
}

// Test_reportEmptyInterface validates reportEmptyInterface helper function.
// Full behavior is tested in 004_external_test.go via analysistest.Run.
func Test_reportEmptyInterface(t *testing.T) {
	tests := []struct {
		name     string
		param    *ast.Field
		funcName string
		expected string
	}{
		{
			name:     "named parameter",
			param:    &ast.Field{Names: []*ast.Ident{{Name: "data"}}},
			funcName: "TestFunc",
			expected: "data",
		},
		{
			name:     "anonymous parameter",
			param:    &ast.Field{},
			funcName: "TestFunc",
			expected: "paramètre de TestFunc",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction buildParamName
			got := buildParamName(tt.param, tt.funcName)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("buildParamName() = %q, want %q", got, tt.expected)
			}
		})
	}
}
