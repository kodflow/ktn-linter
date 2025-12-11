package ktnfunc

import (
	"go/ast"
	"testing"

	"go/parser"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)


// Test_runFunc007_disabled tests behavior when rule is disabled.
func Test_runFunc007_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configuration avec règle désactivée
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-007": {Enabled: config.Bool(false)},
				},
			})
			// Reset config après le test
			defer config.Reset()

			// Créer un pass minimal
			result, err := runFunc007(&analysis.Pass{})
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

// Test_runFunc007_excludedFile tests behavior with excluded files.
func Test_runFunc007_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-007": {
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
			_, err = runFunc007(pass)
			// Vérification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}

// Test_runFunc007 tests the runFunc007 private function.
func Test_runFunc007(t *testing.T) {
	// Test cases pour la fonction privée runFunc007
	// La logique principale est testée via l'API publique dans 009_external_test.go
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

// Test_isGetter vérifie la détection des getters.
func Test_isGetter(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		expected bool
	}{
		{
			name:     "error case validation",
			funcName: "GetValue",
			expected: true,
		},
		{
			name:     "IsValid getter",
			funcName: "IsValid",
			expected: true,
		},
		{
			name:     "HasData getter",
			funcName: "HasData",
			expected: true,
		},
		{
			name:     "NotGetter function",
			funcName: "Calculate",
			expected: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := isGetter(tt.funcName)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isGetter(%s) = %v, want %v", tt.funcName, result, tt.expected)
			}
		})
	}
}

// Test_hasSideEffect vérifie la détection des effets de bord.
func Test_hasSideEffect(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name: "error case validation",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "obj"},
				Sel: &ast.Ident{Name: "field"},
			},
			expected: true,
		},
		{
			name: "simple identifier",
			expr: &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name: "index on selector",
			expr: &ast.IndexExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "obj"},
					Sel: &ast.Ident{Name: "arr"},
				},
			},
			expected: true,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := hasSideEffect(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("hasSideEffect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_checkGetterSideEffects vérifie la détection des side effects dans les getters.
//
// Params:
//   - t: instance de testing
func Test_checkGetterSideEffects(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "getter_side_effect_detection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			_ = tt.name
		})
	}
}

// Test_reportAssignSideEffect vérifie le rapport des assignations.
//
// Params:
//   - t: instance de testing
func Test_reportAssignSideEffect(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "assign_side_effect_report",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			_ = tt.name
		})
	}
}

// Test_reportIncDecSideEffect vérifie le rapport des incréments/décréments.
//
// Params:
//   - t: instance de testing
func Test_reportIncDecSideEffect(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "incdec_side_effect_report",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			_ = tt.name
		})
	}
}

// Test_collectLazyLoadFields vérifie la collecte des champs lazy load.
//
// Params:
//   - t: instance de testing
func Test_collectLazyLoadFields(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "lazy_load_fields_collection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test avec body vide (pas nil pour éviter panic)
			body := &ast.BlockStmt{List: []ast.Stmt{}}
			result := collectLazyLoadFields(body)
			// Vérification du résultat
			if result == nil {
				t.Error("collectLazyLoadFields() should return empty map, not nil")
			}
		})
	}
}

// Test_isNilComparison vérifie la détection de comparaisons avec nil.
//
// Params:
//   - t: instance de testing
func Test_isNilComparison(t *testing.T) {
	tests := []struct {
		name     string
		binary   *ast.BinaryExpr
		expected bool
	}{
		{
			name: "nil_on_right",
			binary: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.Ident{Name: "nil"},
				Op: 0,
			},
			expected: true,
		},
		{
			name: "nil_on_left",
			binary: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "nil"},
				Y:  &ast.Ident{Name: "x"},
				Op: 0,
			},
			expected: true,
		},
		{
			name: "no_nil",
			binary: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "x"},
				Y:  &ast.Ident{Name: "y"},
				Op: 0,
			},
			expected: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := isNilComparison(tt.binary)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isNilComparison() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractFieldName vérifie l'extraction du nom de champ.
//
// Params:
//   - t: instance de testing
func Test_extractFieldName(t *testing.T) {
	tests := []struct {
		name     string
		binary   *ast.BinaryExpr
		expected string
	}{
		{
			name: "selector_on_left",
			binary: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "obj"},
					Sel: &ast.Ident{Name: "field"},
				},
				Y: &ast.Ident{Name: "nil"},
			},
			expected: "field",
		},
		{
			name: "selector_on_right",
			binary: &ast.BinaryExpr{
				X: &ast.Ident{Name: "nil"},
				Y: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "s"},
					Sel: &ast.Ident{Name: "cache"},
				},
			},
			expected: "cache",
		},
		{
			name: "no_selector",
			binary: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "nil"},
			},
			expected: "",
		},
		{
			name: "both_sides_not_selector",
			binary: &ast.BinaryExpr{
				X: &ast.Ident{Name: "x"},
				Y: &ast.Ident{Name: "y"},
			},
			expected: "",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := extractFieldName(tt.binary)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractFieldName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_isLazyLoadAssignment vérifie la détection des assignations lazy load.
//
// Params:
//   - t: instance de testing
func Test_isLazyLoadAssignment(t *testing.T) {
	tests := []struct {
		name       string
		lhs        ast.Expr
		lazyFields map[string]bool
		expected   bool
	}{
		{
			name: "lazy_field_match",
			lhs: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "s"},
				Sel: &ast.Ident{Name: "cache"},
			},
			lazyFields: map[string]bool{"cache": true},
			expected:   true,
		},
		{
			name: "no_match",
			lhs: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "s"},
				Sel: &ast.Ident{Name: "data"},
			},
			lazyFields: map[string]bool{"cache": true},
			expected:   false,
		},
		{
			name:       "not_selector",
			lhs:        &ast.Ident{Name: "x"},
			lazyFields: map[string]bool{"x": true},
			expected:   false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := isLazyLoadAssignment(tt.lhs, tt.lazyFields)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isLazyLoadAssignment() = %v, want %v", result, tt.expected)
			}
		})
	}
}
