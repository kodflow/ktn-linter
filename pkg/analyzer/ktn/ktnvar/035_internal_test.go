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

// Test_isBlankOrNil tests detection of blank identifier or nil.
func Test_isBlankOrNil(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "nil expression",
			expr:     nil,
			expected: true,
		},
		{
			name:     "blank identifier",
			expr:     &ast.Ident{Name: "_"},
			expected: true,
		},
		{
			name:     "normal identifier",
			expr:     &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name:     "non-identifier",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "1"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isBlankOrNil(tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isBlankOrNil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_matchesContainsPattern tests detection of contains pattern.
func Test_matchesContainsPattern(t *testing.T) {
	tests := []struct {
		name       string
		rangeStmt  *ast.RangeStmt
		returnStmt *ast.ReturnStmt
		expected   bool
	}{
		{
			name: "key is not blank",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "v"},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "no value variable",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: nil,
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "nil body",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body:  nil,
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "empty body",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body:  &ast.BlockStmt{List: []ast.Stmt{}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "multiple statements in body",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.IfStmt{},
					&ast.ExprStmt{},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "body not if statement",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "if does not match pattern",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.IfStmt{
						Else: &ast.BlockStmt{},
						Cond: &ast.Ident{Name: "cond"},
					},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
			expected:   false,
		},
		{
			name: "return not false",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.IfStmt{
						Cond: &ast.BinaryExpr{
							Op: token.EQL,
							X:  &ast.Ident{Name: "v"},
							Y:  &ast.Ident{Name: "target"},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "true"}}},
						}},
					},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "true"}}},
			expected:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := matchesContainsPattern(tt.rangeStmt, tt.returnStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("matchesContainsPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_matchesIfReturnTruePattern tests detection of if return true pattern.
func Test_matchesIfReturnTruePattern(t *testing.T) {
	tests := []struct {
		name       string
		ifStmt     *ast.IfStmt
		rangeValue ast.Expr
		expected   bool
	}{
		{
			name: "has else",
			ifStmt: &ast.IfStmt{
				Else: &ast.BlockStmt{},
				Cond: &ast.BinaryExpr{Op: token.EQL},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "cond not binary",
			ifStmt: &ast.IfStmt{
				Cond: &ast.Ident{Name: "x"},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "not equality operator",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "range value not in comparison",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.Ident{Name: "y"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "nil body",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: nil,
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "empty body",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "multiple statements in body",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{},
					&ast.ExprStmt{},
				}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "body not return",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "returns false",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "false"}}},
				}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "valid pattern with range value on left",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "v"},
					Y:  &ast.Ident{Name: "target"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "true"}}},
				}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
		{
			name: "valid pattern with range value on right",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "target"},
					Y:  &ast.Ident{Name: "v"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "true"}}},
				}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := matchesIfReturnTruePattern(tt.ifStmt, tt.rangeValue)
			// Verify result
			if result != tt.expected {
				t.Errorf("matchesIfReturnTruePattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_containsRangeValue tests detection of range value in comparison.
func Test_containsRangeValue(t *testing.T) {
	tests := []struct {
		name       string
		binaryExpr *ast.BinaryExpr
		rangeValue ast.Expr
		expected   bool
	}{
		{
			name: "range value not identifier",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "y"},
			},
			rangeValue: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			expected:   false,
		},
		{
			name: "found on left",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "v"},
				Y: &ast.Ident{Name: "target"},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
		{
			name: "found on right",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "target"},
				Y: &ast.Ident{Name: "v"},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
		{
			name: "not found",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "y"},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
		{
			name: "left not identifier",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
				Y: &ast.Ident{Name: "v"},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
		{
			name: "right not identifier",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "v"},
				Y: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   true,
		},
		{
			name: "neither is identifier",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
				Y: &ast.CallExpr{Fun: &ast.Ident{Name: "g"}},
			},
			rangeValue: &ast.Ident{Name: "v"},
			expected:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := containsRangeValue(tt.binaryExpr, tt.rangeValue)
			// Verify result
			if result != tt.expected {
				t.Errorf("containsRangeValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_returnsLiteralBool tests detection of boolean literals in return.
func Test_returnsLiteralBool(t *testing.T) {
	tests := []struct {
		name       string
		returnStmt *ast.ReturnStmt
		expected   bool
		wantTrue   bool
	}{
		{
			name:       "no results",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{}},
			expected:   false,
			wantTrue:   true,
		},
		{
			name: "multiple results",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.Ident{Name: "true"},
				&ast.Ident{Name: "nil"},
			}},
			expected: false,
			wantTrue: true,
		},
		{
			name: "not identifier",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			}},
			expected: false,
			wantTrue: true,
		},
		{
			name: "returns true - expect true",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.Ident{Name: "true"},
			}},
			expected: true,
			wantTrue: true,
		},
		{
			name: "returns true - expect false",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.Ident{Name: "true"},
			}},
			expected: false,
			wantTrue: false,
		},
		{
			name: "returns false - expect false",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.Ident{Name: "false"},
			}},
			expected: true,
			wantTrue: false,
		},
		{
			name: "returns false - expect true",
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{
				&ast.Ident{Name: "false"},
			}},
			expected: false,
			wantTrue: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := returnsLiteralBool(tt.returnStmt, tt.wantTrue)
			// Verify result
			if result != tt.expected {
				t.Errorf("returnsLiteralBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_analyzeBodyForContainsPattern tests body analysis for contains pattern.
func Test_analyzeBodyForContainsPattern(t *testing.T) {
	tests := []struct {
		name string
		body *ast.BlockStmt
	}{
		{
			name: "empty body",
			body: &ast.BlockStmt{List: []ast.Stmt{}},
		},
		{
			name: "single statement (not enough)",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{},
			}},
		},
		{
			name: "non-range first statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				&ast.ReturnStmt{},
			}},
		},
		{
			name: "range not followed by return",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.RangeStmt{Body: &ast.BlockStmt{}},
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			defer func() {
				// Recover if panic happens
				if r := recover(); r != nil {
					t.Errorf("analyzeBodyForContainsPattern() panicked: %v", r)
				}
			}()
			// Pass nil for pass and cfg - function should handle gracefully
			analyzeBodyForContainsPattern(nil, tt.body, nil)
		})
	}
}

// Test_runVar035_ruleDisabled tests runVar035 when rule is disabled.
func Test_runVar035_ruleDisabled(t *testing.T) {
	// Save the current config
	oldCfg := config.Get()

	// Create new config with rule disabled
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar035] = &config.RuleConfig{Enabled: config.Bool(false)}
	config.Set(newCfg)
	// Ensure restoration at the end
	defer config.Set(oldCfg)

	// Create minimal pass
	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Analyzer: Analyzer035,
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: &inspector.Inspector{},
		},
	}

	// Run should return early due to disabled rule
	result, err := runVar035(pass)

	// Verify no error
	if err != nil {
		t.Errorf("runVar035() error = %v, want nil", err)
	}

	// Verify nil result
	if result != nil {
		t.Errorf("runVar035() result = %v, want nil", result)
	}
}

// Test_checkContainsPattern_fileExcluded tests file exclusion logic directly.
func Test_checkContainsPattern_fileExcluded(t *testing.T) {
	// Save and restore config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Configure rule with file exclusion
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar035] = &config.RuleConfig{
		Enabled: config.Bool(true),
		Exclude: []string{"excluded.go"},
	}
	config.Set(newCfg)

	// Parse a valid source file to create proper AST
	src := `package test

func testFunc() {
	slice := []int{1, 2, 3}
	for i := 0; i < len(slice); i++ {
		if slice[i] == 2 {
			return
		}
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create inspector
	insp := inspector.New([]*ast.File{file})

	// Create pass with the file
	pass := &analysis.Pass{
		Analyzer: Analyzer035,
		Fset:     fset,
	}

	cfg := config.Get()

	// Call checkContainsPattern directly - this should trigger file exclusion branch
	checkContainsPattern(pass, insp, cfg)

	// If we reach here without panic, the file exclusion branch was exercised
}

// Test_checkContainsPattern_funcLit tests checkContainsPattern with FuncLit.
func Test_checkContainsPattern_funcLit(t *testing.T) {
	// Parse source with function literal
	src := `package test

func wrapper() {
	fn := func() {
		slice := []int{1, 2, 3}
		for i := 0; i < len(slice); i++ {
			if slice[i] == 2 {
				return
			}
		}
	}
	_ = fn
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create inspector
	insp := inspector.New([]*ast.File{file})

	// Create pass
	pass := &analysis.Pass{
		Analyzer: Analyzer035,
		Fset:     fset,
	}

	// Create config
	cfg := config.Get()

	// Run checkContainsPattern - should handle FuncLit branch
	checkContainsPattern(pass, insp, cfg)

	// No assertion needed - just verifying no panic and FuncLit branch is covered
}

// Test_checkContainsPattern_funcDeclNilBody tests checkContainsPattern with FuncDecl nil body.
func Test_checkContainsPattern_funcDeclNilBody(t *testing.T) {
	// Parse source with external function declaration (no body)
	src := `package test

// External function declaration (no body)
func ExternalFunc(x int) int
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create inspector
	insp := inspector.New([]*ast.File{file})

	// Create pass
	pass := &analysis.Pass{
		Analyzer: Analyzer035,
		Fset:     fset,
	}

	// Create config
	cfg := config.Get()

	// Run checkContainsPattern - should handle nil body branch
	checkContainsPattern(pass, insp, cfg)

	// No assertion needed - just verifying no panic and nil body branch is covered
}
