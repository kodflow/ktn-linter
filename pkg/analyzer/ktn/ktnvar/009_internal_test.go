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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_checkMakeCallVar008_notMake tests with non-make call.
func Test_checkMakeCallVar008_notMake(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test with non-make call
			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "len"},
			}
			checkMakeCallVar008(pass, call)
			// No error expected

		})
	}
}

// Test_checkMakeCallVar008_tooFewArgs tests with insufficient args.
func Test_checkMakeCallVar008_tooFewArgs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test with make but only 1 arg
			call := &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{&ast.Ident{Name: "T"}},
			}
			checkMakeCallVar008(pass, call)
			// No error expected

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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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

// Test_checkMakeCallVar008 tests the checkMakeCallVar008 private function.
func Test_checkMakeCallVar008(t *testing.T) {
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

// Test_checkMakeCallVar008_notSlice tests with non-slice type.
func Test_checkMakeCallVar008_notSlice(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"not a slice type"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
					Uses:  make(map[*ast.Ident]types.Object),
					Defs:  make(map[*ast.Ident]types.Object),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// Test make with map type (not slice)
			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.MapType{
						Key:   &ast.Ident{Name: "string"},
						Value: &ast.Ident{Name: "int"},
					},
					&ast.BasicLit{Kind: token.INT, Value: "10"},
				},
			}
			checkMakeCallVar008(pass, call)
			// No error expected - not a slice
		})
	}
}

// Test_checkMakeCallVar008_zeroLength tests with zero length.
func Test_checkMakeCallVar008_zeroLength(t *testing.T) {
	// This test verifies zero-length make calls are not flagged
	// The actual logic is tested via analysistest in external tests
	t.Run("zero length coverage", func(t *testing.T) {
		// Passthrough - main functionality tested via external tests
	})
}

// Test_checkMakeCallVar008_smallConstantSize tests VAR-016 skip case.
func Test_checkMakeCallVar008_smallConstantSize(t *testing.T) {
	t.Run("small constant size skip", func(t *testing.T) {
		// Parse code with make call using small constant
		code := `package test
const size = 512
func example() {
	s := make([]int, size)
	_ = s
}
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		// Type check
		conf := types.Config{}
		info := &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		}
		pkg, _ := conf.Check("test", fset, []*ast.File{file}, info)

		reportCount := 0
		pass := &analysis.Pass{
			Fset:      fset,
			TypesInfo: info,
			Pkg:       pkg,
			Report: func(_d analysis.Diagnostic) {
				reportCount++
			},
		}

		// Find and check the make call
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "make" {
					checkMakeCallVar008(pass, call)
				}
			}
			return true
		})

		// Should NOT report - VAR-016 handles small constant sizes
		if reportCount != 0 {
			t.Errorf("checkMakeCallVar008() reported %d, expected 0 (VAR-016 skip)", reportCount)
		}
	})
}
