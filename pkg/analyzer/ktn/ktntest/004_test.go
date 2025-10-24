package ktntest

import (
	"go/ast"
	"testing"
)

// TestIsTestFunction tests isTestFunction helper
func TestIsTestFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		{"test function", "TestFoo", true},
		{"test with underscore", "Test_Foo", true},
		{"regular function", "Foo", false},
		{"helper", "testHelper", false},
		{"benchmark not detected", "BenchmarkFoo", false}, // isTestFunction only checks "Test"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: tt.funcName},
			}

			got := isTestFunction(funcDecl)
			if got != tt.want {
				t.Errorf("isTestFunction(%q) = %v, want %v", tt.funcName, got, tt.want)
			}
		})
	}
}

// TestIsErrorIndicatorName tests isErrorIndicatorName helper
func TestIsErrorIndicatorName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"wantErr", true},      // contains "err"
		{"expectErr", true},    // contains "err"
		{"expectedErr", true},  // contains "err"
		{"shouldErr", true},    // contains "err"
		{"wantError", true},    // contains "error"
		{"expectError", true},  // contains "error"
		{"hasError", true},     // contains "error"
		{"isError", true},      // contains "error"
		{"err", true},          // contains "err"
		{"error", true},        // contains "error"
		{"invalid", true},      // contains "invalid"
		{"fail", true},         // contains "fail"
		{"want", false},        // no error indicator
		{"value", false},       // no error indicator
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isErrorIndicatorName(tt.name)
			if got != tt.want {
				t.Errorf("isErrorIndicatorName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// Note: hasErrorTestCases is tested indirectly via the main analyzer tests
// Direct unit testing is complex due to AST structure requirements
