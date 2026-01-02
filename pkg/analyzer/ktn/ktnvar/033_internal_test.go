package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_isValidIfStructureFor033 tests validation of if structure.
func Test_isValidIfStructureFor033(t *testing.T) {
	tests := []struct {
		name     string
		ifStmt   *ast.IfStmt
		expected bool
	}{
		{
			name: "has init",
			ifStmt: &ast.IfStmt{
				Init: &ast.AssignStmt{},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			expected: false,
		},
		{
			name: "has else",
			ifStmt: &ast.IfStmt{
				Else: &ast.BlockStmt{},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			expected: false,
		},
		{
			name: "nil body",
			ifStmt: &ast.IfStmt{
				Body: nil,
			},
			expected: false,
		},
		{
			name: "empty body",
			ifStmt: &ast.IfStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
			expected: false,
		},
		{
			name: "multiple statements in body",
			ifStmt: &ast.IfStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{},
					&ast.ExprStmt{},
				}},
			},
			expected: false,
		},
		{
			name: "valid structure",
			ifStmt: &ast.IfStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isValidIfStructureFor033(tt.ifStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("isValidIfStructureFor033() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractVarFromCondition033 tests variable extraction from condition.
func Test_extractVarFromCondition033(t *testing.T) {
	tests := []struct {
		name    string
		cond    ast.Expr
		wantNil bool
	}{
		{
			name:    "not binary expression",
			cond:    &ast.Ident{Name: "x"},
			wantNil: true,
		},
		{
			name: "not NEQ operator",
			cond: &ast.BinaryExpr{
				Op: token.EQL,
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
			},
			wantNil: true,
		},
		{
			name: "neither side is zero value",
			cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.Ident{Name: "y"},
			},
			wantNil: true,
		},
		{
			name: "valid pattern with zero on right",
			cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
			},
			wantNil: false,
		},
		{
			name: "valid pattern with zero on left",
			cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				Y:  &ast.Ident{Name: "x"},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractVarFromCondition033(tt.cond)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("extractVarFromCondition033() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_extractVarAndZeroFor033 tests extraction of variable and zero value.
func Test_extractVarAndZeroFor033(t *testing.T) {
	tests := []struct {
		name       string
		binaryExpr *ast.BinaryExpr
		wantVar    bool
		wantZero   bool
	}{
		{
			name: "zero on right",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.BasicLit{Kind: token.INT, Value: "0"},
			},
			wantVar:  true,
			wantZero: true,
		},
		{
			name: "zero on left",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.BasicLit{Kind: token.INT, Value: "0"},
				Y: &ast.Ident{Name: "x"},
			},
			wantVar:  true,
			wantZero: true,
		},
		{
			name: "neither is zero",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "y"},
			},
			wantVar:  false,
			wantZero: false,
		},
		{
			name: "empty string on right",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "s"},
				Y: &ast.BasicLit{Kind: token.STRING, Value: `""`},
			},
			wantVar:  true,
			wantZero: true,
		},
		{
			name: "backtick empty string on right",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "s"},
				Y: &ast.BasicLit{Kind: token.STRING, Value: "``"},
			},
			wantVar:  true,
			wantZero: true,
		},
		{
			name: "nil on right",
			binaryExpr: &ast.BinaryExpr{
				X: &ast.Ident{Name: "ptr"},
				Y: &ast.Ident{Name: "nil"},
			},
			wantVar:  true,
			wantZero: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			varExpr, zeroExpr := extractVarAndZeroFor033(tt.binaryExpr)
			gotVar := varExpr != nil
			gotZero := zeroExpr != nil
			// Verify variable result
			if gotVar != tt.wantVar {
				t.Errorf("extractVarAndZeroFor033() var = %v, want %v", gotVar, tt.wantVar)
			}
			// Verify zero result
			if gotZero != tt.wantZero {
				t.Errorf("extractVarAndZeroFor033() zero = %v, want %v", gotZero, tt.wantZero)
			}
		})
	}
}

// Test_isSimpleZeroValueFor033 tests detection of zero values.
func Test_isSimpleZeroValueFor033(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "integer zero",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "0"},
			expected: true,
		},
		{
			name:     "integer non-zero",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "1"},
			expected: false,
		},
		{
			name:     "empty string double quotes",
			expr:     &ast.BasicLit{Kind: token.STRING, Value: `""`},
			expected: true,
		},
		{
			name:     "empty string backticks",
			expr:     &ast.BasicLit{Kind: token.STRING, Value: "``"},
			expected: true,
		},
		{
			name:     "non-empty string",
			expr:     &ast.BasicLit{Kind: token.STRING, Value: `"hello"`},
			expected: false,
		},
		{
			name:     "nil identifier",
			expr:     &ast.Ident{Name: "nil"},
			expected: true,
		},
		{
			name:     "other identifier",
			expr:     &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name:     "float literal",
			expr:     &ast.BasicLit{Kind: token.FLOAT, Value: "0.0"},
			expected: false,
		},
		{
			name:     "call expression",
			expr:     &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isSimpleZeroValueFor033(tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isSimpleZeroValueFor033() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_validateIfBodyReturns033 tests validation of if body returns.
func Test_validateIfBodyReturns033(t *testing.T) {
	tests := []struct {
		name     string
		body     *ast.BlockStmt
		varExpr  ast.Expr
		expected bool
	}{
		{
			name: "not a return statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name: "valid return of variable",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "x"}}},
			}},
			varExpr:  &ast.Ident{Name: "x"},
			expected: true,
		},
		{
			name: "return different variable",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "y"}}},
			}},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := validateIfBodyReturns033(tt.body, tt.varExpr)
			// Verify result
			if result != tt.expected {
				t.Errorf("validateIfBodyReturns033() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_returnsVariableFor033 tests if return statement returns the variable.
func Test_returnsVariableFor033(t *testing.T) {
	tests := []struct {
		name       string
		returnStmt *ast.ReturnStmt
		varExpr    ast.Expr
		expected   bool
	}{
		{
			name: "no results",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{},
			},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name: "multiple results",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{&ast.Ident{Name: "x"}, &ast.Ident{Name: "y"}},
			},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name: "varExpr not identifier",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			varExpr:  &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			expected: false,
		},
		{
			name: "result not identifier",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "f"}}},
			},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name: "names match",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{&ast.Ident{Name: "x"}},
			},
			varExpr:  &ast.Ident{Name: "x"},
			expected: true,
		},
		{
			name: "names differ",
			returnStmt: &ast.ReturnStmt{
				Results: []ast.Expr{&ast.Ident{Name: "y"}},
			},
			varExpr:  &ast.Ident{Name: "x"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := returnsVariableFor033(tt.returnStmt, tt.varExpr)
			// Verify result
			if result != tt.expected {
				t.Errorf("returnsVariableFor033() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_matchesCmpOrPattern tests detection of cmp.Or pattern.
func Test_matchesCmpOrPattern(t *testing.T) {
	tests := []struct {
		name       string
		ifStmt     *ast.IfStmt
		returnStmt *ast.ReturnStmt
		expected   bool
	}{
		{
			name: "invalid if structure (has init)",
			ifStmt: &ast.IfStmt{
				Init: &ast.AssignStmt{},
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "x"}}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "def"}}},
			expected:   false,
		},
		{
			name: "invalid condition (not NEQ)",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.EQL,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "x"}}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "def"}}},
			expected:   false,
		},
		{
			name: "body does not return variable",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "y"}}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "def"}}},
			expected:   false,
		},
		{
			name: "empty return results",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "x"}}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{}},
			expected:   false,
		},
		{
			name: "valid pattern",
			ifStmt: &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  &ast.Ident{Name: "x"},
					Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "x"}}},
				}},
			},
			returnStmt: &ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "def"}}},
			expected:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := matchesCmpOrPattern(tt.ifStmt, tt.returnStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("matchesCmpOrPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_analyzeBodyForCmpOrPattern tests body analysis for cmp.Or pattern.
func Test_analyzeBodyForCmpOrPattern(t *testing.T) {
	// This test verifies the function doesn't panic with various inputs
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
			name: "non-if first statement",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				&ast.ReturnStmt{},
			}},
		},
		{
			name: "if not followed by return",
			body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.IfStmt{Body: &ast.BlockStmt{}},
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic - we can't verify reports without pass
			defer func() {
				// Recover if panic happens
				if r := recover(); r != nil {
					t.Errorf("analyzeBodyForCmpOrPattern() panicked: %v", r)
				}
			}()
			// Pass nil for pass and cfg - function should handle gracefully
			analyzeBodyForCmpOrPattern(nil, tt.body, nil)
		})
	}
}

// Test_runVar033_ruleDisabled tests that analyzer returns early when rule is disabled.
func Test_runVar033_ruleDisabled(t *testing.T) {
	// Save the current config
	oldCfg := config.Get()

	// Create new config with rule disabled
	newCfg := config.DefaultConfig()
	falseVal := false
	newCfg.Rules[ruleCodeVar033] = &config.RuleConfig{Enabled: &falseVal}
	config.Set(newCfg)
	// Ensure restoration at the end
	defer config.Set(oldCfg)

	// Run analyzer with testhelper
	diags := testhelper.RunAnalyzer(t, Analyzer033, "testdata/src/var033/bad.go")

	// With rule disabled, should have 0 errors
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when rule disabled, got %d", len(diags))
	}
}

// Test_checkCmpOrPattern_funcLit tests function literals are analyzed.
func Test_checkCmpOrPattern_funcLit(t *testing.T) {
	// Test file with function literal containing cmp.Or pattern
	src := `package test

const DefaultVal = 10

func wrapper() func(int) int {
	return func(x int) int {
		if x != 0 {
			return x
		}
		return DefaultVal
	}
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a mock pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer to populate ResultOf
	inspResult, err := inspect.Analyzer.Run(pass)
	// Check for inspector errors
	if err != nil {
		t.Fatalf("Failed to run inspect: %v", err)
	}
	pass.ResultOf[inspect.Analyzer] = inspResult

	// Get config
	cfg := config.Get()

	// Get the inspector
	insp := inspResult.(*inspector.Inspector)

	// Run the checkCmpOrPattern function
	checkCmpOrPattern(pass, insp, cfg)

	// Should detect the pattern in function literal
	if len(diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic for func literal, got %d", len(diagnostics))
	}
}

// Test_checkCmpOrPattern_funcDeclNilBody tests external function declarations.
func Test_checkCmpOrPattern_funcDeclNilBody(t *testing.T) {
	// Test file with external function declaration (no body)
	src := `package test

// External function declaration (no body)
func ExternalFunc(x int) int
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a mock pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer to populate ResultOf
	inspResult, err := inspect.Analyzer.Run(pass)
	// Check for inspector errors
	if err != nil {
		t.Fatalf("Failed to run inspect: %v", err)
	}
	pass.ResultOf[inspect.Analyzer] = inspResult

	// Get config
	cfg := config.Get()

	// Get the inspector
	insp := inspResult.(*inspector.Inspector)

	// Run the checkCmpOrPattern function - should not panic on nil body
	checkCmpOrPattern(pass, insp, cfg)

	// Should not detect any issues (no body to analyze)
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for nil body, got %d", len(diagnostics))
	}
}

// Test_checkCmpOrPattern_fileExcluded tests file exclusion.
func Test_checkCmpOrPattern_fileExcluded(t *testing.T) {
	// Test source
	src := `package test

const DefaultVal = 10

func getValue(x int) int {
	if x != 0 {
		return x
	}
	return DefaultVal
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded_test_file.go", src, parser.ParseComments)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a mock pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer to populate ResultOf
	inspResult, err := inspect.Analyzer.Run(pass)
	// Check for inspector errors
	if err != nil {
		t.Fatalf("Failed to run inspect: %v", err)
	}
	pass.ResultOf[inspect.Analyzer] = inspResult

	// Create config with file exclusion for this rule
	cfg := config.DefaultConfig()
	cfg.Rules[ruleCodeVar033] = &config.RuleConfig{
		Exclude: []string{"excluded_test_file.go"},
	}

	// Get the inspector
	insp := inspResult.(*inspector.Inspector)

	// Run the checkCmpOrPattern function
	checkCmpOrPattern(pass, insp, cfg)

	// Should not detect any issues (file excluded)
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for excluded file, got %d", len(diagnostics))
	}
}
