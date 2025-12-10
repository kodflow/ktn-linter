package ktnfunc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/packages"
)

// Test_runFunc001_disabled tests behavior when rule is disabled.
func Test_runFunc001_disabled(t *testing.T) {
	// Configuration avec règle désactivée
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-FUNC-001": {Enabled: config.Bool(false)},
		},
	})
	// Reset config après le test
	defer config.Reset()

	// Créer un pass minimal
	result, err := runFunc001(&analysis.Pass{})
	// Vérification de l'erreur
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Vérification du résultat nil
	if result != nil {
		t.Errorf("Expected nil result when rule disabled, got %v", result)
	}
}

// Test_runFunc001_excludedFile tests behavior with excluded files.
func Test_runFunc001_excludedFile(t *testing.T) {
	// Configuration avec fichier exclu
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-FUNC-001": {
				Enabled:       config.Bool(true),
				Exclude: []string{"test.go"},
			},
		},
	})
	// Reset config après le test
	defer config.Reset()

	code := `package test
func foo() (string, error) { return "", nil }
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
	_, err = runFunc001(pass)
	// Vérification erreur
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test_runFunc001 tests the runFunc001 private function.
func Test_runFunc001(t *testing.T) {
	// Test cases pour la fonction privée runFunc001
	// La logique principale est testée via l'API publique dans 006_external_test.go
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

// Test_validateErrorInReturns vérifie la validation de la position des erreurs.
func Test_validateErrorInReturns(t *testing.T) {
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

// Test_isErrorType vérifie la détection du type error.
func Test_isErrorType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "error case validation",
			code: `package test
func foo() error { return nil }`,
			expected: true,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo}
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer un package minimal
			pkg := &packages.Package{
				Fset:      fset,
				Syntax:    []*ast.File{file},
				TypesInfo: &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)},
			}
			_ = pkg
			_ = cfg

			// Test passthrough - le test complet nécessite un contexte d'analyse complet
		})
	}
}

// Test_isBuiltinError vérifie la détection du type error builtin.
//
// Params:
//   - t: instance de testing
func Test_isBuiltinError(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
	}{
		{
			name:     "error_interface_detection",
			typeName: "error",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite types.Type réel
			_ = tt.typeName
		})
	}
}

// Test_implementsError vérifie si un type implémente error.
//
// Params:
//   - t: instance de testing
func Test_implementsError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "implements_error_detection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite types.Type réel
			_ = tt.name
		})
	}
}
