package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIsZeroLiteral(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"zero literal", "0", true},
		{"non-zero literal", "5", false},
		{"float zero", "0.0", false}, // FLOAT type, not INT
		{"string", `"0"`, false},     // STRING type, not INT
		{"hex zero", "0x0", false},   // INT but value "0x0", not "0"
		{"ident not lit", "x", false}, // Not a BasicLit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := IsZeroLiteral(expr)
			if got != tt.expected {
				t.Errorf("IsZeroLiteral(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

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
			got := IsReferenceType(expr)
			if got != tt.expected {
				t.Errorf("IsReferenceType(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

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
			got := IsStructType(expr)
			if got != tt.expected {
				t.Errorf("IsStructType(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestIsSliceType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"slice type", "[]int", true},
		{"array type", "[5]int", false},
		{"map type", "map[string]int", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, _ := parser.ParseExpr(tt.code)
			got := IsSliceType(expr)
			if got != tt.expected {
				t.Errorf("IsSliceType(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

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
			got := GetTypeName(expr)
			if got != tt.expected {
				t.Errorf("GetTypeName(%s) = %s, want %s", tt.code, got, tt.expected)
			}
		})
	}
}

func TestGetTypeNameWithUnsupportedType(t *testing.T) {
	// Test avec un type non supporté
	expr := &ast.FuncType{
		Params: &ast.FieldList{},
	}
	got := GetTypeName(expr)
	if got != "T" {
		t.Errorf("GetTypeName(unsupported) = %s, want T", got)
	}
}

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
			got := IsMakeSliceZero(expr)
			if got != tt.expected {
				t.Errorf("IsMakeSliceZero(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestIsMakeSliceZeroWithNonCallExpr(t *testing.T) {
	// Test avec une expression qui n'est pas un appel
	expr := &ast.Ident{Name: "test"}
	got := IsMakeSliceZero(expr)
	if got != false {
		t.Errorf("IsMakeSliceZero(non-call) = %v, want false", got)
	}
}

func TestIsMakeSliceZeroWithNonMake(t *testing.T) {
	// Test avec un appel qui n'est pas make
	fset := token.NewFileSet()
	code := "append(s, 1)"
	expr, _ := parser.ParseExprFrom(fset, "", code, 0)
	got := IsMakeSliceZero(expr)
	if got != false {
		t.Errorf("IsMakeSliceZero(non-make) = %v, want false", got)
	}
}
