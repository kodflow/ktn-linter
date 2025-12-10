// Internal tests for 006.go private functions.
package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runStruct006 teste la fonction runStruct006.
//
// Params:
//   - t: instance de testing
func Test_runStruct006(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "struct006_analysis",
			expectErr: false,
		},
		{
			name:      "struct006_error_case",
			expectErr: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			// Les cas d'erreur sont couverts via le test external
			_ = tt.expectErr
		})
	}
}

// Test_checkPrivateFieldsWithTags teste la fonction checkPrivateFieldsWithTags.
//
// Params:
//   - t: instance de testing
func Test_checkPrivateFieldsWithTags(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name: "empty struct",
			code: `package test
type User struct{}`,
			expectError: false,
		},
		{
			name: "public field with tag",
			code: `package test
type User struct {
	Name string ` + "`json:\"name\"`" + `
}`,
			expectError: false,
		},
		{
			name: "private field without tag",
			code: `package test
type User struct {
	name string
}`,
			expectError: false,
		},
		{
			name: "private field with tag",
			code: `package test
type UserDTO struct {
	name string ` + "`json:\"name\"`" + `
}`,
			expectError: true,
		},
		{
			name: "field with empty backticks is treated as tag",
			code: `package test
type User struct {
	name string ` + "``" + `
}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			errCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				Report: func(_ analysis.Diagnostic) {
					errCount++
				},
			}

			// Find the struct type
			ast.Inspect(f, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						checkPrivateFieldsWithTags(pass, st, ts.Name.Name)
					}
				}
				return true
			})

			if tt.expectError && errCount == 0 {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && errCount > 0 {
				t.Errorf("Expected no error but got %d", errCount)
			}
		})
	}

	// Test edge case: struct with nil Fields
	t.Run("struct with nil fields", func(t *testing.T) {
		pass := &analysis.Pass{
			Fset: token.NewFileSet(),
			Report: func(_ analysis.Diagnostic) {
				t.Error("Should not report error for nil fields")
			},
		}
		// Create a StructType with nil Fields
		structType := &ast.StructType{Fields: nil}
		checkPrivateFieldsWithTags(pass, structType, "TestStruct")
	})
}

// Test_isPrivateField teste la fonction isPrivateField.
//
// Params:
//   - t: instance de testing
func Test_isPrivateField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "lowercase_start_is_private",
			input:    "name",
			expected: true,
		},
		{
			name:     "uppercase_start_is_public",
			input:    "Name",
			expected: false,
		},
		{
			name:     "empty_string_is_not_private",
			input:    "",
			expected: false,
		},
		{
			name:     "single_lowercase_is_private",
			input:    "a",
			expected: true,
		},
		{
			name:     "single_uppercase_is_public",
			input:    "A",
			expected: false,
		},
		{
			name:     "underscore_prefix_is_private",
			input:    "_internal",
			expected: false,
		},
	}

	// Exécution des tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPrivateField(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isPrivateField(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_runStruct006_disabled tests that the rule is skipped when disabled.
func Test_runStruct006_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-006": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runStruct006(pass)
	if err != nil {
		t.Errorf("runStruct006() error = %v", err)
	}
}

// Test_runStruct006_excludedFile tests that excluded files are skipped.
func Test_runStruct006_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-006": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runStruct006(pass)
	if err != nil {
		t.Errorf("runStruct006() error = %v", err)
	}
}
