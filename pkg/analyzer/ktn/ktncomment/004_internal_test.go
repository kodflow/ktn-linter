// Internal tests for 004.go - variable comment analyzer.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runComment004 tests the runComment004 function configuration.
// The actual analyzer is tested via analysistest in 004_external_test.go.
//
// Params:
//   - t: testing context
func Test_runComment004(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment004 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer004 is properly configured
			if Analyzer004 == nil {
				t.Error("Analyzer004 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer004.Name != "ktncomment004" {
				t.Errorf("Analyzer004.Name = %q, want %q", Analyzer004.Name, "ktncomment004")
			}
			// Check analyzer requires inspect
			if len(Analyzer004.Requires) == 0 {
				t.Error("Analyzer004 should require inspect.Analyzer")
			}
		})
	}
}

// Test_runComment004_ruleDisabled tests behavior when rule is disabled.
//
// Params:
//   - t: testing context
func Test_runComment004_ruleDisabled(t *testing.T) {
	// Import config package for test
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-COMMENT-004": {Enabled: config.Bool(false)},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	code := `package test
var myVar = 1`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment004(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment004 failed: %v", err)
	}

	// Should report no errors when rule disabled
	if errorCount != 0 {
		t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
	}
}

// Test_runComment004_fileExcluded tests behavior when file is excluded.
//
// Params:
//   - t: testing context
func Test_runComment004_fileExcluded(t *testing.T) {
	// Import config package for test
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-COMMENT-004": {
				Enabled: config.Bool(true),
				Exclude: []string{"*.go"},
			},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	code := `package test
var myVar = 1`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment004(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment004 failed: %v", err)
	}

	// Should report no errors when file excluded
	if errorCount != 0 {
		t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
	}
}

// Test_checkFileDeclarations tests checkFileDeclarations function.
//
// Params:
//   - t: testing context
func Test_checkFileDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		{
			name: "var with comment",
			code: `package test
// myVar is a variable
var myVar = 1`,
			wantErrs: 0,
		},
		{
			name: "var without comment",
			code: `package test
var myVar = 1`,
			wantErrs: 1,
		},
		{
			name: "const declaration skipped",
			code: `package test
const myConst = 1`,
			wantErrs: 0,
		},
		{
			name: "function declaration skipped",
			code: `package test
func myFunc() {}`,
			wantErrs: 0,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			errorCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					errorCount++
				},
			}

			// Run function
			checkFileDeclarations(pass, file)

			// Check error count
			if errorCount != tt.wantErrs {
				t.Errorf("checkFileDeclarations() errors = %d, want %d", errorCount, tt.wantErrs)
			}
		})
	}
}
