// Internal tests for analyzer 013 in ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runFunc013 tests the private runFunc013 function
func Test_runFunc013(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		{
			name: "nil return for slice",
			code: `package test
func GetSlice() []int {
	return nil
}`,
			wantErrs: 1,
		},
		{
			name: "nil return for map",
			code: `package test
func GetMap() map[string]int {
	return nil
}`,
			wantErrs: 1,
		},
		{
			name: "empty slice return",
			code: `package test
func GetSlice() []int {
	return []int{}
}`,
			wantErrs: 0,
		},
		{
			name: "empty map return",
			code: `package test
func GetMap() map[string]int {
	return map[string]int{}
}`,
			wantErrs: 0,
		},
		{
			name: "nil return for error is OK",
			code: `package test
func GetError() error {
	return nil
}`,
			wantErrs: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Check type checking success
			if err != nil {
				t.Logf("type check error (may be expected): %v", err)
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
				ResultOf:  make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runFunc013(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runFunc013 failed: %v", err)
			}

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_isSliceOrMapTypeFunc013 tests the isSliceOrMapTypeFunc013 function
func Test_isSliceOrMapTypeFunc013(t *testing.T) {
	code := `package test
type MySlice []int
type MyMap map[string]int
type MyStruct struct{}
type MyInt int
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	// Create type checker
	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check("test", fset, []*ast.File{file}, info)
	// Check type checking success
	if err != nil {
		t.Logf("type check error: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: info,
	}

	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		{"slice type", "MySlice", true},
		{"map type", "MyMap", true},
		{"struct type", "MyStruct", false},
		{"int type", "MyInt", false},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var expr ast.Expr
			// Find type expression
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for type spec
				if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == tt.typeName {
					expr = ts.Type
					return false
				}
				return true
			})

			// Check expression was found
			if expr == nil {
				t.Fatalf("could not find type %s", tt.typeName)
			}

			got := isSliceOrMapTypeFunc013(pass, expr)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isSliceOrMapTypeFunc013(%s) = %v, want %v", tt.typeName, got, tt.want)
			}
		})
	}
}

// Test_isNilIdentFunc013 tests the isNilIdentFunc013 function
func Test_isNilIdentFunc013(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"nil identifier", "nil", true},
		{"other identifier", "x", false},
		{"number literal", "42", false},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			expr, err := parser.ParseExpr(tt.code)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse expression: %v", err)
			}

			_ = fset // unused but kept for consistency
			got := isNilIdentFunc013(expr)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isNilIdentFunc013(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// Test_checkNilReturnsFunc013 tests the checkNilReturnsFunc013 private function.
func Test_checkNilReturnsFunc013(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		{
			name: "function with no body (interface method)",
			code: `package test
type Iface interface {
	GetSlice() []int
}`,
			wantErrs: 0,
		},
		{
			name: "function with body and nil return",
			code: `package test
func GetSlice() []int {
	return nil
}`,
			wantErrs: 1,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Find function declaration and test
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					checkNilReturnsFunc013(pass, funcDecl)
				}
				// Continue traversal
				return true
			})

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_collectSliceMapReturnTypesFunc013 tests the collectSliceMapReturnTypesFunc013 private function.
func Test_collectSliceMapReturnTypesFunc013(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantLen  int
		wantVals []string
	}{
		{
			name: "function with slice return",
			code: `package test
func GetSlice() []int {
	return []int{}
}`,
			wantLen:  1,
			wantVals: []string{"[]int{}"},
		},
		{
			name: "function with map return",
			code: `package test
func GetMap() map[string]int {
	return map[string]int{}
}`,
			wantLen:  1,
			wantVals: []string{"map[string]int{}"},
		},
		{
			name: "function with no results",
			code: `package test
func NoReturn() {
	return
}`,
			wantLen:  0,
			wantVals: []string{},
		},
		{
			name: "function with multiple named returns",
			code: `package test
func MultipleNamed() (a, b []int) {
	return []int{}, []int{}
}`,
			wantLen:  2,
			wantVals: []string{"[]int{}", "[]int{}"},
		},
		{
			name: "function with unnamed return",
			code: `package test
func UnnamedReturn() []string {
	return []string{}
}`,
			wantLen:  1,
			wantVals: []string{"[]string{}"},
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Find function declaration
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					result := collectSliceMapReturnTypesFunc013(pass, funcDecl)
					// Verify length
					if len(result) != tt.wantLen {
						t.Errorf("expected %d return types, got %d", tt.wantLen, len(result))
					}
					// Verify values match
					for i, want := range tt.wantVals {
						if i < len(result) && result[i] != want {
							t.Logf("return type %d: got %q, want %q", i, result[i], want)
						}
					}
				}
				// Continue traversal
				return true
			})
		})
	}
}

// Test_getSuggestionForTypeFunc013 tests the getSuggestionForTypeFunc013 private function.
func Test_getSuggestionForTypeFunc013(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "slice type suggestion",
			code: `package test
type MySlice []int`,
			want: "[]int{}",
		},
		{
			name: "map type suggestion",
			code: `package test
type MyMap map[string]int`,
			want: "map[string]int{}",
		},
		{
			name: "non-slice/map type returns empty string",
			code: `package test
type MyInt int`,
			want: "",
		},
		{
			name: "struct type returns empty string",
			code: `package test
type MyStruct struct{}`,
			want: "",
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Defs:  make(map[*ast.Ident]types.Object),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			// Find the type and test suggestion
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for type spec
				if ts, ok := n.(*ast.TypeSpec); ok {
					if tv, ok := info.Types[ts.Type]; ok {
						got := getSuggestionForTypeFunc013(tv.Type)
						// Check result matches expectation
						if got != tt.want {
							t.Errorf("getSuggestionForTypeFunc013() = %q, want %q", got, tt.want)
						}
					}
				}
				// Continue traversal
				return true
			})
		})
	}
}

// Test_isSliceOrMapTypeFunc013WithNilTypeInfo tests isSliceOrMapTypeFunc013 with nil type info.
func Test_isSliceOrMapTypeFunc013WithNilTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", `package test
			type MySlice []int`, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create empty type info (no types recorded)
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Find type spec and test with missing type info
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for type spec
				if ts, ok := n.(*ast.TypeSpec); ok {
					got := isSliceOrMapTypeFunc013(pass, ts.Type)
					// Should return false when type info is nil
					if got != false {
						t.Errorf("isSliceOrMapTypeFunc013 with nil type info = %v, want false", got)
					}
				}
				// Continue traversal
				return true
			})

		})
	}
}

// Test_runFunc013WithDisabledRule tests runFunc013 when rule is disabled.
func Test_runFunc013WithDisabledRule(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Set configuration to disable rule
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					ruleCodeFunc013: {Enabled: config.Bool(false)},
				},
			})
			// Reset config after test
			defer config.Reset()

			code := `package test
			func GetSlice() []int {
			return nil
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Check type checking
			if err != nil {
				t.Logf("type check error: %v", err)
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
				ResultOf:  make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runFunc013(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runFunc013 failed: %v", err)
			}

			// Should report 0 errors when rule is disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runFunc013WithFileExclusion tests runFunc013 with excluded files.
func Test_runFunc013WithFileExclusion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Set configuration to exclude specific file
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					ruleCodeFunc013: {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config after test
			defer config.Reset()

			code := `package test
			func GetSlice() []int {
			return nil
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Check type checking
			if err != nil {
				t.Logf("type check error: %v", err)
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
				ResultOf:  make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runFunc013(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runFunc013 failed: %v", err)
			}

			// Should report 0 errors when file is excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}

// Test_collectSliceMapReturnTypesFunc013WithNilTypeInfo tests edge cases.
func Test_collectSliceMapReturnTypesFunc013WithNilTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func GetData() []int {
			return []int{}
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create empty type info (simulates missing type information)
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Find function declaration and test
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					result := collectSliceMapReturnTypesFunc013(pass, funcDecl)
					// Should return empty string when type info is missing
					if len(result) != 1 || result[0] != "" {
						t.Errorf("expected [\"\"], got %v", result)
					}
				}
				// Continue traversal
				return true
			})

		})
	}
}

// Test_checkNilReturnsFunc013WithEmptyTypeInfo tests checkNilReturnsFunc013 with empty type info.
func Test_checkNilReturnsFunc013WithEmptyTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func GetData() []int {
			return nil
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create empty type info (simulates missing type information)
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Find function declaration and test checkNilReturnsFunc013
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					checkNilReturnsFunc013(pass, funcDecl)
				}
				// Continue traversal
				return true
			})

			// Should report 0 errors when type info is empty
			// because typeInfo will be "" and won't trigger the report
			if errorCount != 0 {
				t.Errorf("expected 0 errors when type info is empty, got %d", errorCount)
			}

		})
	}
}

// Test_checkNilReturnsFunc013NonReturnStatement tests that non-return statements are skipped.
func Test_checkNilReturnsFunc013NonReturnStatement(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func GetData() []int {
			x := 1
			if x > 0 {
				x = 2
			}
			return []int{x}
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Create type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{file},
				TypesInfo: info,
			}

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Find function declaration and test checkNilReturnsFunc013
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					checkNilReturnsFunc013(pass, funcDecl)
				}
				// Continue traversal
				return true
			})

			// Should report 0 errors for valid code
			if errorCount != 0 {
				t.Errorf("expected 0 errors, got %d", errorCount)
			}

		})
	}
}
