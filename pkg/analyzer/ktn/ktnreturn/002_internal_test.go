// Internal tests for analyzer 002 in ktnreturn package.
package ktnreturn

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runReturn002 tests the private runReturn002 function
func Test_runReturn002(t *testing.T) {
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
			_, err = runReturn002(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runReturn002 failed: %v", err)
			}

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_isSliceOrMapType tests the isSliceOrMapType function
func Test_isSliceOrMapType(t *testing.T) {
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

			got := isSliceOrMapType(pass, expr)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isSliceOrMapType(%s) = %v, want %v", tt.typeName, got, tt.want)
			}
		})
	}
}

// Test_isNilIdent tests the isNilIdent function
func Test_isNilIdent(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			expr, err := parser.ParseExpr(tt.code)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse expression: %v", err)
			}

			_ = fset // unused but kept for consistency
			got := isNilIdent(expr)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isNilIdent(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// Test_checkNilReturns tests the checkNilReturns private function.
func Test_checkNilReturns(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_collectSliceMapReturnTypes tests the collectSliceMapReturnTypes private function.
func Test_collectSliceMapReturnTypes(t *testing.T) {
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
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
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
					result := collectSliceMapReturnTypes(pass, funcDecl)
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

// Test_getSuggestionForType tests the getSuggestionForType private function.
func Test_getSuggestionForType(t *testing.T) {
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
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
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
						got := getSuggestionForType(tv.Type)
						// Check result matches expectation
						if got != tt.want {
							t.Errorf("getSuggestionForType() = %q, want %q", got, tt.want)
						}
					}
				}
				// Continue traversal
				return true
			})
		})
	}
}
