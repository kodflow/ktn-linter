package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// TestIsEmptyInterface tests the isEmptyInterface function.
func TestIsEmptyInterface(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		interfaceType *ast.InterfaceType
		expected      bool
	}{
		{
			name:          "nil methods",
			interfaceType: &ast.InterfaceType{Methods: nil},
			expected:      true,
		},
		{
			name: "empty methods list",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{List: nil},
			},
			expected: true,
		},
		{
			name: "empty methods slice",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{List: []*ast.Field{}},
			},
			expected: true,
		},
		{
			name: "non-empty interface",
			interfaceType: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Method"}},
						},
					},
				},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isEmptyInterface(tt.interfaceType)
			// Check result
			if result != tt.expected {
				t.Errorf("isEmptyInterface() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckEmptyInterface tests the checkEmptyInterface function.
func TestCheckEmptyInterface(t *testing.T) {
	// Test: non-empty interface
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: "Method"}},
				},
			},
		},
	}
	// Should not panic
	checkEmptyInterface(nil, interfaceType)
}

// Test_runVar024_ruleDisabled tests analyzer returns early when disabled.
func Test_runVar024_ruleDisabled(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar024]

	// Disable the rule
	cfg.Rules[ruleCodeVar024] = &config.RuleConfig{Enabled: config.Bool(false)}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar024)
		} else {
			cfg.Rules[ruleCodeVar024] = originalRule
		}
	}()

	// Run analyzer with testhelper
	diags := testhelper.RunAnalyzer(t, Analyzer024, "testdata/src/var024/bad.go")

	// With rule disabled, should have 0 errors
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when rule disabled, got %d", len(diags))
	}
}

// Test_runVar024_fileExcluded tests file exclusion.
func Test_runVar024_fileExcluded(t *testing.T) {
	// Test source with interface{}
	src := `package test

var x interface{} = nil
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded_file.go", src, parser.ParseComments)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a mock pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer to populate ResultOf
	inspResult, err := inspect.Analyzer.Run(pass)
	// Check for inspector errors
	if err != nil {
		t.Fatalf("Failed to run inspect: %v", err)
	}
	pass.ResultOf[inspect.Analyzer] = inspResult

	// Get config and add file exclusion
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar024]

	// Add the exclusion
	cfg.Rules[ruleCodeVar024] = &config.RuleConfig{Exclude: []string{"excluded_file.go"}}
	// Ensure cleanup
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar024)
		} else {
			cfg.Rules[ruleCodeVar024] = originalRule
		}
	}()

	// Run the actual analyzer
	_, runErr := runVar024(pass)
	// Check for analyzer errors
	if runErr != nil {
		t.Fatalf("runVar024() returned error: %v", runErr)
	}

	// Should not detect any issues (file excluded)
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for excluded file, got %d", len(diagnostics))
	}
}

// Test_checkEmptyInterface_report tests checkEmptyInterface reporting.
func Test_checkEmptyInterface_report(t *testing.T) {
	// Test source with interface{}
	src := `package test

var x interface{} = nil
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a mock pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer to populate ResultOf
	inspResult, err := inspect.Analyzer.Run(pass)
	// Check for inspector errors
	if err != nil {
		t.Fatalf("Failed to run inspect: %v", err)
	}
	pass.ResultOf[inspect.Analyzer] = inspResult

	// Run the analyzer
	_, runErr := runVar024(pass)
	// Check for analyzer errors
	if runErr != nil {
		t.Fatalf("runVar024() returned error: %v", runErr)
	}

	// Should detect the empty interface{}
	if len(diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diagnostics))
	}
}
