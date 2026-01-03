package ktnvar

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"math"
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

// Test_runVar013_nilTypesInfo tests runVar013 with nil TypesInfo.
func Test_runVar013_nilTypesInfo(t *testing.T) {
	// Reset config for clean state
	config.Reset()

	code := `package test
	func foo(x int) {}
	`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Vérifier l'erreur de parsing
	if err != nil || file == nil {
		t.Fatalf("failed to parse test code: %v", err)
	}
	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo: nil, // nil TypesInfo
		Report:    func(_d analysis.Diagnostic) {},
	}

	result, err := runVar013(pass)
	// Should return nil, nil for nil TypesInfo
	if err != nil {
		t.Errorf("runVar013() error = %v, expected nil", err)
	}
	// Verify result
	if result != nil {
		t.Errorf("runVar013() = %v, expected nil", result)
	}
}

// Test_checkFuncParams009_negativeMaxBytes tests with maxBytes <= 0.
func Test_checkFuncParams009_negativeMaxBytes(t *testing.T) {
	reportCount := 0
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
		TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	params := &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{{Name: "x"}},
				Type:  &ast.Ident{Name: "int"},
			},
		},
	}

	// Test with maxBytes = 0 (should be clamped to default)
	checkFuncParams009(pass, params, 0, false)

	// Test with maxBytes = -1 (should be clamped to default)
	checkFuncParams009(pass, params, -1, false)

	// Should not report for basic int type
	if reportCount != 0 {
		t.Errorf("checkFuncParams009() reported %d, expected 0", reportCount)
	}
}

// Test_getStructSize009 tests the private getStructSize009 function.
func Test_getStructSize009(t *testing.T) {
	tests := []struct {
		name         string
		setupPass    func() (*analysis.Pass, ast.Expr)
		expectedSize int64
	}{
		{
			name: "pointer type returns -1",
			setupPass: func() (*analysis.Pass, ast.Expr) {
				pkg := types.NewPackage("test", "test")
				pass := &analysis.Pass{
					Pkg: pkg,
					TypesInfo: &types.Info{
						Types: make(map[ast.Expr]types.TypeAndValue),
					},
				}
				// Pointer type
				return pass, &ast.StarExpr{X: &ast.Ident{Name: "int"}}
			},
			expectedSize: -1,
		},
		{
			name: "nil type info returns -1",
			setupPass: func() (*analysis.Pass, ast.Expr) {
				pkg := types.NewPackage("test", "test")
				pass := &analysis.Pass{
					Pkg: pkg,
					TypesInfo: &types.Info{
						Types: make(map[ast.Expr]types.TypeAndValue),
					},
				}
				// Unknown type
				return pass, &ast.Ident{Name: "UnknownType"}
			},
			expectedSize: -1,
		},
		{
			name: "non-struct type returns -1",
			setupPass: func() (*analysis.Pass, ast.Expr) {
				pkg := types.NewPackage("test", "test")
				typeIdent := &ast.Ident{Name: "int"}
				pass := &analysis.Pass{
					Pkg: pkg,
					TypesInfo: &types.Info{
						Types: map[ast.Expr]types.TypeAndValue{
							typeIdent: {Type: types.Typ[types.Int]},
						},
					},
				}
				return pass, typeIdent
			},
			expectedSize: -1,
		},
		{
			name: "struct with valid size",
			setupPass: func() (*analysis.Pass, ast.Expr) {
				pkg := types.NewPackage("test", "test")
				structType := types.NewStruct(
					[]*types.Var{
						types.NewVar(0, pkg, "a", types.Typ[types.Int64]),
					},
					nil,
				)
				obj := types.NewTypeName(0, pkg, "TestStruct", structType)
				namedType := types.NewNamed(obj, structType, nil)
				typeIdent := &ast.Ident{Name: "TestStruct"}
				pass := &analysis.Pass{
					Pkg: pkg,
					TypesInfo: &types.Info{
						Types: map[ast.Expr]types.TypeAndValue{
							typeIdent: {Type: namedType},
						},
					},
					TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
				}
				return pass, typeIdent
			},
			expectedSize: 8,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			pass, expr := tt.setupPass()
			size := getStructSize009(pass, expr)
			// Verify expected size
			if size != tt.expectedSize {
				t.Errorf("getStructSize009() = %d, expected %d", size, tt.expectedSize)
			}
		})
	}
}

// Test_getStructSize009_nilTypesSizes tests with nil TypesSizes (fallback path).
func Test_getStructSize009_nilTypesSizes(t *testing.T) {
	// Create a struct type
	pkg := types.NewPackage("test", "test")
	structType := types.NewStruct(
		[]*types.Var{
			types.NewVar(0, pkg, "a", types.Typ[types.Int]),
			types.NewVar(0, pkg, "b", types.Typ[types.Int]),
			types.NewVar(0, pkg, "c", types.Typ[types.Int]),
		},
		nil,
	)
	obj := types.NewTypeName(0, pkg, "TestStruct", structType)
	namedType := types.NewNamed(obj, structType, nil)

	typeIdent := &ast.Ident{Name: "TestStruct"}
	pass := &analysis.Pass{
		Pkg: pkg,
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{
				typeIdent: {Type: namedType},
			},
		},
		TypesSizes: nil, // nil TypesSizes - uses fallback
	}

	size := getStructSize009(pass, typeIdent)
	// Fallback: 8 bytes per field * 3 fields = 24 bytes
	if size != 24 {
		t.Errorf("getStructSize009() = %d, expected 24", size)
	}
}

// Test_isExternalType009_nilObjPkg tests with nil obj.Pkg().
func Test_isExternalType009_nilObjPkg(t *testing.T) {
	// Create a named type with nil package (universe scope type)
	obj := types.NewTypeName(0, nil, "error", nil)

	pass := &analysis.Pass{
		Pkg: types.NewPackage("test", "test"),
	}

	// Test with named type that has nil Pkg
	result := isExternalType009(obj.Type(), pass)
	// Should return false for nil Pkg
	if result {
		t.Errorf("isExternalType009() = true, expected false for nil Pkg")
	}
}

// Test_checkParamType009_largeStruct tests with large struct that triggers report.
func Test_checkParamType009_largeStruct(t *testing.T) {
	// Reset config for clean state
	config.Reset()

	// Create a large struct type (>64 bytes)
	pkg := types.NewPackage("test", "test")
	fields := make([]*types.Var, 10)
	// Create 10 int fields = 80 bytes on 64-bit
	for i := range 10 {
		fields[i] = types.NewVar(0, pkg, fmt.Sprintf("f%d", i), types.Typ[types.Int64])
	}
	structType := types.NewStruct(fields, nil)
	obj := types.NewTypeName(0, pkg, "BigStruct", structType)
	namedType := types.NewNamed(obj, structType, nil)

	typeIdent := &ast.Ident{Name: "BigStruct"}
	reportCount := 0
	pass := &analysis.Pass{
		Pkg: pkg,
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{
				typeIdent: {Type: namedType},
			},
		},
		TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	checkParamType009(pass, typeIdent, token.NoPos, 64, false)

	// Should report for large struct
	if reportCount != 1 {
		t.Errorf("checkParamType009() reported %d, expected 1", reportCount)
	}
}

// Test_checkParamType009_ellipsisNilElt tests with nil Ellipsis.Elt.
func Test_checkParamType009_ellipsisNilElt(t *testing.T) {
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with Ellipsis that has nil Elt
	typ := &ast.Ellipsis{Elt: nil}
	checkParamType009(pass, typ, token.NoPos, 64, false)
	// Should not panic
}

// Test_getStructSize009_zeroSizeFallback tests when Sizeof returns 0.
func Test_getStructSize009_zeroSizeFallback(t *testing.T) {
	// Create an empty struct (0 fields = 0 bytes)
	pkg := types.NewPackage("test", "test")
	structType := types.NewStruct(nil, nil) // Empty struct
	obj := types.NewTypeName(0, pkg, "EmptyStruct", structType)
	namedType := types.NewNamed(obj, structType, nil)

	typeIdent := &ast.Ident{Name: "EmptyStruct"}
	pass := &analysis.Pass{
		Pkg: pkg,
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{
				typeIdent: {Type: namedType},
			},
		},
		// Use real TypesSizes - empty struct returns 0
		TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
	}

	size := getStructSize009(pass, typeIdent)
	// Empty struct should return 0 from Sizeof, triggering fallback (0 fields * 8 = 0)
	if size != 0 {
		t.Errorf("getStructSize009() = %d, expected 0 for empty struct", size)
	}
}

// Test_isExternalType009_nilObj tests with named type that has nil Obj.
func Test_isExternalType009_nilObj(t *testing.T) {
	pass := &analysis.Pass{
		Pkg: types.NewPackage("test", "test"),
	}

	// Create a named type with nil obj - use interface type as a workaround
	// The predeclared error type has pkg = nil
	errorType := types.Universe.Lookup("error").Type()

	result := isExternalType009(errorType, pass)
	// Error is a named type but from universe scope (pkg = nil)
	// Should return false for nil Pkg
	if result {
		t.Errorf("isExternalType009() = true, expected false for universe type")
	}
}

// Test_runVar013_withFuncNoParams tests with function that has nil Params.
func Test_runVar013_withFuncNoParams(t *testing.T) {
	// Reset config for clean state
	config.Reset()

	code := `package test
	func foo() {}
	`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Type check the code
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, _ := conf.Check("test", fset, []*ast.File{file}, info)

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo:  info,
		TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	result, err := runVar013(pass)
	if err != nil {
		t.Errorf("runVar013() error = %v", err)
	}
	if result != nil {
		t.Errorf("runVar013() result = %v, expected nil", result)
	}
	// Should not report for function with no params
	if reportCount != 0 {
		t.Errorf("runVar013() reported %d, expected 0", reportCount)
	}
}

// Test_checkParamType009_displaySizeOverflow tests displaySize overflow guard.
func Test_checkParamType009_displaySizeOverflow(t *testing.T) {
	// Reset config for clean state
	config.Reset()

	// Create a mock struct type that reports very large size
	pkg := types.NewPackage("test", "test")
	structType := types.NewStruct(nil, nil)
	obj := types.NewTypeName(0, pkg, "HugeStruct", structType)
	namedType := types.NewNamed(obj, structType, nil)

	typeIdent := &ast.Ident{Name: "HugeStruct"}
	reportCount := 0
	var reportedMsg string

	// Create custom TypesSizes that returns a huge size
	hugeSizes := &mockHugeSizes{}

	pass := &analysis.Pass{
		Pkg: pkg,
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{
				typeIdent: {Type: namedType},
			},
		},
		TypesSizes: hugeSizes,
		Report: func(d analysis.Diagnostic) {
			reportCount++
			reportedMsg = d.Message
		},
	}

	// Should trigger overflow guard path
	checkParamType009(pass, typeIdent, token.NoPos, 64, false)

	// Should report for struct exceeding threshold
	if reportCount != 1 {
		t.Errorf("checkParamType009() reported %d, expected 1", reportCount)
	}
	// Verify message was generated (overflow path reached)
	if reportedMsg == "" {
		t.Error("checkParamType009() did not generate report message")
	}
}

// mockHugeSizes implements types.Sizes interface returning huge sizes.
type mockHugeSizes struct{}

func (m *mockHugeSizes) Alignof(T types.Type) int64 {
	return 8
}

func (m *mockHugeSizes) Offsetsof(fields []*types.Var) []int64 {
	offsets := make([]int64, len(fields))
	// Offsets are sequential
	for i := range fields {
		offsets[i] = int64(i) * 8
	}
	return offsets
}

func (m *mockHugeSizes) Sizeof(T types.Type) int64 {
	// Return a very large size to trigger overflow guard
	// On 64-bit systems math.MaxInt is 9223372036854775807
	// Return a value larger than 64 bytes but not overflowing
	return math.MaxInt64
}

// Test_getStructSize009_sizesNonPositive tests when Sizeof returns 0 for non-empty struct.
func Test_getStructSize009_sizesNonPositive(t *testing.T) {
	// Create a struct type with fields
	pkg := types.NewPackage("test", "test")
	structType := types.NewStruct(
		[]*types.Var{
			types.NewVar(0, pkg, "a", types.Typ[types.Int]),
			types.NewVar(0, pkg, "b", types.Typ[types.Int]),
		},
		nil,
	)
	obj := types.NewTypeName(0, pkg, "TestStruct", structType)
	namedType := types.NewNamed(obj, structType, nil)

	typeIdent := &ast.Ident{Name: "TestStruct"}

	// Create custom TypesSizes that returns 0 size
	zeroSizes := &mockZeroSizes{}

	pass := &analysis.Pass{
		Pkg: pkg,
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{
				typeIdent: {Type: namedType},
			},
		},
		TypesSizes: zeroSizes, // Returns 0 for Sizeof
	}

	size := getStructSize009(pass, typeIdent)
	// When Sizeof returns 0, fallback is used: 2 fields * 8 bytes = 16
	if size != 16 {
		t.Errorf("getStructSize009() = %d, expected 16 (fallback)", size)
	}
}

// mockZeroSizes implements types.Sizes returning 0 for Sizeof.
type mockZeroSizes struct{}

func (m *mockZeroSizes) Alignof(T types.Type) int64 {
	return 8
}

func (m *mockZeroSizes) Offsetsof(fields []*types.Var) []int64 {
	offsets := make([]int64, len(fields))
	return offsets
}

func (m *mockZeroSizes) Sizeof(T types.Type) int64 {
	return 0 // Return 0 to trigger fallback path
}

// Test_runVar013_fileExcludedWithLargeStruct tests file exclusion with large struct that would trigger.
func Test_runVar013_fileExcludedWithLargeStruct(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-013": {
				Exclude: []string{"excluded.go"},
			},
		},
	})
	defer config.Reset()

	// Code that would normally trigger the rule (large struct passed by value)
	code := `package test
type LargeStruct struct {
	a, b, c, d, e, f, g, h, i, j int64
}
func foo(s LargeStruct) {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Type check the code
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, _ := conf.Check("test", fset, []*ast.File{file}, info)

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo:  info,
		TypesSizes: types.SizesFor(runtime.Compiler, runtime.GOARCH),
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar013(pass)
	if err != nil {
		t.Errorf("runVar013() error = %v", err)
	}

	// Should not report when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar013() reported %d, expected 0 when file excluded", reportCount)
	}
}
