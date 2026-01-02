package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
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
	// Test: not a selector
	callNotSelector := &ast.CallExpr{
		Fun: &ast.Ident{Name: "Min"},
	}
	checkMathMinMax(nil, callNotSelector)

	// Test: X is not an identifier
	callNonIdentX := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.CallExpr{Fun: &ast.Ident{Name: "getMath"}},
			Sel: &ast.Ident{Name: "Min"},
		},
	}
	checkMathMinMax(nil, callNonIdentX)

	// Test: not the math package
	callNotMath := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "other"},
			Sel: &ast.Ident{Name: "Min"},
		},
	}
	checkMathMinMax(nil, callNotMath)

	// Test: not Min or Max
	callNotMinMax := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "math"},
			Sel: &ast.Ident{Name: "Abs"},
		},
	}
	checkMathMinMax(nil, callNotMinMax)
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
