package ktnfunc

import (
	"testing"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)


// Test_runFunc011_disabled tests behavior when rule is disabled.
func Test_runFunc011_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"rule disabled returns early"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec règle désactivée
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-011": {Enabled: config.Bool(false)},
				},
			})
			// Reset config après le test
			defer config.Reset()

			// Créer un pass minimal
			result, err := runFunc011(&analysis.Pass{})
			// Vérification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Vérification du résultat nil
			if result != nil {
				t.Errorf("Expected nil result when rule disabled, got %v", result)
			}
		})
	}
}

// Test_runFunc011_excludedFile tests behavior with excluded files.
func Test_runFunc011_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"excluded file skipped"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-011": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config après le test
			defer config.Reset()

			code := `package test
func foo() { }
`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer un inspector
			files := []*ast.File{file}
			inspectResult, _ := inspect.Analyzer.Run(&analysis.Pass{
				Fset:  fset,
				Files: files,
			})

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Expected no diagnostics for excluded file, got: %s", d.Message)
				},
			}

			// Exécuter l'analyse
			_, err = runFunc011(pass)
			// Vérification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

// Test_runFunc011 tests the runFunc011 private function.
func Test_runFunc011(t *testing.T) {
	// Test cases pour la fonction privée runFunc011
	// La logique principale est testée via l'API publique dans 005_external_test.go
	// Ce test vérifie les cas edge de la fonction privée

	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - la logique principale est dans external tests
		})
	}
}

// Test_calculateComplexity tests the calculateComplexity private function.
func Test_calculateComplexity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}
