package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
)

// Test_isCallWithArg1 verifie la detection d'un appel avec argument 1.
func Test_isCallWithArg1(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "call with arg 1",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "1"},
				},
			},
			expected: true,
		},
		{
			name: "call with arg 2",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "2"},
				},
			},
			expected: false,
		},
		{
			name: "call with no args",
			call: &ast.CallExpr{
				Args: []ast.Expr{},
			},
			expected: false,
		},
		{
			name: "call with multiple args",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "1"},
					&ast.BasicLit{Kind: token.INT, Value: "2"},
				},
			},
			expected: false,
		},
		{
			name: "call with string arg",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.STRING, Value: "\"1\""},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isCallWithArg1(tt.call)
			// Verification du resultat
			if result != tt.expected {
				t.Errorf("isCallWithArg1() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_getReceiverName verifie l'extraction du nom du receiver.
func Test_getReceiverName(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple ident",
			expr:     &ast.Ident{Name: "wg"},
			expected: "wg",
		},
		{
			name:     "non-ident expr",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "1"},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := getReceiverName(tt.expr)
			// Verification du resultat
			if result != tt.expected {
				t.Errorf("getReceiverName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractFuncLiteral verifie l'extraction de fonction literal.
func Test_extractFuncLiteral(t *testing.T) {
	tests := []struct {
		name     string
		goStmt   *ast.GoStmt
		expected bool
	}{
		{
			name: "go with func literal",
			goStmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Body: &ast.BlockStmt{},
					},
				},
			},
			expected: true,
		},
		{
			name: "go with ident (non-literal)",
			goStmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.Ident{Name: "myFunc"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := extractFuncLiteral(tt.goStmt)
			// Verification du resultat
			if (result != nil) != tt.expected {
				t.Errorf("extractFuncLiteral() returned %v, want non-nil=%v", result, tt.expected)
			}
		})
	}
}
