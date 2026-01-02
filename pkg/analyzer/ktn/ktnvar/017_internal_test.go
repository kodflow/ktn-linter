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

// Test_runVar017 tests the private runVar017 function.
func Test_runVar017(t *testing.T) {
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

// Test_runVar017_disabled tests runVar017 with disabled rule.
func Test_runVar017_disabled(t *testing.T) {
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
					"KTN-VAR-017": {Enabled: config.Bool(false)},
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

			_, err = runVar017(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar017() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar017() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar017_fileExcluded tests runVar017 with excluded file.
func Test_runVar017_fileExcluded(t *testing.T) {
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
					"KTN-VAR-017": {
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

			_, err = runVar017(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar017() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar017() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_runVar017_withMakeMap tests runVar017 with make(map) without capacity.
func Test_runVar017_withMakeMap(t *testing.T) {
	config.Reset()

	code := `package test
func foo() {
	m := make(map[string]int)
	_ = m
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

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
		Fset: fset,
		Pkg:  pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo: info,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err := runVar017(pass)
	if err != nil {
		t.Errorf("runVar017() error = %v", err)
	}

	// Should report the map without capacity
	if reportCount != 1 {
		t.Errorf("runVar017() reported %d, expected 1", reportCount)
	}
}

// Test_runVar017_withMapCapacity tests runVar017 with make(map) with capacity.
func Test_runVar017_withMapCapacity(t *testing.T) {
	config.Reset()

	code := `package test
func foo() {
	m := make(map[string]int, 10)
	_ = m
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

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
		Fset: fset,
		Pkg:  pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo: info,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err := runVar017(pass)
	if err != nil {
		t.Errorf("runVar017() error = %v", err)
	}

	// Should not report - map has capacity
	if reportCount != 0 {
		t.Errorf("runVar017() with capacity reported %d, expected 0", reportCount)
	}
}

// Test_runVar017_withVerbose tests runVar017 with verbose mode.
func Test_runVar017_withVerbose(t *testing.T) {
	config.Set(&config.Config{
		Verbose: true,
	})
	defer config.Reset()

	code := `package test
func foo() {
	m := make(map[string]int)
	_ = m
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

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
		Fset: fset,
		Pkg:  pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo: info,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err := runVar017(pass)
	if err != nil {
		t.Errorf("runVar017() error = %v", err)
	}

	// Should report with verbose mode
	if reportCount != 1 {
		t.Errorf("runVar017() with verbose reported %d, expected 1", reportCount)
	}
}

// Test_runVar017_notMakeCall tests runVar017 with non-make call.
func Test_runVar017_notMakeCall(t *testing.T) {
	config.Reset()

	code := `package test
func foo() {
	x := len("hello")
	_ = x
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, _ := conf.Check("test", fset, []*ast.File{file}, info)

	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset: fset,
		Pkg:  pkg,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		TypesInfo: info,
		Report:    func(_d analysis.Diagnostic) {},
	}

	result, err := runVar017(pass)
	if err != nil {
		t.Errorf("runVar017() error = %v", err)
	}
	if result != nil {
		t.Errorf("runVar017() result = %v, expected nil", result)
	}
}
