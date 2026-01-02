package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"runtime"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar013 tests the private runVar013 function.
func Test_runVar013(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name: "valid code without large structs",
			code: `package test

func foo(x int) {}
`,
			expectError: false,
		},
		{
			name: "code with receiver",
			code: `package test

type T struct{ x int }
func (t T) foo() {}
`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Reset config for clean state
			config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})

			// Properly type-check the code to populate TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Defs:  make(map[*ast.Ident]types.Object),
				Uses:  make(map[*ast.Ident]types.Object),
			}
			conf := types.Config{}
			pkg, checkErr := conf.Check("test", fset, []*ast.File{file}, info)
			// Fail fast on type-check errors
			if checkErr != nil {
				t.Fatalf("failed to type-check: %v", checkErr)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Pkg:   pkg,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				TypesInfo:  info,
				TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
				Report:     func(_d analysis.Diagnostic) {},
			}

			_, err = runVar013(pass)
			// Verify error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("runVar013() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// Test_checkFuncParams009 tests the private checkFuncParams009 function.
func Test_checkFuncParams009(t *testing.T) {
	tests := []struct {
		name         string
		params       *ast.FieldList
		expectReport bool
	}{
		{
			name:         "nil params",
			params:       nil,
			expectReport: false,
		},
		{
			name: "empty params list",
			params: &ast.FieldList{
				List: []*ast.Field{},
			},
			expectReport: false,
		},
		{
			name: "single int param",
			params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "x"}},
						Type:  &ast.Ident{Name: "int"},
					},
				},
			},
			expectReport: false,
		},
		{
			name: "param without names",
			params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			expectReport: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reportCount := 0
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
					Defs:  make(map[*ast.Ident]types.Object),
					Uses:  make(map[*ast.Ident]types.Object),
				},
				TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkFuncParams009(pass, tt.params, 64, false)

			// Verify report expectation
			if (reportCount > 0) != tt.expectReport {
				t.Errorf("checkFuncParams009() reported %d, expectReport %v", reportCount, tt.expectReport)
			}
		})
	}
}

// Test_checkParamType009 tests the private checkParamType009 function.
func Test_checkParamType009(t *testing.T) {
	tests := []struct {
		name         string
		typ          ast.Expr
		expectReport bool
	}{
		{
			name:         "pointer type skipped",
			typ:          &ast.StarExpr{X: &ast.Ident{Name: "BigStruct"}},
			expectReport: false,
		},
		{
			name:         "unknown type skipped",
			typ:          &ast.Ident{Name: "UnknownType"},
			expectReport: false,
		},
		{
			name:         "basic int type skipped",
			typ:          &ast.Ident{Name: "int"},
			expectReport: false,
		},
		{
			name:         "variadic param unwrapped",
			typ:          &ast.Ellipsis{Elt: &ast.Ident{Name: "int"}},
			expectReport: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reportCount := 0
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
					Defs:  make(map[*ast.Ident]types.Object),
					Uses:  make(map[*ast.Ident]types.Object),
				},
				TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkParamType009(pass, tt.typ, token.NoPos, 64, false)

			// Verify report expectation
			if (reportCount > 0) != tt.expectReport {
				t.Errorf("checkParamType009() reported %d, expectReport %v", reportCount, tt.expectReport)
			}
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
		tt := tt // Capture range variable
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
			checkParamType009(pass, typ, token.NoPos, 64, false)
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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test with type that won't have TypeOf info
			typ := &ast.Ident{Name: "UnknownType"}
			checkParamType009(pass, typ, token.NoPos, 64, false)
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
		tt := tt // Capture range variable
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

			checkParamType009(pass, typeIdent, token.NoPos, 64, false)
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
		tt := tt // Capture range variable
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

			checkParamType009(pass, typeIdent, token.NoPos, 64, false)
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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

// Test_runVar013_disabled tests runVar013 with disabled rule.
func Test_runVar013_disabled(t *testing.T) {
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
					"KTN-VAR-013": {Enabled: config.Bool(false)},
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

			_, err = runVar013(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar013() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar013() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar013_fileExcluded tests runVar013 with excluded file.
func Test_runVar013_fileExcluded(t *testing.T) {
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
					"KTN-VAR-013": {
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

			_, err = runVar013(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar013() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar013() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_isExternalType009 tests the isExternalType009 private function.
func Test_isExternalType009(t *testing.T) {
	tests := []struct {
		name     string
		typeInfo types.Type
		passPkg  *types.Package
		expected bool
	}{
		{
			name:     "nil type returns false",
			typeInfo: nil,
			passPkg:  types.NewPackage("test/pkg", "pkg"),
			expected: false,
		},
		{
			name:     "basic type returns false",
			typeInfo: types.Typ[types.Int],
			passPkg:  types.NewPackage("test/pkg", "pkg"),
			expected: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				Pkg: tt.passPkg,
			}

			result := isExternalType009(tt.typeInfo, pass)
			// Verify result
			if result != tt.expected {
				t.Errorf("isExternalType009() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
