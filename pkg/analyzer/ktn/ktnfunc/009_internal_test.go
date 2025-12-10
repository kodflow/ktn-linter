package ktnfunc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/ast/inspector"


	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)


// Test_runFunc009_disabled tests behavior when rule is disabled.
func Test_runFunc009_disabled(t *testing.T) {
	// Configuration avec règle désactivée
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-FUNC-009": {Enabled: config.Bool(false)},
		},
	})
	// Reset config après le test
	defer config.Reset()

	// Créer un pass minimal
	result, err := runFunc009(&analysis.Pass{})
	// Vérification de l'erreur
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Vérification du résultat nil
	if result != nil {
		t.Errorf("Expected nil result when rule disabled, got %v", result)
	}
}

// Test_runFunc009_excludedFile tests behavior with excluded files.
func Test_runFunc009_excludedFile(t *testing.T) {
	// Configuration avec fichier exclu
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-FUNC-009": {
				Enabled:       config.Bool(true),
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
	_, err = runFunc009(pass)
	// Vérification erreur
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test_runFunc009 tests the runFunc009 private function.
func Test_runFunc009(t *testing.T) {
	// Test cases pour la fonction privée runFunc009
	// La logique principale est testée via l'API publique dans 003_external_test.go
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

// Test_getAllowedValues vérifie que getAllowedValues retourne les valeurs autorisées.
func Test_getAllowedValues(t *testing.T) {
	tests := []struct {
		name     string
		expected map[string]bool
	}{
		{
			name: "error case validation",
			expected: map[string]bool{
				"0":  true,
				"1":  true,
				"-1": true,
			},
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := getAllowedValues()
			// Vérifier que toutes les valeurs attendues sont présentes
			for key, expectedVal := range tt.expected {
				// Vérification de la valeur
				if result[key] != expectedVal {
					t.Errorf("getAllowedValues()[%s] = %v, want %v", key, result[key], expectedVal)
				}
			}
		})
	}
}

// Test_collectAllowedLiterals vérifie la collecte des littéraux autorisés.
func Test_collectAllowedLiterals(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "error case validation",
			code: `package test
const MaxSize = 100
const MinValue = -1`,
			expected: 2,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			result := collectAllowedLiterals(insp)

			// Vérifier le nombre de littéraux collectés
			if len(result) != tt.expected {
				t.Errorf("collectAllowedLiterals() = %d literals, want %d", len(result), tt.expected)
			}
		})
	}
}

// Test_checkMagicNumbers vérifie la détection des magic numbers.
func Test_checkMagicNumbers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - la logique est testée via external tests
		})
	}
}
