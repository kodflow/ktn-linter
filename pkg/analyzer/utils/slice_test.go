package utils

import (
	"go/ast"
	"go/parser"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestIsSliceType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "slice type",
			code: `[]int`,
			want: true,
		},
		{
			name: "array type",
			code: `[5]int`,
			want: false,
		},
		{
			name: "map type",
			code: `map[string]int`,
			want: false,
		},
		{
			name: "channel type",
			code: `chan int`,
			want: false,
		},
		{
			name: "identifier",
			code: `MyType`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := IsSliceType(expr)
			if got != tt.want {
				t.Errorf("IsSliceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMapType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "map type",
			code: `map[string]int`,
			want: true,
		},
		{
			name: "slice type",
			code: `[]int`,
			want: false,
		},
		{
			name: "array type",
			code: `[5]int`,
			want: false,
		},
		{
			name: "channel type",
			code: `chan int`,
			want: false,
		},
		{
			name: "identifier",
			code: `MyType`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := IsMapType(expr)
			if got != tt.want {
				t.Errorf("IsMapType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSliceOrMapType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "slice type",
			code: `[]int`,
			want: true,
		},
		{
			name: "map type",
			code: `map[string]int`,
			want: true,
		},
		{
			name: "array type",
			code: `[5]int`,
			want: false,
		},
		{
			name: "channel type",
			code: `chan int`,
			want: false,
		},
		{
			name: "identifier",
			code: `MyType`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := IsSliceOrMapType(expr)
			if got != tt.want {
				t.Errorf("IsSliceOrMapType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptySliceLiteral(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "empty slice literal",
			code: `[]int{}`,
			want: true,
		},
		{
			name: "non-empty slice literal",
			code: `[]int{1, 2, 3}`,
			want: false,
		},
		{
			name: "empty map literal",
			code: `map[string]int{}`,
			want: false,
		},
		{
			name: "empty struct literal",
			code: `MyStruct{}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			lit, ok := expr.(*ast.CompositeLit)
			if !ok {
				t.Fatalf("Expression is not a CompositeLit")
			}

			got := IsEmptySliceLiteral(lit)
			if got != tt.want {
				t.Errorf("IsEmptySliceLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsByteSlice(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "byte slice",
			code: `[]byte`,
			want: true,
		},
		{
			name: "uint8 slice",
			code: `[]uint8`,
			want: true,
		},
		{
			name: "int slice",
			code: `[]int`,
			want: false,
		},
		{
			name: "string slice",
			code: `[]string`,
			want: false,
		},
		{
			name: "array of bytes",
			code: `[10]byte`,
			want: false,
		},
		{
			name: "map type",
			code: `map[string]byte`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			got := IsByteSlice(expr)
			if got != tt.want {
				t.Errorf("IsByteSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test functions with Pass parameter using a mock Pass
func TestIsSliceTypeWithPass(t *testing.T) {
	// Create a minimal Pass with TypesInfo
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	// Test with TypesInfo
	sliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}}
	pass.TypesInfo.Types[sliceExpr] = types.TypeAndValue{
		Type: types.NewSlice(types.Typ[types.Int]),
	}

	got := IsSliceTypeWithPass(pass, sliceExpr)
	if !got {
		t.Errorf("IsSliceTypeWithPass() with TypesInfo = false, want true")
	}

	// Test fallback to AST checking
	pass2 := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	got2 := IsSliceTypeWithPass(pass2, sliceExpr)
	if !got2 {
		t.Errorf("IsSliceTypeWithPass() fallback to AST = false, want true")
	}
}

func TestIsMapTypeWithPass(t *testing.T) {
	// Create a minimal Pass with TypesInfo
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	// Test with TypesInfo
	mapExpr := &ast.MapType{
		Key:   &ast.Ident{Name: "string"},
		Value: &ast.Ident{Name: "int"},
	}
	pass.TypesInfo.Types[mapExpr] = types.TypeAndValue{
		Type: types.NewMap(types.Typ[types.String], types.Typ[types.Int]),
	}

	got := IsMapTypeWithPass(pass, mapExpr)
	if !got {
		t.Errorf("IsMapTypeWithPass() with TypesInfo = false, want true")
	}

	// Test fallback to AST checking
	pass2 := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	got2 := IsMapTypeWithPass(pass2, mapExpr)
	if !got2 {
		t.Errorf("IsMapTypeWithPass() fallback to AST = false, want true")
	}
}

func TestIsByteSliceWithPass(t *testing.T) {
	// Create a minimal Pass with TypesInfo
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	// Test with TypesInfo
	byteSliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "byte"}}
	pass.TypesInfo.Types[byteSliceExpr] = types.TypeAndValue{
		Type: types.NewSlice(types.Typ[types.Byte]),
	}

	got := IsByteSliceWithPass(pass, byteSliceExpr)
	if !got {
		t.Errorf("IsByteSliceWithPass() with TypesInfo = false, want true")
	}

	// Test fallback to AST checking
	pass2 := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	got2 := IsByteSliceWithPass(pass2, byteSliceExpr)
	if !got2 {
		t.Errorf("IsByteSliceWithPass() fallback to AST = false, want true")
	}
}

func TestIsSliceOrMapTypeWithPass(t *testing.T) {
	// Create a minimal Pass with TypesInfo
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	// Test with slice
	sliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}}
	pass.TypesInfo.Types[sliceExpr] = types.TypeAndValue{
		Type: types.NewSlice(types.Typ[types.Int]),
	}

	got := IsSliceOrMapTypeWithPass(pass, sliceExpr)
	if !got {
		t.Errorf("IsSliceOrMapTypeWithPass() with slice = false, want true")
	}

	// Test with map
	mapExpr := &ast.MapType{
		Key:   &ast.Ident{Name: "string"},
		Value: &ast.Ident{Name: "int"},
	}
	pass.TypesInfo.Types[mapExpr] = types.TypeAndValue{
		Type: types.NewMap(types.Typ[types.String], types.Typ[types.Int]),
	}

	got2 := IsSliceOrMapTypeWithPass(pass, mapExpr)
	if !got2 {
		t.Errorf("IsSliceOrMapTypeWithPass() with map = false, want true")
	}
}
// TestIsByteSliceNonIdent tests IsByteSlice with non-identifier element type
func TestIsByteSliceNonIdent(t *testing.T) {
	// Create a slice with SelectorExpr element (e.g., pkg.Type)
	sliceExpr := &ast.ArrayType{
		Len: nil,
		Elt: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "Type"},
		},
	}

	got := IsByteSlice(sliceExpr)
	if got != false {
		t.Errorf("IsByteSlice([]pkg.Type) = %v, want false", got)
	}
}
