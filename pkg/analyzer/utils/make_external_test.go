package utils_test

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
)

func TestIsMakeCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "make slice",
			code: `make([]int, 0)`,
			want: true,
		},
		{
			name: "make map",
			code: `make(map[string]int)`,
			want: true,
		},
		{
			name: "make channel",
			code: `make(chan int, 10)`,
			want: true,
		},
		{
			name: "not make call",
			code: `append(slice, 1)`,
			want: false,
		},
		{
			name: "method call",
			code: `obj.make()`,
			want: false,
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

			got := utils.IsMakeCall(call)
			if got != tt.want {
				t.Errorf("utils.IsMakeCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMakeCallWithLength(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		minArgs int
		want    bool
	}{
		{
			name:    "make with length",
			code:    `make([]int, 10)`,
			minArgs: 2,
			want:    true,
		},
		{
			name:    "make with length and capacity",
			code:    `make([]int, 10, 20)`,
			minArgs: 3,
			want:    true,
		},
		{
			name:    "make without enough args",
			code:    `make(map[string]int)`,
			minArgs: 2,
			want:    false,
		},
		{
			name:    "not make call",
			code:    `append(slice, 1, 2, 3)`,
			minArgs: 2,
			want:    false,
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

			got := utils.IsMakeCallWithLength(call, tt.minArgs)
			if got != tt.want {
				t.Errorf("utils.IsMakeCallWithLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMakeSliceCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "make slice with length",
			code: `make([]int, 10)`,
			want: true,
		},
		{
			name: "make slice with capacity",
			code: `make([]string, 0, 10)`,
			want: true,
		},
		{
			name: "make map",
			code: `make(map[string]int)`,
			want: false,
		},
		{
			name: "make channel",
			code: `make(chan int)`,
			want: false,
		},
		{
			name: "not make call",
			code: `append([]int{}, 1)`,
			want: false,
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

			got := utils.IsMakeSliceCall(call)
			if got != tt.want {
				t.Errorf("utils.IsMakeSliceCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMakeMapCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "make map without capacity",
			code: `make(map[string]int)`,
			want: true,
		},
		{
			name: "make map with capacity",
			code: `make(map[string]bool, 100)`,
			want: true,
		},
		{
			name: "make slice",
			code: `make([]int, 10)`,
			want: false,
		},
		{
			name: "make channel",
			code: `make(chan int)`,
			want: false,
		},
		{
			name: "not make call",
			code: `len(myMap)`,
			want: false,
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

			got := utils.IsMakeMapCall(call)
			if got != tt.want {
				t.Errorf("utils.IsMakeMapCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMakeByteSliceCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "make byte slice",
			code: `make([]byte, 1024)`,
			want: true,
		},
		{
			name: "make byte slice with capacity",
			code: `make([]byte, 0, 1024)`,
			want: true,
		},
		{
			name: "make int slice",
			code: `make([]int, 10)`,
			want: false,
		},
		{
			name: "make string slice",
			code: `make([]string, 0)`,
			want: false,
		},
		{
			name: "not make call",
			code: `[]byte("hello")`,
			want: false,
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

			got := utils.IsMakeByteSliceCall(call)
			if got != tt.want {
				t.Errorf("utils.IsMakeByteSliceCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsMakeCallsWithNoArgs tests make calls without arguments
func TestIsMakeCallsWithNoArgs(t *testing.T) {
	// Create a make call without args artificially
	makeCall := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: []ast.Expr{}, // Empty args
	}

	tests := []struct {
		name     string
		callFunc func(*ast.CallExpr) bool
		funcName string
	}{
		{"IsMakeSliceCall with no args", utils.IsMakeSliceCall, "IsMakeSliceCall"},
		{"IsMakeMapCall with no args", utils.IsMakeMapCall, "IsMakeMapCall"},
		{"IsMakeByteSliceCall with no args", utils.IsMakeByteSliceCall, "IsMakeByteSliceCall"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Appel fonction
			got := tt.callFunc(makeCall)
			// Vérification résultat
			if got != false {
				t.Errorf("%s(make with no args) = %v, want false", tt.funcName, got)
			}
		})
	}
}
