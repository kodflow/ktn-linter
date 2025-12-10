// Internal tests for ident.go - testing private behavior.
package utils

import (
	"go/ast"
	"go/parser"
	"testing"
)

// Test_extractVarName_IndexExpr tests ExtractVarName with index expressions.
//
// Params:
//   - t: testing context
func Test_extractVarName_IndexExpr(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple index expression",
			code: "items[i]",
			want: "items[...]",
		},
		{
			name: "nested index expression",
			code: "matrix[i][j]",
			want: "matrix[...][...]",
		},
		{
			name: "invalid index expression",
			code: "123[i]",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.code)
			got := ExtractVarName(expr)
			// Check result
			if got != tt.want {
				t.Errorf("ExtractVarName(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// Test_extractVarName_SelectorExpr tests ExtractVarName with selector expressions.
//
// Params:
//   - t: testing context
func Test_extractVarName_SelectorExpr(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple selector",
			code: "obj.field",
			want: "obj.field",
		},
		{
			name: "nested selector",
			code: "obj.nested.field",
			want: "obj.nested.field",
		},
		{
			name: "nested selector with method call",
			code: "obj.Method().field",
			want: "field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.code)
			got := ExtractVarName(expr)
			// Check result
			if got != tt.want {
				t.Errorf("ExtractVarName(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// Test_extractVarName_StarExpr tests ExtractVarName with pointer dereference.
//
// Params:
//   - t: testing context
func Test_extractVarName_StarExpr(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple dereference",
			code: "*ptr",
			want: "*ptr",
		},
		{
			name: "double dereference",
			code: "**ptr",
			want: "**ptr",
		},
		{
			name: "invalid dereference",
			code: "*123",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.code)
			got := ExtractVarName(expr)
			// Check result
			if got != tt.want {
				t.Errorf("ExtractVarName(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// Test_extractVarName_ParenExpr tests ExtractVarName with parentheses.
//
// Params:
//   - t: testing context
func Test_extractVarName_ParenExpr(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "parenthesized identifier",
			code: "(x)",
			want: "x",
		},
		{
			name: "nested parentheses",
			code: "((x))",
			want: "x",
		},
		{
			name: "parenthesized selector",
			code: "(obj.field)",
			want: "obj.field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.code)
			got := ExtractVarName(expr)
			// Check result
			if got != tt.want {
				t.Errorf("ExtractVarName(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// Test_extractVarName_UnknownExpr tests ExtractVarName with unknown expressions.
//
// Params:
//   - t: testing context
func Test_extractVarName_UnknownExpr(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "literal number",
			code: "42",
			want: "",
		},
		{
			name: "literal string",
			code: `"hello"`,
			want: "",
		},
		{
			name: "function call",
			code: "foo()",
			want: "",
		},
		{
			name: "binary expression",
			code: "a + b",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.code)
			got := ExtractVarName(expr)
			// Check result
			if got != tt.want {
				t.Errorf("ExtractVarName(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// parseExpr parses a Go expression for testing.
//
// Params:
//   - t: testing context
//   - code: expression code
//
// Returns:
//   - ast.Expr: parsed expression
func parseExpr(t *testing.T, code string) ast.Expr {
	t.Helper()
	expr, err := parser.ParseExpr(code)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse expression %q: %v", code, err)
	}
	// Return expression
	return expr
}
