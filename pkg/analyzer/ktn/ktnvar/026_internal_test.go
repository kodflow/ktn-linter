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
)

// TestHasReturnInElse tests the hasReturnInElse function.
func TestHasReturnInElse(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		elseStmt ast.Stmt
		expected bool
	}{
		{
			name:     "nil block stmt",
			elseStmt: (*ast.BlockStmt)(nil),
			expected: false,
		},
		{
			name:     "empty block stmt",
			elseStmt: &ast.BlockStmt{List: nil},
			expected: false,
		},
		{
			name:     "block with empty list",
			elseStmt: &ast.BlockStmt{List: []ast.Stmt{}},
			expected: false,
		},
		{
			name: "block with return stmt",
			elseStmt: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{},
				},
			},
			expected: true,
		},
		{
			name: "block with non-return stmt",
			elseStmt: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				},
			},
			expected: false,
		},
		{
			name:     "non-block stmt",
			elseStmt: &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := hasReturnInElse(tt.elseStmt)
			// Check result
			if result != tt.expected {
				t.Errorf("hasReturnInElse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckMathMinMax tests checkMathMinMax function.
func TestCheckMathMinMax(t *testing.T) {
	tests := []struct {
		name         string
		callExpr     *ast.CallExpr
		expectReport bool
	}{
		{
			name: "not a selector",
			callExpr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "Min"},
			},
			expectReport: false,
		},
		{
			name: "X is not an identifier",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.CallExpr{Fun: &ast.Ident{Name: "getMath"}},
					Sel: &ast.Ident{Name: "Min"},
				},
			},
			expectReport: false,
		},
		{
			name: "not the math package",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "other"},
					Sel: &ast.Ident{Name: "Min"},
				},
			},
			expectReport: false,
		},
		{
			name: "not Min or Max",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "math"},
					Sel: &ast.Ident{Name: "Abs"},
				},
			},
			expectReport: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reported := false
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reported = true
				},
			}
			// Call function
			checkMathMinMax(pass, tt.callExpr)
			// Check result
			if reported != tt.expectReport {
				t.Errorf("checkMathMinMax() reported = %v, want %v", reported, tt.expectReport)
			}
		})
	}
}

// TestGetBuiltinName tests getBuiltinName function.
func TestGetBuiltinName(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Min to min",
			input:    "Min",
			expected: "min",
		},
		{
			name:     "Max to max",
			input:    "Max",
			expected: "max",
		},
		{
			name:     "other returns max",
			input:    "Other",
			expected: "max",
		},
	}
	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := getBuiltinName(tt.input)
			// Check result
			if result != tt.expected {
				t.Errorf("getBuiltinName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckIfMinMaxPattern tests checkIfMinMaxPattern function.
func TestCheckIfMinMaxPattern(t *testing.T) {
	// Test: condition is not a min/max condition
	ifNotMinMax := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.EQL,
			X:  &ast.Ident{Name: "a"},
			Y:  &ast.Ident{Name: "b"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{&ast.ReturnStmt{}},
		},
	}
	checkIfMinMaxPattern(nil, ifNotMinMax)

	// Test: body doesn't have return
	ifNoReturn := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.LSS,
			X:  &ast.Ident{Name: "a"},
			Y:  &ast.Ident{Name: "b"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{&ast.ExprStmt{X: &ast.Ident{Name: "x"}}},
		},
	}
	checkIfMinMaxPattern(nil, ifNoReturn)

	// Test: no matching return (no else)
	ifNoElse := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.LSS,
			X:  &ast.Ident{Name: "a"},
			Y:  &ast.Ident{Name: "b"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{&ast.ReturnStmt{}},
		},
		Else: nil,
	}
	checkIfMinMaxPattern(nil, ifNoElse)
}

// TestIsMinMaxCondition tests isMinMaxCondition function.
func TestIsMinMaxCondition(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		cond     ast.Expr
		expected bool
	}{
		{
			name: "less than",
			cond: &ast.BinaryExpr{
				Op: token.LSS,
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
			},
			expected: true,
		},
		{
			name: "greater than",
			cond: &ast.BinaryExpr{
				Op: token.GTR,
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
			},
			expected: true,
		},
		{
			name: "less or equal",
			cond: &ast.BinaryExpr{
				Op: token.LEQ,
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
			},
			expected: false,
		},
		{
			name: "greater or equal",
			cond: &ast.BinaryExpr{
				Op: token.GEQ,
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
			},
			expected: false,
		},
		{
			name: "equal",
			cond: &ast.BinaryExpr{
				Op: token.EQL,
				X:  &ast.Ident{Name: "a"},
				Y:  &ast.Ident{Name: "b"},
			},
			expected: false,
		},
		{
			name:     "not binary",
			cond:     &ast.Ident{Name: "x"},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isMinMaxCondition(tt.cond)
			// Check result
			if result != tt.expected {
				t.Errorf("isMinMaxCondition() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestHasReturnInBody tests hasReturnInBody function.
func TestHasReturnInBody(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		body     *ast.BlockStmt
		expected bool
	}{
		{
			name:     "nil body",
			body:     nil,
			expected: false,
		},
		{
			name:     "empty body",
			body:     &ast.BlockStmt{List: nil},
			expected: false,
		},
		{
			name:     "empty list",
			body:     &ast.BlockStmt{List: []ast.Stmt{}},
			expected: false,
		},
		{
			name: "first stmt is return",
			body: &ast.BlockStmt{
				List: []ast.Stmt{&ast.ReturnStmt{}},
			},
			expected: true,
		},
		{
			name: "first stmt is not return",
			body: &ast.BlockStmt{
				List: []ast.Stmt{&ast.ExprStmt{X: &ast.Ident{Name: "x"}}},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := hasReturnInBody(tt.body)
			// Check result
			if result != tt.expected {
				t.Errorf("hasReturnInBody() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestHasMatchingReturn tests hasMatchingReturn function.
func TestHasMatchingReturn(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		ifStmt   *ast.IfStmt
		expected bool
	}{
		{
			name: "no else",
			ifStmt: &ast.IfStmt{
				Cond: &ast.Ident{Name: "x"},
				Body: &ast.BlockStmt{},
				Else: nil,
			},
			expected: false,
		},
		{
			name: "else with return",
			ifStmt: &ast.IfStmt{
				Cond: &ast.Ident{Name: "x"},
				Body: &ast.BlockStmt{},
				Else: &ast.BlockStmt{
					List: []ast.Stmt{&ast.ReturnStmt{}},
				},
			},
			expected: true,
		},
		{
			name: "else without return",
			ifStmt: &ast.IfStmt{
				Cond: &ast.Ident{Name: "x"},
				Body: &ast.BlockStmt{},
				Else: &ast.BlockStmt{
					List: []ast.Stmt{&ast.ExprStmt{X: &ast.Ident{Name: "y"}}},
				},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := hasMatchingReturn(tt.ifStmt)
			// Check result
			if result != tt.expected {
				t.Errorf("hasMatchingReturn() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_runVar026_ruleDisabled tests analyzer returns early when disabled.
func Test_runVar026_ruleDisabled(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar026]

	// Disable the rule
	cfg.Rules[ruleCodeVar026] = &config.RuleConfig{Enabled: config.Bool(false)}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar026)
		} else {
			cfg.Rules[ruleCodeVar026] = originalRule
		}
	}()

	// Run analyzer with testhelper
	diags := testhelper.RunAnalyzer(t, Analyzer026, "testdata/src/var026/bad.go")

	// With rule disabled, should have 0 errors
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when rule disabled, got %d", len(diags))
	}
}

// Test_runVar026_fileExcluded tests file exclusion.
func Test_runVar026_fileExcluded(t *testing.T) {
	// Test source with math.Min
	src := `package test

import "math"

func getMin(a, b float64) float64 {
	return math.Min(a, b)
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded_file.go", src, parser.ParseComments)
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

	// Get config and add file exclusion
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar026]

	// Add the exclusion
	cfg.Rules[ruleCodeVar026] = &config.RuleConfig{Exclude: []string{"excluded_file.go"}}
	// Ensure cleanup
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar026)
		} else {
			cfg.Rules[ruleCodeVar026] = originalRule
		}
	}()

	// Run the actual analyzer
	_, runErr := runVar026(pass)
	// Check for analyzer errors
	if runErr != nil {
		t.Fatalf("runVar026() returned error: %v", runErr)
	}

	// Should not detect any issues (file excluded)
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for excluded file, got %d", len(diagnostics))
	}
}

// Test_checkIfMinMaxPattern_report tests checkIfMinMaxPattern reporting.
func Test_checkIfMinMaxPattern_report(t *testing.T) {
	// Test source with if/else min/max pattern
	src := `package test

func getMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
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

	// Run the analyzer
	_, runErr := runVar026(pass)
	// Check for analyzer errors
	if runErr != nil {
		t.Fatalf("runVar026() returned error: %v", runErr)
	}

	// Should detect the if/else min/max pattern
	if len(diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diagnostics))
	}
}
