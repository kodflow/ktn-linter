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
		name      string
		body      *ast.BlockStmt
		vars      *rangeVariables036
		wantFound bool
	}{
		{
			name:      "nil body",
			body:      nil,
			vars:      vars,
			wantFound: false,
		},
		{
			name:      "empty body",
			body:      &ast.BlockStmt{List: []ast.Stmt{}},
			vars:      vars,
			wantFound: false,
		},
		{
			name: "no matching statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
			vars:      vars,
			wantFound: false,
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
			vars:      vars,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, found := findIndexPatternInBody036(tt.body, tt.vars)
			// Verify result
			if found != tt.wantFound {
				t.Errorf("findIndexPatternInBody036() found = %v, want %v", found, tt.wantFound)
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
			name: "valid pattern with value on left",
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
		{
			name: "valid pattern with value on right",
			stmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "target"},
					Y:  &ast.Ident{Name: "v"},
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
		wantFound bool
	}{
		{
			name: "invalid range variables",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			wantFound: false,
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
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, found := checkRangeForIndexPattern(tt.rangeStmt)
			// Verify result
			if found != tt.wantFound {
				t.Errorf("checkRangeForIndexPattern() found = %v, want %v", found, tt.wantFound)
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
		{
			name: "range at last position",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
					&ast.RangeStmt{Body: &ast.BlockStmt{}},
				}},
			},
		},
		{
			name: "range with valid pattern and return -1",
			funcDecl: &ast.FuncDecl{
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.RangeStmt{
						Key:   &ast.Ident{Name: "i"},
						Value: &ast.Ident{Name: "v"},
						Body:  &ast.BlockStmt{List: []ast.Stmt{}},
					},
					&ast.ReturnStmt{Results: []ast.Expr{
						&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.INT, Value: "1"}},
					}},
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

// Test_runVar036_ruleDisabled tests runVar036 when rule is disabled.
func Test_runVar036_ruleDisabled(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar036]

	// Disable the rule
	cfg.Rules[ruleCodeVar036] = &config.RuleConfig{Enabled: config.Bool(false)}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar036)
		} else {
			cfg.Rules[ruleCodeVar036] = originalRule
		}
	}()

	// Parse code with index search pattern
	src := `package test

func findIndex(s []int, target int) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer036,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar036(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar036() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when disabled
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics when rule disabled, got %d", len(diagnostics))
	}
}

// Test_runVar036_fileExcluded tests runVar036 when file is excluded.
func Test_runVar036_fileExcluded(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar036]

	// Set up rule with file exclusion
	cfg.Rules[ruleCodeVar036] = &config.RuleConfig{
		Exclude: []string{"excluded.go"},
	}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar036)
		} else {
			cfg.Rules[ruleCodeVar036] = originalRule
		}
	}()

	// Parse code with index search pattern
	src := `package test

func findIndex(s []int, target int) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer036,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar036(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar036() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when file excluded
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics when file excluded, got %d", len(diagnostics))
	}
}

// Test_runVar036_nilBody tests runVar036 with function having nil body.
func Test_runVar036_nilBody(t *testing.T) {
	// Parse code with external function declaration (nil body)
	src := `package test

// External function declaration
func externalFunc(s []int, target int) int
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer036,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar036(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar036() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports for nil body
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics for nil body, got %d", len(diagnostics))
	}
}
