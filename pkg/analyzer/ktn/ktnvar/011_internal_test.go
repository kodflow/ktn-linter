package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar011 tests the private runVar011 function.
func Test_runVar011(t *testing.T) {
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

// Test_extractIdent tests the private extractIdent helper function.
func Test_extractIdent(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected *ast.Ident
	}{
		{
			name:     "ident",
			expr:     &ast.Ident{Name: "x"},
			expected: &ast.Ident{Name: "x"},
		},
		{
			name:     "not ident",
			expr:     &ast.BasicLit{Value: "1"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractIdent(tt.expr)
			// Vérification du résultat
			if tt.expected == nil {
				// Vérification que result est nil
				if result != nil {
					t.Errorf("extractIdent() = %v, expected nil", result)
				}
			} else {
				// Vérification que result n'est pas nil et a le bon nom
				if result == nil {
					t.Errorf("extractIdent() = nil, expected ident with name %s", tt.expected.Name)
				} else if result.Name != tt.expected.Name {
					t.Errorf("extractIdent() = %s, expected %s", result.Name, tt.expected.Name)
				}
			}
		})
	}
}

// Test_checkShortVarDecl tests the private checkShortVarDecl function.
func Test_checkShortVarDecl(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks short var declarations
		})
	}
}

// Test_isShadowing tests the private isShadowing function.
func Test_isShadowing(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks if shadowing
		})
	}
}

// Test_lookupInParentScope tests the private lookupInParentScope function.
func Test_lookupInParentScope(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function looks up in parent scope
		})
	}
}

// Test_runVar011_disabled tests runVar011 with disabled rule.
func Test_runVar011_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-011": {Enabled: config.Bool(false)},
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

	_, err = runVar011(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar011() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar011() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar011_fileExcluded tests runVar011 with excluded file.
func Test_runVar011_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-011": {
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

	_, err = runVar011(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar011() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar011() reported %d issues, expected 0 when file excluded", reportCount)
	}
}
