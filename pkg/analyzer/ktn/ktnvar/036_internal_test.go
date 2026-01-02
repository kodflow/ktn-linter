package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
)

// Test_functionReturnsInt tests detection of function returning int.
func Test_functionReturnsInt(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		expected bool
	}{
		{
			name: "no results",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: nil,
				},
			},
			expected: false,
		},
		{
			name: "multiple results",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "int"}},
							{Type: &ast.Ident{Name: "error"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "result not ident",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.StarExpr{X: &ast.Ident{Name: "int"}}},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "result not int",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "string"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "valid int return",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ident{Name: "int"}},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := functionReturnsInt(tt.funcDecl)
			// Verify result
			if result != tt.expected {
				t.Errorf("functionReturnsInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isReturnMinusOne tests detection of return -1.
func Test_isReturnMinusOne(t *testing.T) {
	tests := []struct {
		name     string
		stmt     ast.Stmt
		expected bool
	}{
		{
			name:     "not return statement",
			stmt:     &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			expected: false,
		},
		{
			name: "no results",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{},
			},
			expected: false,
		},
		{
			name: "multiple results",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.INT, Value: "1"}},
					&ast.Ident{Name: "nil"},
				},
			},
			expected: false,
		},
		{
			name: "not unary expression",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "1"},
				},
			},
			expected: false,
		},
		{
			name: "not SUB operator",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.ADD, X: &ast.BasicLit{Kind: token.INT, Value: "1"}},
				},
			},
			expected: false,
		},
		{
			name: "operand not basic lit",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.SUB, X: &ast.Ident{Name: "x"}},
				},
			},
			expected: false,
		},
		{
			name: "operand not INT",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.FLOAT, Value: "1.0"}},
				},
			},
			expected: false,
		},
		{
			name: "operand not 1",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.INT, Value: "2"}},
				},
			},
			expected: false,
		},
		{
			name: "valid return -1",
			stmt: &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.INT, Value: "1"}},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isReturnMinusOne(tt.stmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("isReturnMinusOne() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractRangeVariables036 tests extraction of range variables.
func Test_extractRangeVariables036(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		wantNil   bool
	}{
		{
			name: "nil key",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			wantNil: true,
		},
		{
			name: "nil value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: nil,
			},
			wantNil: true,
		},
		{
			name: "key not ident",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
				Value: &ast.Ident{Name: "v"},
			},
			wantNil: true,
		},
		{
			name: "value not ident",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			},
			wantNil: true,
		},
		{
			name: "valid range variables",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "v"},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractRangeVariables036(tt.rangeStmt)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("extractRangeVariables036() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_findIndexPatternInBody036 tests finding index pattern in body.
func Test_findIndexPatternInBody036(t *testing.T) {
	vars := &rangeVariables036{indexName: "i", valueName: "v"}

	tests := []struct {
		name    string
		body    *ast.BlockStmt
		vars    *rangeVariables036
		wantNil bool
	}{
		{
			name:    "nil body",
			body:    nil,
			vars:    vars,
			wantNil: true,
		},
		{
			name:    "empty body",
			body:    &ast.BlockStmt{List: []ast.Stmt{}},
			vars:    vars,
			wantNil: true,
		},
		{
			name: "no matching statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
			vars:    vars,
			wantNil: true,
		},
		{
			name: "matching if statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.EQL,
						X:  &ast.Ident{Name: "v"},
						Y:  &ast.Ident{Name: "target"},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{
						&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "i"}}},
					}},
				},
			}},
			vars:    vars,
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := findIndexPatternInBody036(tt.body, tt.vars)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("findIndexPatternInBody036() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_checkIfEqualReturnIndex tests checking if statement pattern.
func Test_checkIfEqualReturnIndex(t *testing.T) {
	tests := []struct {
		name      string
		stmt      ast.Stmt
		indexName string
		valueName string
		expected  bool
	}{
		{
			name:      "not if statement",
			stmt:      &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			indexName: "i",
			valueName: "v",
			expected:  false,
		},
		{
			name: "cond not binary",
			stmt: &ast.IfStmt{
				Cond: &ast.Ident{Name: "cond"},
				Body: &ast.BlockStmt{},
			},
			indexName: "i",
			valueName: "v",
			expected:  false,
		},
		{
			name: "not equality operator",
			stmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{},
			},
			indexName: "i",
			valueName: "v",
			expected:  false,
		},
		{
			name: "value not in comparison",
			stmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.Ident{Name: "y"},
				},
				Body: &ast.BlockStmt{},
			},
			indexName: "i",
			valueName: "v",
			expected:  false,
		},
		{
			name: "body does not return index",
			stmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "other"}}},
				}},
			},
			indexName: "i",
			valueName: "v",
			expected:  false,
		},
		{
			name: "valid pattern",
			stmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "i"}}},
				}},
			},
			indexName: "i",
			valueName: "v",
			expected:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := checkIfEqualReturnIndex(tt.stmt, tt.indexName, tt.valueName)
			// Verify result
			if result != tt.expected {
				t.Errorf("checkIfEqualReturnIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_expressionUsesIdent tests detection of identifier usage.
func Test_expressionUsesIdent(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		identNme string
		expected bool
	}{
		{
			name:     "not identifier",
			expr:     &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			identNme: "x",
			expected: false,
		},
		{
			name:     "wrong name",
			expr:     &ast.Ident{Name: "y"},
			identNme: "x",
			expected: false,
		},
		{
			name:     "matching name",
			expr:     &ast.Ident{Name: "x"},
			identNme: "x",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := expressionUsesIdent(tt.expr, tt.identNme)
			// Verify result
			if result != tt.expected {
				t.Errorf("expressionUsesIdent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_ifBodyReturnsIndex tests detection of index return in if body.
func Test_ifBodyReturnsIndex(t *testing.T) {
	tests := []struct {
		name      string
		body      *ast.BlockStmt
		indexName string
		expected  bool
	}{
		{
			name:      "nil body",
			body:      nil,
			indexName: "i",
			expected:  false,
		},
		{
			name:      "empty body",
			body:      &ast.BlockStmt{List: []ast.Stmt{}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "no return statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "return with no results",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{}},
			}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "return with multiple results",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.Ident{Name: "i"},
					&ast.Ident{Name: "nil"},
				}},
			}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "return not identifier",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
				}},
			}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "return wrong identifier",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.Ident{Name: "other"},
				}},
			}},
			indexName: "i",
			expected:  false,
		},
		{
			name: "valid return index",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.Ident{Name: "i"},
				}},
			}},
			indexName: "i",
			expected:  true,
		},
		{
			name: "return index after other statements",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.Ident{Name: "i"},
				}},
			}},
			indexName: "i",
			expected:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := ifBodyReturnsIndex(tt.body, tt.indexName)
			// Verify result
			if result != tt.expected {
				t.Errorf("ifBodyReturnsIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_checkRangeForIndexPattern tests detection of index pattern in range.
func Test_checkRangeForIndexPattern(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		wantNil   bool
	}{
		{
			name: "invalid range variables",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			wantNil: true,
		},
		{
			name: "valid with pattern",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "v"},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.IfStmt{
						Cond: &ast.BinaryExpr{
							Op: token.EQL,
							X:  &ast.Ident{Name: "v"},
							Y:  &ast.Ident{Name: "target"},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "i"}}},
						}},
					},
				}},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := checkRangeForIndexPattern(tt.rangeStmt)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("checkRangeForIndexPattern() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_analyzeFunctionForIndexPattern tests function analysis.
func Test_analyzeFunctionForIndexPattern(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "function not returning int",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{Results: nil},
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
		},
		{
			name: "body too short",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{},
				}},
			},
		},
		{
			name: "no range statement",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
					&ast.ReturnStmt{},
				}},
			},
		},
		{
			name: "range not followed by return -1",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.RangeStmt{Body: &ast.BlockStmt{}},
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			defer func() {
				// Recover if panic happens
				if r := recover(); r != nil {
					t.Errorf("analyzeFunctionForIndexPattern() panicked: %v", r)
				}
			}()
			// Pass nil for pass and cfg - function should handle gracefully
			analyzeFunctionForIndexPattern(nil, tt.funcDecl, nil)
		})
	}
}
