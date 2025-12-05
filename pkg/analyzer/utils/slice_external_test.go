package utils_test

import (
	"go/ast"
	"go/constant"
	"go/parser"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
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

			got := utils.IsSliceType(expr)
			if got != tt.want {
				t.Errorf("utils.IsSliceType() = %v, want %v", got, tt.want)
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

			got := utils.IsMapType(expr)
			if got != tt.want {
				t.Errorf("utils.IsMapType() = %v, want %v", got, tt.want)
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

			got := utils.IsSliceOrMapType(expr)
			if got != tt.want {
				t.Errorf("utils.IsSliceOrMapType() = %v, want %v", got, tt.want)
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

			got := utils.IsEmptySliceLiteral(lit)
			if got != tt.want {
				t.Errorf("utils.IsEmptySliceLiteral() = %v, want %v", got, tt.want)
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

			got := utils.IsByteSlice(expr)
			if got != tt.want {
				t.Errorf("utils.IsByteSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test functions with Pass parameter using a mock Pass
func TestIsSliceTypeWithPass(t *testing.T) {
	tests := []struct {
		name         string
		withTypeInfo bool
	}{
		{name: "with TypesInfo", withTypeInfo: true},
		{name: "fallback to AST", withTypeInfo: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			sliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}}
			if tt.withTypeInfo {
				pass.TypesInfo.Types[sliceExpr] = types.TypeAndValue{
					Type: types.NewSlice(types.Typ[types.Int]),
				}
			}

			if !utils.IsSliceTypeWithPass(pass, sliceExpr) {
				t.Errorf("utils.IsSliceTypeWithPass() = false, want true")
			}
		})
	}
}

func TestIsMapTypeWithPass(t *testing.T) {
	tests := []struct {
		name         string
		withTypeInfo bool
	}{
		{name: "with TypesInfo", withTypeInfo: true},
		{name: "fallback to AST", withTypeInfo: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			mapExpr := &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			}
			if tt.withTypeInfo {
				pass.TypesInfo.Types[mapExpr] = types.TypeAndValue{
					Type: types.NewMap(types.Typ[types.String], types.Typ[types.Int]),
				}
			}

			if !utils.IsMapTypeWithPass(pass, mapExpr) {
				t.Errorf("utils.IsMapTypeWithPass() = false, want true")
			}
		})
	}
}

func TestIsByteSliceWithPass(t *testing.T) {
	tests := []struct {
		name         string
		withTypeInfo bool
	}{
		{name: "with TypesInfo", withTypeInfo: true},
		{name: "fallback to AST", withTypeInfo: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			byteSliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "byte"}}
			if tt.withTypeInfo {
				pass.TypesInfo.Types[byteSliceExpr] = types.TypeAndValue{
					Type: types.NewSlice(types.Typ[types.Byte]),
				}
			}

			if !utils.IsByteSliceWithPass(pass, byteSliceExpr) {
				t.Errorf("utils.IsByteSliceWithPass() = false, want true")
			}
		})
	}
}

func TestIsSliceOrMapTypeWithPass(t *testing.T) {
	tests := []struct {
		name    string
		isSlice bool
	}{
		{name: "with slice", isSlice: true},
		{name: "with map", isSlice: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			var expr ast.Expr
			if tt.isSlice {
				sliceExpr := &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}}
				pass.TypesInfo.Types[sliceExpr] = types.TypeAndValue{
					Type: types.NewSlice(types.Typ[types.Int]),
				}
				expr = sliceExpr
			} else {
				mapExpr := &ast.MapType{
					Key:   &ast.Ident{Name: "string"},
					Value: &ast.Ident{Name: "int"},
				}
				pass.TypesInfo.Types[mapExpr] = types.TypeAndValue{
					Type: types.NewMap(types.Typ[types.String], types.Typ[types.Int]),
				}
				expr = mapExpr
			}

			if !utils.IsSliceOrMapTypeWithPass(pass, expr) {
				t.Errorf("utils.IsSliceOrMapTypeWithPass() = false, want true")
			}
		})
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

	got := utils.IsByteSlice(sliceExpr)
	if got != false {
		t.Errorf("utils.IsByteSlice([]pkg.Type) = %v, want false", got)
	}
}

// TestIsSmallConstantSize tests the IsSmallConstantSize function
func TestIsSmallConstantSize(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "nil value", want: false},
		{name: "non-constant", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Create a variable expression (non-constant)
			expr := &ast.Ident{Name: "n"}

			got := utils.IsSmallConstantSize(pass, expr)
			if got != tt.want {
				t.Errorf("utils.IsSmallConstantSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsSmallConstantSizeWithConstants tests with actual constant values
func TestIsSmallConstantSizeWithConstants(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  bool
	}{
		{"zero", 0, false},
		{"small positive 100", 100, true},
		{"boundary 1024", 1024, true},
		{"over boundary 1025", 1025, false},
		{"negative", -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &ast.BasicLit{Value: "x"}
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: map[ast.Expr]types.TypeAndValue{
						expr: {Value: constant.MakeInt64(tt.value)},
					},
				},
			}
			got := utils.IsSmallConstantSize(pass, expr)
			// Check result
			if got != tt.want {
				t.Errorf("IsSmallConstantSize(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
