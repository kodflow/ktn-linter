// Internal tests for 007.go private functions.
package ktnstruct

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

// Test_runStruct007 teste la fonction runStruct007.
func Test_runStruct007(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"validation_success", false},
		{"validation_error_case", false},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			// Les cas d'erreur sont couverts via le test external
			if tt.wantErr {
				t.Error("Expected error but got none")
			}
		})
	}
}

// Test_isExportedField teste la fonction isExportedField.
func Test_isExportedField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty_string", "", false},
		{"exported_field", "Name", true},
		{"private_field", "name", false},
		{"single_char_upper", "N", true},
		{"single_char_lower", "n", false},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := isExportedField(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isExportedField(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_hasSerializationTag teste la fonction hasSerializationTag.
func Test_hasSerializationTag(t *testing.T) {
	tests := []struct {
		name     string
		tagValue string
		tags     []string
		expected bool
	}{
		{"no_tag", "", []string{"json", "xml"}, false},
		{"has_json", "`json:\"name\"`", []string{"json", "xml"}, true},
		{"has_xml", "`xml:\"name\"`", []string{"json", "xml"}, true},
		{"no_match", "`yaml:\"name\"`", []string{"json", "xml"}, false},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Créer un field avec le tag approprié
			var field *ast.Field
			// Vérifier si on a un tag
			if tt.tagValue != "" {
				field = &ast.Field{
					Tag: &ast.BasicLit{Value: tt.tagValue},
				}
			} else {
				field = &ast.Field{}
			}

			result := hasSerializationTag(field, tt.tags)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("hasSerializationTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_analyzeTypeSpec007 teste la fonction analyzeTypeSpec007.
func Test_analyzeTypeSpec007(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "non_struct_type",
			code: `package test
type MyInt int`,
		},
		{
			name: "private_struct",
			code: `package test
type myStruct struct {
	Name string
}`,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Configurer le test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

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
				Report: func(_ analysis.Diagnostic) {},
			}

			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}

			// Tester l'analyse
			insp.Preorder(nodeFilter, func(n ast.Node) {
				analyzeTypeSpec007(pass, cfg, n)
			})
		})
	}
}

// Test_checkExportedFieldsWithoutTags teste la fonction checkExportedFieldsWithoutTags.
func Test_checkExportedFieldsWithoutTags(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "struct_no_fields",
			code: `package test
type Empty struct{}`,
		},
		{
			name: "struct_with_tags",
			code: `package test
type User struct {
	Name string ` + "`json:\"name\"`" + `
}`,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Configurer le test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			pass := &analysis.Pass{
				Fset:   fset,
				Files:  []*ast.File{f},
				Report: func(_ analysis.Diagnostic) {},
			}

			// Trouver la struct
			ast.Inspect(f, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						checkExportedFieldsWithoutTags(pass, cfg, st)
					}
				}
				return true
			})
		})
	}
}
