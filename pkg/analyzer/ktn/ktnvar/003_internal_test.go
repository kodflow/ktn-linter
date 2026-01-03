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

// Test_runVar003 tests the runVar003 function.
func Test_runVar003(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		ruleEnabled bool
		expectCount int
	}{
		{
			name:        "enabled with violation",
			code:        `package test; var my_var int = 42`,
			ruleEnabled: true,
			expectCount: 1,
		},
		{
			name:        "enabled without violation",
			code:        `package test; var myVar int = 42`,
			ruleEnabled: true,
			expectCount: 0,
		},
		{
			name:        "disabled",
			code:        `package test; var my_var int = 42`,
			ruleEnabled: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.ruleEnabled {
				config.Reset()
			} else {
				config.Set(&config.Config{
					Rules: map[string]*config.RuleConfig{
						"KTN-VAR-003": {Enabled: config.Bool(false)},
					},
				})
			}
			defer config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset:     fset,
				ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
				Report:   func(_ analysis.Diagnostic) { reportCount++ },
			}

			_, _ = runVar003(pass)

			if reportCount != tt.expectCount {
				t.Errorf("runVar003() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar003Names tests the checkVar003Names function.
func Test_checkVar003Names(t *testing.T) {
	tests := []struct {
		name        string
		varName     string
		expectCount int
	}{
		{
			name:        "blank identifier",
			varName:     "_",
			expectCount: 0,
		},
		{
			name:        "snake_case",
			varName:     "my_var",
			expectCount: 1,
		},
		{
			name:        "camelCase",
			varName:     "myVar",
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			defer config.Reset()

			fset := token.NewFileSet()
			reportCount := 0

			pass := &analysis.Pass{
				Fset:   fset,
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			valueSpec := &ast.ValueSpec{
				Names: []*ast.Ident{{Name: tt.varName}},
			}

			checkVar003Names(pass, valueSpec)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar003Names() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// TestHasUnderscore003 tests the hasUnderscore003 function.
func TestHasUnderscore003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "blank identifier",
			input:    "_",
			expected: false,
		},
		{
			name:     "camelCase",
			input:    "myVariable",
			expected: false,
		},
		{
			name:     "PascalCase",
			input:    "MyVariable",
			expected: false,
		},
		{
			name:     "snake_case",
			input:    "my_variable",
			expected: true,
		},
		{
			name:     "SCREAMING_SNAKE_CASE",
			input:    "MY_VARIABLE",
			expected: true,
		},
		{
			name:     "mixed_Case",
			input:    "My_Variable",
			expected: true,
		},
		{
			name:     "single underscore prefix",
			input:    "_private",
			expected: true,
		},
		{
			name:     "acronym",
			input:    "HTTPStatus",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := hasUnderscore003(tt.input)
			if result != tt.expected {
				t.Errorf("hasUnderscore003(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_runVar003_disabled tests runVar003 with disabled rule.
func Test_runVar003_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-003": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, err = runVar003(pass)
	if err != nil {
		t.Fatalf("runVar003() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar003_fileExcluded tests runVar003 with excluded file.
func Test_runVar003_fileExcluded(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-003": {Exclude: []string{"test.go"}},
		},
	})
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, err = runVar003(pass)
	if err != nil {
		t.Fatalf("runVar003() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 when file excluded", reportCount)
	}
}

// Test_runVar003_snakeCase tests runVar003 detecting snake_case.
func Test_runVar003_snakeCase(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should report one issue
	if reportCount != 1 {
		t.Errorf("runVar003() reported %d issues, expected 1", reportCount)
	}
}

// Test_checkVar003Names_blankIdent tests checkVar003Names with blank identifier.
func Test_checkVar003Names_blankIdent(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
var _ int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should not report for blank identifier
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 for blank identifier", reportCount)
	}
}

// Test_runVar003_nonVarDecl tests skipping non-var declarations.
func Test_runVar003_nonVarDecl(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Code with const (not var)
	code := `package test
const MY_CONST = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should not report for const declarations
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 for const", reportCount)
	}
}

// Test_runVar003_typeDecl tests skipping type declarations.
func Test_runVar003_typeDecl(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Code with type (not var)
	code := `package test
type My_Type int
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should not report for type declarations (VAR-003 checks vars only)
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 for type", reportCount)
	}
}
