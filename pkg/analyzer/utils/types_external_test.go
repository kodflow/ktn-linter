package utils_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
)

// TestIsZeroLiteral tests the functionality of the corresponding implementation.
func TestIsZeroLiteral(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"zero literal", "0", true},
		{"non-zero literal", "5", false},
		{"float zero", "0.0", false},  // FLOAT type, not INT
		{"string", `"0"`, false},      // STRING type, not INT
		{"hex zero", "0x0", false},    // INT but value "0x0", not "0"
		{"ident not lit", "x", false}, // Not a BasicLit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := utils.IsZeroLiteral(expr)
			if got != tt.expected {
				t.Errorf("utils.IsZeroLiteral(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

// TestIsReferenceType tests the functionality of the corresponding implementation.
func TestIsReferenceType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"slice type", "[]int", true},
		{"map type", "map[string]int", true},
		{"chan type", "chan int", true},
		{"array type", "[5]int", false},
		{"struct type", "struct{}", false},
		{"ident with map", "mymap", true},
		{"ident with chan", "mychan", true},
		{"ident with slice", "myslice", true},
		{"regular ident", "myint", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := utils.IsReferenceType(expr)
			if got != tt.expected {
				t.Errorf("utils.IsReferenceType(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

// TestIsStructType tests the functionality of the corresponding implementation.
func TestIsStructType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"struct type", "struct{}", true},
		{"ident (named type)", "MyStruct", true},
		{"selector (imported type)", "pkg.MyStruct", true},
		{"slice type", "[]int", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := utils.IsStructType(expr)
			if got != tt.expected {
				t.Errorf("utils.IsStructType(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

// TestGetTypeName tests the functionality of the corresponding implementation.
func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"slice", "[]int", "[]int"},
		{"map", "map[string]int", "map[string]int"},
		{"chan", "chan int", "chan int"},
		{"ident", "MyType", "MyType"},
		{"selector", "pkg.Type", "pkg.Type"},
		{"pointer", "*int", "*int"},
		{"nested slice", "[][]string", "[][]string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := utils.GetTypeName(expr)
			if got != tt.expected {
				t.Errorf("utils.GetTypeName(%s) = %s, want %s", tt.code, got, tt.expected)
			}
		})
	}
}

// TestGetTypeNameWithUnsupportedType tests the functionality of the corresponding implementation.
func TestGetTypeNameWithUnsupportedType(t *testing.T) {
	// Test avec un type non supporté
	expr := &ast.FuncType{
		Params: &ast.FieldList{},
	}
	got := utils.GetTypeName(expr)
	if got != "T" {
		t.Errorf("utils.GetTypeName(unsupported) = %s, want T", got)
	}
}

// TestIsMakeSliceZero tests the functionality of the corresponding implementation.
func TestIsMakeSliceZero(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"make slice with 0", "make([]int, 0)", true},
		{"make slice with 0,0", "make([]int, 0, 0)", true},
		{"make slice with non-zero length", "make([]int, 5)", false},
		{"make slice with non-zero capacity", "make([]int, 0, 10)", false},
		{"make map", "make(map[string]int)", false},
		{"not make", "[]int{}", false},
		{"make with one arg", "make([]int)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := utils.IsMakeSliceZero(expr)
			if got != tt.expected {
				t.Errorf("utils.IsMakeSliceZero(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

// TestIsMakeSliceZeroWithNonCallExpr tests the functionality of the corresponding implementation.
func TestIsMakeSliceZeroWithNonCallExpr(t *testing.T) {
	// Test avec une expression qui n'est pas un appel
	expr := &ast.Ident{Name: "test"}
	got := utils.IsMakeSliceZero(expr)
	if got != false {
		t.Errorf("utils.IsMakeSliceZero(non-call) = %v, want false", got)
	}
}

// TestIsMakeSliceZeroWithNonMake tests the functionality of the corresponding implementation.
func TestIsMakeSliceZeroWithNonMake(t *testing.T) {
	// Test avec un appel qui n'est pas make
	fset := token.NewFileSet()
	code := "append(s, 1)"
	expr, _ := parser.ParseExprFrom(fset, "", code, 0)
	got := utils.IsMakeSliceZero(expr)
	if got != false {
		t.Errorf("utils.IsMakeSliceZero(non-make) = %v, want false", got)
	}
}

// TestHasPositiveLength tests the functionality of HasPositiveLength
func TestHasPositiveLength(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"positive literal", "5", true},
		{"zero literal", "0", false},
		{"variable", "n", true},     // Variables assumed positive
		{"expression", "1+1", true}, // Non-literal assumed positive
		{"hex zero", "0x0", true},   // Not "0" string, treated as positive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			// Call with nil pass to test AST-only path
			got := utils.HasPositiveLength(nil, expr)
			if got != tt.expected {
				t.Errorf("utils.HasPositiveLength(%q) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}

	// Test avec pass non-nil mais TypesInfo nil
	t.Run("pass without typesinfo", func(t *testing.T) {
		pass := &analysis.Pass{
			TypesInfo: nil,
		}
		expr, _ := parser.ParseExpr("5")
		got := utils.HasPositiveLength(pass, expr)
		if !got {
			t.Errorf("HasPositiveLength with pass but no TypesInfo should return true")
		}
	})
}
