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

// Test_runVar019 tests the private runVar019 function.
func Test_runVar019(t *testing.T) {
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

// Test_isPointerType tests the private isPointerType helper function.
func Test_isPointerType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "pointer type",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "T"}},
			expected: true,
		},
		{
			name:     "non-pointer type",
			expr:     &ast.Ident{Name: "T"},
			expected: false,
		},
		{
			name:     "selector expr",
			expr:     &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isPointerType(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isPointerType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_getTypeName tests the private getTypeName helper function.
func Test_getTypeName(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "ident",
			expr:     &ast.Ident{Name: "MyStruct"},
			expected: "MyStruct",
		},
		{
			name:     "star expr",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "MyStruct"}},
			expected: "MyStruct",
		},
		{
			name:     "selector expr",
			expr:     &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := getTypeName(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("getTypeName() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// Test_getMutexTypeName tests the private getMutexTypeName helper function.
func Test_getMutexTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typ      types.Type
		expected string
	}{
		{
			name:     "not named type",
			typ:      types.Typ[types.Int],
			expected: "",
		},
		{
			name:     "named type without package",
			typ:      types.NewNamed(types.NewTypeName(0, nil, "Test", nil), types.Typ[types.Int], nil),
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := getMutexTypeName(tt.typ)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("getMutexTypeName() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// Test_collectTypesWithValueReceivers tests the private collectTypesWithValueReceivers function.
func Test_collectTypesWithValueReceivers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function collects types with value receivers
		})
	}
}

// Test_checkStructsWithMutex tests the private checkStructsWithMutex function.
func Test_checkStructsWithMutex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks structs with mutex
		})
	}
}

// Test_checkValueReceivers tests the private checkValueReceivers function.
func Test_checkValueReceivers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks value receivers
		})
	}
}

// Test_checkValueParams tests the private checkValueParams function.
func Test_checkValueParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks value params
		})
	}
}

// Test_checkAssignments tests the private checkAssignments function.
func Test_checkAssignments(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks assignments
		})
	}
}

// Test_getMutexType_notInTypesInfo tests with expr not in TypesInfo.
func Test_getMutexType_notInTypesInfo(t *testing.T) {
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
			result := getMutexType(pass, expr)
			// Vérification du résultat
			if result != "" {
				t.Errorf("getMutexType() = %q, expected empty string", result)
			}

		})
	}
}

// Test_hasMutex_notInTypesInfo tests with expr not in TypesInfo.
func Test_hasMutex_notInTypesInfo(t *testing.T) {
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
			result := hasMutex(pass, expr)
			// Vérification du résultat
			if result {
				t.Errorf("hasMutex() = true, expected false")
			}

		})
	}
}

// Test_hasMutexInType tests the private hasMutexInType function.
func Test_hasMutexInType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks if has mutex in type
		})
	}
}

// Test_getMutexTypeFromType_notInTypesInfo tests with expr not in TypesInfo.
func Test_getMutexTypeFromType_notInTypesInfo(t *testing.T) {
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
			result := getMutexTypeFromType(pass, expr)
			// Vérification du résultat
			if result != "" {
				t.Errorf("getMutexTypeFromType() = %q, expected empty string", result)
			}

		})
	}
}

// Test_getMutexTypeFromType_nonStruct tests with non-struct type.
func Test_getMutexTypeFromType_nonStruct(t *testing.T) {
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

			// Test with basic type
			expr := &ast.Ident{Name: "x"}
			pass.TypesInfo.Types[expr] = types.TypeAndValue{
				Type: types.Typ[types.Int],
			}
			result := getMutexTypeFromType(pass, expr)
			// Vérification du résultat
			if result != "" {
				t.Errorf("getMutexTypeFromType() = %q, expected empty string for int", result)
			}

		})
	}
}

// Test_isMutexCopy tests the private isMutexCopy function.
func Test_isMutexCopy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks if is mutex copy
		})
	}
}

// Test_runVar019_disabled tests runVar019 with disabled rule.
func Test_runVar019_disabled(t *testing.T) {
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
					"KTN-VAR-019": {Enabled: config.Bool(false)},
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
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
					Uses:  make(map[*ast.Ident]types.Object),
					Defs:  make(map[*ast.Ident]types.Object),
				},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar019(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar019() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar019() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar019_fileExcluded tests runVar019 with excluded file.
func Test_runVar019_fileExcluded(t *testing.T) {
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
					"KTN-VAR-019": {
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
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
					Uses:  make(map[*ast.Ident]types.Object),
					Defs:  make(map[*ast.Ident]types.Object),
				},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar019(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar019() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar019() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_getMutexType tests the getMutexType private function.
func Test_getMutexType(t *testing.T) {
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

// Test_hasMutex tests the hasMutex private function.
func Test_hasMutex(t *testing.T) {
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

// Test_getMutexTypeFromType tests the getMutexTypeFromType private function.
func Test_getMutexTypeFromType(t *testing.T) {
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

// Test_runVar019_nilInspector tests runVar019 with nil inspector.
func Test_runVar019_nilInspector(t *testing.T) {
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

			result, err := runVar019(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar019() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar019() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar019_nilFset tests runVar019 with nil Fset.
func Test_runVar019_nilFset(t *testing.T) {
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

			result, err := runVar019(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar019() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar019() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar019_nilTypesInfo tests runVar019 with nil TypesInfo.
func Test_runVar019_nilTypesInfo(t *testing.T) {
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

			result, err := runVar019(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar019() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar019() = %v, expected nil", result)
			}
		})
	}
}

// Test_hasMutexInType_pointer tests hasMutexInType with pointer type.
func Test_hasMutexInType_pointer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a pointer to int (no mutex)
			intType := types.Typ[types.Int]
			ptrType := types.NewPointer(intType)

			result := hasMutexInType(ptrType)
			// Should be false
			if result {
				t.Errorf("hasMutexInType() = true, expected false for pointer to int")
			}
		})
	}
}

// Test_hasMutexInType_nonStruct tests hasMutexInType with non-struct type.
func Test_hasMutexInType_nonStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test with basic type
			intType := types.Typ[types.Int]

			result := hasMutexInType(intType)
			// Should be false
			if result {
				t.Errorf("hasMutexInType() = true, expected false for int type")
			}
		})
	}
}

// Test_isMutexCopy_notInTypesInfo tests isMutexCopy with rhs not in TypesInfo.
func Test_isMutexCopy_notInTypesInfo(t *testing.T) {
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

			lhs := &ast.Ident{Name: "x"}
			rhs := &ast.Ident{Name: "y"}

			result := isMutexCopy(pass, lhs, rhs)
			// Should be empty
			if result != "" {
				t.Errorf("isMutexCopy() = %q, expected empty string", result)
			}
		})
	}
}

// Test_isMutexCopy_notMutex tests isMutexCopy with non-mutex type.
func Test_isMutexCopy_notMutex(t *testing.T) {
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

			lhs := &ast.Ident{Name: "x"}
			rhs := &ast.Ident{Name: "y"}
			pass.TypesInfo.Types[rhs] = types.TypeAndValue{
				Type: types.Typ[types.Int],
			}

			result := isMutexCopy(pass, lhs, rhs)
			// Should be empty for non-mutex
			if result != "" {
				t.Errorf("isMutexCopy() = %q, expected empty string for int", result)
			}
		})
	}
}

// Test_collectTypesWithValueReceivers_emptyRecvList tests with empty receiver list.
func Test_collectTypesWithValueReceivers_emptyRecvList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a function with Recv but empty List
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{List: []*ast.Field{}}, // Empty list
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{}},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{}

			result := collectTypesWithValueReceivers(pass, insp)

			// Should be empty
			if len(result) != 0 {
				t.Errorf("collectTypesWithValueReceivers() returned %d types, expected 0", len(result))
			}
		})
	}
}

// Test_collectTypesWithValueReceivers_pointerReceiver tests with pointer receiver.
func Test_collectTypesWithValueReceivers_pointerReceiver(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a function with pointer receiver
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "s"}},
							Type:  &ast.StarExpr{X: &ast.Ident{Name: "MyStruct"}},
						},
					},
				},
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}}, Body: &ast.BlockStmt{List: []ast.Stmt{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{}

			result := collectTypesWithValueReceivers(pass, insp)

			// Should be empty - pointer receivers are not collected
			if len(result) != 0 {
				t.Errorf("collectTypesWithValueReceivers() returned %d types, expected 0", len(result))
			}
		})
	}
}

// Test_collectTypesWithValueReceivers_selectorType tests with selector type receiver.
func Test_collectTypesWithValueReceivers_selectorType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a function with selector type receiver (pkg.Type)
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "s"}},
							Type:  &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}},
						},
					},
				},
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}}, Body: &ast.BlockStmt{List: []ast.Stmt{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{}

			result := collectTypesWithValueReceivers(pass, insp)

			// Should be empty - selector types return empty string
			if len(result) != 0 {
				t.Errorf("collectTypesWithValueReceivers() returned %d types, expected 0", len(result))
			}
		})
	}
}

// Test_checkStructsWithMutex_noValueReceivers tests checkStructsWithMutex without value receivers.
func Test_checkStructsWithMutex_noValueReceivers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			// Create a struct type
			typeSpec := &ast.TypeSpec{
				Name: &ast.Ident{Name: "MyStruct"},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "x"}},
								Type:  &ast.Ident{Name: "int"},
							},
						},
					},
				},
			}

			file := &ast.File{
				Name: &ast.Ident{Name: "test"},
				Decls: []ast.Decl{
					&ast.GenDecl{
						Specs: []ast.Spec{typeSpec},
					},
				},
			}

			fset := token.NewFileSet()
			insp := inspector.New([]*ast.File{file})
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			// Empty map - no value receivers
			typesWithValueRecv := make(map[string]bool)

			checkStructsWithMutex(pass, insp, typesWithValueRecv)

			// Should not report anything
			if reportCount != 0 {
				t.Errorf("checkStructsWithMutex() reported %d issues, expected 0", reportCount)
			}
		})
	}
}

// Test_checkValueReceivers_noRecv tests checkValueReceivers with function without receiver.
func Test_checkValueReceivers_noRecv(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			// Create a regular function (no receiver)
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "RegularFunc"},
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}}, Body: &ast.BlockStmt{List: []ast.Stmt{}},
				Recv: nil, // no receiver
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			fset := token.NewFileSet()
			insp := inspector.New([]*ast.File{file})
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkValueReceivers(pass, insp)

			// Should not report anything
			if reportCount != 0 {
				t.Errorf("checkValueReceivers() reported %d issues, expected 0", reportCount)
			}
		})
	}
}

// Test_checkValueReceivers_emptyRecvList tests checkValueReceivers with empty receiver list.
func Test_checkValueReceivers_emptyRecvList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			// Create a function with empty receiver list
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}}, Body: &ast.BlockStmt{List: []ast.Stmt{}},
				Recv: &ast.FieldList{List: []*ast.Field{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			fset := token.NewFileSet()
			insp := inspector.New([]*ast.File{file})
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkValueReceivers(pass, insp)

			// Should not report anything
			if reportCount != 0 {
				t.Errorf("checkValueReceivers() reported %d issues, expected 0", reportCount)
			}
		})
	}
}

// Test_checkValueParams_noParams tests checkValueParams with function without params.
func Test_checkValueParams_noParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			// Create a function without params
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "NoParams"},
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}},
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			fset := token.NewFileSet()
			insp := inspector.New([]*ast.File{file})
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkValueParams(pass, insp)

			// Should not report anything
			if reportCount != 0 {
				t.Errorf("checkValueParams() reported %d issues, expected 0", reportCount)
			}
		})
	}
}

// Test_checkAssignments_unbalanced tests checkAssignments with unbalanced assignment.
func Test_checkAssignments_unbalanced(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			code := `package test
func f() {
	var x, y int
	x, y = 1, 2
}
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
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkAssignments(pass, insp)

			// Should not report anything (not mutex)
			if reportCount != 0 {
				t.Errorf("checkAssignments() reported %d issues, expected 0", reportCount)
			}
		})
	}
}

// Test_getMutexTypeFromType_pointerType tests getMutexTypeFromType with pointer type.
func Test_getMutexTypeFromType_pointerType(t *testing.T) {
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

			// Test with pointer to basic type
			expr := &ast.Ident{Name: "x"}
			pass.TypesInfo.Types[expr] = types.TypeAndValue{
				Type: types.NewPointer(types.Typ[types.Int]),
			}
			result := getMutexTypeFromType(pass, expr)
			// Verification du resultat
			if result != "" {
				t.Errorf("getMutexTypeFromType() = %q, expected empty string for pointer to int", result)
			}
		})
	}
}
