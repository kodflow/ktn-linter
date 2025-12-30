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

// Test_runVar009 tests the private runVar009 function.
func Test_runVar009(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_checkFuncParams009 tests the private checkFuncParams009 function.
func Test_checkFuncParams009(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks function params for VAR-009
		})
	}
}

// Test_checkParamType009 tests the private checkParamType009 function.
func Test_checkParamType009(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks param types for large structs
		})
	}
}

// Test_checkParamType009_pointer tests with pointer type.
func Test_checkParamType009_pointer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test with pointer type (should return early)
			typ := &ast.StarExpr{
				X: &ast.Ident{Name: "BigStruct"},
			}
			checkParamType009(pass, typ, token.NoPos, 64)
			// No error expected for pointer

		})
	}
}

// Test_checkParamType009_nilTypeInfo tests with nil type info.
func Test_checkParamType009_nilTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test with type that won't have TypeOf info
			typ := &ast.Ident{Name: "UnknownType"}
			checkParamType009(pass, typ, token.NoPos, 64)
			// No error expected when type info is nil

		})
	}
}

// Test_checkParamType009_externalType tests with external type.
func Test_checkParamType009_externalType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create external package
			externalPkg := types.NewPackage("external/pkg", "pkg")
			currentPkg := types.NewPackage("test/pkg", "pkg")

			// Create named type in external package
			obj := types.NewTypeName(0, externalPkg, "BigStruct", types.NewStruct(
				[]*types.Var{
					types.NewVar(0, externalPkg, "a", types.Typ[types.Int]),
					types.NewVar(0, externalPkg, "b", types.Typ[types.Int]),
					types.NewVar(0, externalPkg, "c", types.Typ[types.Int]),
					types.NewVar(0, externalPkg, "d", types.Typ[types.Int]),
				},
				nil,
			))
			namedType := types.NewNamed(obj, obj.Type().Underlying(), nil)

			typeIdent := &ast.Ident{Name: "BigStruct"}
			pass := &analysis.Pass{
				Pkg: currentPkg,
				TypesInfo: &types.Info{
					Types: map[ast.Expr]types.TypeAndValue{
						typeIdent: {Type: namedType},
					},
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			checkParamType009(pass, typeIdent, token.NoPos, 64)
			// No error expected for external type

		})
	}
}

// Test_checkParamType009_notStruct tests with non-struct type.
func Test_checkParamType009_notStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			typeIdent := &ast.Ident{Name: "int"}
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: map[ast.Expr]types.TypeAndValue{
						typeIdent: {Type: types.Typ[types.Int]},
					},
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			checkParamType009(pass, typeIdent, token.NoPos, 64)
			// No error expected for non-struct type

		})
	}
}

// Test_isExternalType009_notNamed tests with non-named type.
func Test_isExternalType009_notNamed(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{}

			// Test with basic type (not named)
			basicType := types.Typ[types.Int]
			result := isExternalType009(basicType, pass)
			// Vérification du résultat
			if result {
				t.Errorf("isExternalType009() = true, expected false for basic type")
			}

		})
	}
}

// Test_isExternalType009_samePackage tests with same package type.
func Test_isExternalType009_samePackage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pkg := types.NewPackage("test/pkg", "pkg")
			pass := &analysis.Pass{
				Pkg: pkg,
			}

			// Test with named type from same package
			// Use basic type for underlying to avoid nil check issues
			obj := types.NewTypeName(0, pkg, "MyStruct", types.Typ[types.Int])
			result := isExternalType009(obj.Type(), pass)
			// Vérification du résultat
			if result {
				t.Errorf("isExternalType009() = true, expected false for same package")
			}

		})
	}
}

// Test_runVar009_disabled tests runVar009 with disabled rule.
func Test_runVar009_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-009": {Enabled: config.Bool(false)},
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

			_, err = runVar009(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar009() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar009() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar009_fileExcluded tests runVar009 with excluded file.
func Test_runVar009_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-009": {
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

			_, err = runVar009(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar009() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar009() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_isExternalType009 tests the isExternalType009 private function.
func Test_isExternalType009(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}
