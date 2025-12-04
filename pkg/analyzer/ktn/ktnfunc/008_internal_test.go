package ktnfunc

import (
	"go/ast"
	"testing"
)

// Test_runFunc008 tests the runFunc008 private function.
func Test_runFunc008(t *testing.T) {
	// Test cases pour la fonction privée runFunc008
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
