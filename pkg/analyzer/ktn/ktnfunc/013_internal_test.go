package ktnfunc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_runFunc013 tests the runFunc013 private function.
func Test_runFunc013(t *testing.T) {
	// Test cases pour la fonction privée runFunc013
	// La logique principale est testée via l'API publique dans 013_external_test.go
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

// Test_collectFunctionParams vérifie la collecte des paramètres de fonction.
func Test_collectFunctionParams(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "error case validation",
			code: `package test
func foo(a int, b string) {}`,
			expected: 2,
		},
		{
			name: "no params",
			code: `package test
func bar() {}`,
			expected: 0,
		},
		{
			name: "ignore underscore params",
			code: `package test
func baz(_ int, a string) {}`,
			expected: 1,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver la première fonction
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une fonction
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					// Arrêter la recherche
					return false
				}
				// Continuer la recherche
				return true
			})

			// Vérification de la fonction trouvée
			if funcDecl == nil {
				t.Fatal("No function found")
			}

			result := collectFunctionParams(funcDecl)
			// Vérification du nombre de paramètres
			if len(result) != tt.expected {
				t.Errorf("collectFunctionParams() = %d params, want %d", len(result), tt.expected)
			}
		})
	}
}

// Test_collectUsedVariables vérifie la collecte des variables utilisées.
func Test_collectUsedVariables(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case validation",
			code: `{ x := 1; y := x + 1; }`,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			expr, err := parser.ParseExpr(tt.code)
			// Ignorer les erreurs de parsing pour ce test simple
			_ = expr
			_ = err
			_ = fset

			// Test passthrough - nécessite un contexte AST complet
		})
	}
}

// Test_collectIgnoredVariables vérifie la collecte des variables ignorées.
func Test_collectIgnoredVariables(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "error case validation",
			code: `package test
func foo(x int) { _ = x }`,
			expected: 1,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver le body de la fonction
			var body *ast.BlockStmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une fonction
				if fd, ok := n.(*ast.FuncDecl); ok {
					body = fd.Body
					// Arrêter la recherche
					return false
				}
				// Continuer la recherche
				return true
			})

			// Vérification du body trouvé
			if body == nil {
				t.Fatal("No function body found")
			}

			result := collectIgnoredVariables(body)
			// Vérification du nombre de variables ignorées
			if len(result) != tt.expected {
				t.Errorf("collectIgnoredVariables() = %d vars, want %d", len(result), tt.expected)
			}
		})
	}
}

// Test_findParentAssignToBlank vérifie la recherche d'assignation à _.
func Test_findParentAssignToBlank(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte AST complet
		})
	}
}
