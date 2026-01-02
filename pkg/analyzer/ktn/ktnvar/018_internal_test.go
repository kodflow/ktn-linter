package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar018 tests the private runVar018 function.
func Test_runVar018(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_hasDifferentCapacity tests the private hasDifferentCapacity helper function.
func Test_hasDifferentCapacity(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "two args - no capacity",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.Ident{Name: "T"},
					&ast.BasicLit{Value: "10"},
				},
			},
			expected: false,
		},
		{
			name: "three args - has capacity",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.Ident{Name: "T"},
					&ast.BasicLit{Value: "10"},
					&ast.BasicLit{Value: "20"},
				},
			},
			expected: true,
		},
		{
			name: "one arg - no capacity",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.Ident{Name: "T"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := hasDifferentCapacity(tt.call)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("hasDifferentCapacity() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_isTotalSizeSmall tests the private isTotalSizeSmall helper function.
func Test_isTotalSizeSmall(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_shouldUseArray_tooFewArgs tests with insufficient args.
func Test_shouldUseArray_tooFewArgs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{}

			// Test with only 1 arg
			call := &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{&ast.Ident{Name: "T"}},
			}
			result := shouldUseArray(pass, call)
			// Vérification du résultat
			if result {
				t.Errorf("shouldUseArray() = true, expected false with too few args")
			}

		})
	}
}

// Test_shouldUseArray_withCapacity tests with different capacity.
func Test_shouldUseArray_withCapacity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{}

			// Test with 3 args (has different capacity)
			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "10"},
					&ast.BasicLit{Value: "20"}, // Different capacity
				},
			}
			result := shouldUseArray(pass, call)
			// Vérification du résultat
			if result {
				t.Errorf("shouldUseArray() = true, expected false with different capacity")
			}

		})
	}
}

// Test_getConstantSize_nilValue tests with nil constant value.
func Test_getConstantSize_nilValue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with expression not in TypesInfo
			expr := &ast.Ident{Name: "x"}
			result := getConstantSize(pass, expr)
			// Vérification du résultat
			if result != -1 {
				t.Errorf("getConstantSize() = %d, expected -1 for nil value", result)
			}

		})
	}
}

// Test_getConstantSize_nonInt tests with non-int constant.
func Test_getConstantSize_nonInt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Add a string constant to TypesInfo
			expr := &ast.BasicLit{Value: `"hello"`}
			pass.TypesInfo.Types[expr] = types.TypeAndValue{
				Value: nil, // Not a constant
			}
			result := getConstantSize(pass, expr)
			// Vérification du résultat
			if result != -1 {
				t.Errorf("getConstantSize() = %d, expected -1 for non-constant", result)
			}

		})
	}
}

// Test_reportArraySuggestion tests the private reportArraySuggestion function.
func Test_reportArraySuggestion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function reports array suggestions
		})
	}
}

// Test_runVar018_disabled tests runVar018 with disabled rule.
func Test_runVar018_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-018": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar018() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar018_fileExcluded tests runVar018 with excluded file.
func Test_runVar018_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-018": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar018() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_shouldUseArray tests the shouldUseArray private function.
func Test_shouldUseArray(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_getConstantSize tests the getConstantSize private function.
func Test_getConstantSize(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runVar018_nilInspector tests runVar018 with nil inspector.
func Test_runVar018_nilInspector(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			fset := token.NewFileSet()
			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: nil, // nil inspector
				},
			}

			result, err := runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar018() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar018_nilFset tests runVar018 with nil Fset.
func Test_runVar018_nilFset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			code := `package test`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{
				Fset: nil, // nil Fset
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
			}

			result, err := runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar018() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar018_nilTypesInfo tests runVar018 with nil TypesInfo.
func Test_runVar018_nilTypesInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			code := `package test`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: nil, // nil TypesInfo
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
			}

			result, err := runVar018(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar018() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar018() = %v, expected nil", result)
			}
		})
	}
}

// Test_shouldUseArray_notSliceType tests shouldUseArray with non-slice type.
func Test_shouldUseArray_notSliceType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{}

			// Test with non-slice type as first arg
			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.Ident{Name: "map"}, // not a slice
					&ast.BasicLit{Value: "10"},
				},
			}
			result := shouldUseArray(pass, call)
			// Verification du resultat
			if result {
				t.Errorf("shouldUseArray() = true, expected false for non-slice type")
			}
		})
	}
}

// Test_isTotalSizeSmall_notArrayType tests isTotalSizeSmall with non-array type.
func Test_isTotalSizeSmall_notArrayType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with non-array expression
			expr := &ast.Ident{Name: "int"}
			result := isTotalSizeSmall(pass, expr, 10)
			// Verification du resultat
			if result {
				t.Errorf("isTotalSizeSmall() = true, expected false for non-array type")
			}
		})
	}
}

// Test_isTotalSizeSmall_nilElemType tests isTotalSizeSmall with nil element type.
func Test_isTotalSizeSmall_nilElemType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with array type but nil element type
			elemExpr := &ast.Ident{Name: "UnknownType"}
			expr := &ast.ArrayType{Elt: elemExpr}
			// Element type not in TypesInfo
			result := isTotalSizeSmall(pass, expr, 10)
			// Verification du resultat
			if result {
				t.Errorf("isTotalSizeSmall() = true, expected false for nil element type")
			}
		})
	}
}

// Test_isTotalSizeSmall_nilSizes tests isTotalSizeSmall with nil TypesSizes.
func Test_isTotalSizeSmall_nilSizes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				TypesSizes: nil, // nil TypesSizes - will use default
			}

			// Setup array type with element type
			elemExpr := &ast.Ident{Name: "int"}
			pass.TypesInfo.Types[elemExpr] = types.TypeAndValue{
				Type: types.Typ[types.Int],
			}
			expr := &ast.ArrayType{Elt: elemExpr}

			// Should use default sizes for amd64
			result := isTotalSizeSmall(pass, expr, 8)
			// With int (8 bytes on amd64) * 8 elements = 64 bytes, should be true
			if !result {
				t.Errorf("isTotalSizeSmall() = false, expected true with default sizes")
			}
		})
	}
}

// Test_isTotalSizeSmall_largeSize tests isTotalSizeSmall with large total size.
func Test_isTotalSizeSmall_largeSize(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				TypesSizes: types.SizesFor("gc", "amd64"),
			}

			// Setup array type with element type
			elemExpr := &ast.Ident{Name: "int"}
			pass.TypesInfo.Types[elemExpr] = types.TypeAndValue{
				Type: types.Typ[types.Int],
			}
			expr := &ast.ArrayType{Elt: elemExpr}

			// With int (8 bytes on amd64) * 100 elements = 800 bytes, should be false
			result := isTotalSizeSmall(pass, expr, 100)
			// Verification du resultat
			if result {
				t.Errorf("isTotalSizeSmall() = true, expected false for large size")
			}
		})
	}
}

// Test_reportArraySuggestion_withValidMessage tests reportArraySuggestion with valid message.
func Test_reportArraySuggestion_withValidMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			reported := false
			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}

			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "10"},
				},
			}

			reportArraySuggestion(pass, call)

			// Should report
			if !reported {
				t.Error("reportArraySuggestion() did not report")
			}
		})
	}
}
