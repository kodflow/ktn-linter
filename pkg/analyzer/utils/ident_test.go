package utils

import (
	"go/ast"
	"go/parser"
	"testing"

)

func TestIsIdentCall(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		funcName string
		want     bool
	}{
		{
			name:     "make call",
			code:     `make([]int, 0)`,
			funcName: "make",
			want:     true,
		},
		{
			name:     "append call",
			code:     `append(slice, 1)`,
			funcName: "append",
			want:     true,
		},
		{
			name:     "different function",
			code:     `fmt.Println()`,
			funcName: "make",
			want:     false,
		},
		{
			name:     "method call",
			code:     `obj.Method()`,
			funcName: "Method",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			call, ok := expr.(*ast.CallExpr)
			if !ok {
				t.Fatalf("Expression is not a CallExpr")
			}

			got := IsIdentCall(call, tt.funcName)
			if got != tt.want {
				t.Errorf("IsIdentCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBuiltinCall(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		builtinName string
		want        bool
	}{
		{
			name:        "make builtin correct",
			code:        `make([]int, 0)`,
			builtinName: "make",
			want:        true,
		},
		{
			name:        "make builtin wrong name",
			code:        `make([]int, 0)`,
			builtinName: "append",
			want:        false,
		},
		{
			name:        "append builtin correct",
			code:        `append(slice, 1)`,
			builtinName: "append",
			want:        true,
		},
		{
			name:        "len builtin correct",
			code:        `len(slice)`,
			builtinName: "len",
			want:        true,
		},
		{
			name:        "non-builtin function",
			code:        `customFunc()`,
			builtinName: "make",
			want:        false,
		},
		{
			name:        "method call",
			code:        `obj.Method()`,
			builtinName: "Method",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			call, ok := expr.(*ast.CallExpr)
			if !ok {
				t.Fatalf("Expression is not a CallExpr")
			}

			got := IsBuiltinCall(call, tt.builtinName)
			if got != tt.want {
				t.Errorf("IsBuiltinCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIdentName(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple identifier",
			code: `myVar`,
			want: "myVar",
		},
		{
			name: "selector expression",
			code: `pkg.Func`,
			want: "",
		},
		{
			name: "literal",
			code: `123`,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := GetIdentName(expr)
			if got != tt.want {
				t.Errorf("GetIdentName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractVarName(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple identifier",
			code: `myVar`,
			want: "myVar",
		},
		{
			name: "index expression",
			code: `arr[i]`,
			want: "arr[...]",
		},
		{
			name: "selector expression",
			code: `obj.Field`,
			want: "obj.Field",
		},
		{
			name: "complex selector",
			code: `pkg.Obj.Field`,
			want: "pkg.Obj.Field",
		},
		{
			name: "literal",
			code: `"string"`,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := ExtractVarName(expr)
			if got != tt.want {
				t.Errorf("ExtractVarName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtractVarNameExtended tests additional cases for ExtractVarName
func TestExtractVarNameExtended(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "star expression",
			code: `*ptr`,
			want: "*ptr",
		},
		{
			name: "parenthesized expression",
			code: `(x)`,
			want: "x",
		},
		{
			name: "nested parentheses",
			code: `((y))`,
			want: "y",
		},
		{
			name: "function call (unsupported)",
			code: `fn()`,
			want: "",
		},
		{
			name: "star of complex",
			code: `*arr[0]`,
			want: "*arr[...]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := ExtractVarName(expr)
			if got != tt.want {
				t.Errorf("ExtractVarName() = %v, want %v", got, tt.want)
			}
		})
	}
}
