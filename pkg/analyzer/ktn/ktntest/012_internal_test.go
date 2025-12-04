// Internal tests for analyzer 012.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestIsPassthroughTest tests the isPassthroughTest function.
//
// Params:
//   - t: testing context
func TestIsPassthroughTest(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "empty body",
			code: `package test
func TestEmpty(t *testing.T) {}`,
			expected: true,
		},
		{
			name: "with assertion",
			code: `package test
func TestWithAssertion(t *testing.T) {
	t.Error("fail")
}`,
			expected: false,
		},
		{
			name: "no assertion",
			code: `package test
func TestNoAssertion(t *testing.T) {
	x := 1 + 1
	_ = x
}`,
			expected: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver la fonction de test
			var funcDecl *ast.FuncDecl
			// Parcourir les déclarations
			for _, decl := range file.Decls {
				// Vérifier si c'est une FuncDecl
				if fd, ok := decl.(*ast.FuncDecl); ok {
					funcDecl = fd
					// Sortir de la boucle
					break
				}
			}

			// Vérification fonction trouvée
			if funcDecl == nil {
				t.Fatal("no function found")
			}

			result := isPassthroughTest(funcDecl)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("isPassthroughTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckForAssertion tests the checkForAssertion function.
//
// Params:
//   - t: testing context
func TestCheckForAssertion(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "t.Error call",
			code:     `package test; func f(t *testing.T) { t.Error("x") }`,
			expected: true,
		},
		{
			name:     "t.Fatal call",
			code:     `package test; func f(t *testing.T) { t.Fatal("x") }`,
			expected: true,
		},
		{
			name:     "assert.Equal call",
			code:     `package test; func f(t *testing.T) { assert.Equal(t, 1, 1) }`,
			expected: true,
		},
		{
			name:     "require.NoError call",
			code:     `package test; func f(t *testing.T) { require.NoError(t, nil) }`,
			expected: true,
		},
		{
			name:     "no assertion",
			code:     `package test; func f(t *testing.T) { x := 1 }`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver le premier appel
			hasAssertion := false
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si assertion trouvée
				if checkForAssertion(n) {
					hasAssertion = true
					// Arrêter la traversée
					return false
				}
				// Continuer
				return true
			})

			// Vérification résultat
			if hasAssertion != tt.expected {
				t.Errorf("checkForAssertion() = %v, want %v", hasAssertion, tt.expected)
			}
		})
	}
}

// TestIsTestingMethodCall tests the isTestingMethodCall function.
//
// Params:
//   - t: testing context
func TestIsTestingMethodCall(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "t.Error",
			code:     `package test; func f() { t.Error("x") }`,
			expected: true,
		},
		{
			name:     "t.Errorf",
			code:     `package test; func f() { t.Errorf("x") }`,
			expected: true,
		},
		{
			name:     "t.Fatal",
			code:     `package test; func f() { t.Fatal("x") }`,
			expected: true,
		},
		{
			name:     "t.Fatalf",
			code:     `package test; func f() { t.Fatalf("x") }`,
			expected: true,
		},
		{
			name:     "t.Log",
			code:     `package test; func f() { t.Log("x") }`,
			expected: true,
		},
		{
			name:     "not a testing method",
			code:     `package test; func f() { fmt.Println("x") }`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver le premier appel
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un appel
				if callExpr, ok := n.(*ast.CallExpr); ok {
					found = isTestingMethodCall(callExpr)
					// Arrêter la traversée
					return false
				}
				// Continuer
				return true
			})

			// Vérification résultat
			if found != tt.expected {
				t.Errorf("isTestingMethodCall() = %v, want %v", found, tt.expected)
			}
		})
	}
}

// TestIsAssertLibraryCall tests the isAssertLibraryCall function.
//
// Params:
//   - t: testing context
func TestIsAssertLibraryCall(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "assert.Equal",
			code:     `package test; func f() { assert.Equal(t, 1, 1) }`,
			expected: true,
		},
		{
			name:     "require.NoError",
			code:     `package test; func f() { require.NoError(t, nil) }`,
			expected: true,
		},
		{
			name:     "not assert library",
			code:     `package test; func f() { fmt.Println("x") }`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver le premier appel
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un appel
				if callExpr, ok := n.(*ast.CallExpr); ok {
					found = isAssertLibraryCall(callExpr)
					// Arrêter la traversée
					return false
				}
				// Continuer
				return true
			})

			// Vérification résultat
			if found != tt.expected {
				t.Errorf("isAssertLibraryCall() = %v, want %v", found, tt.expected)
			}
		})
	}
}

// Test_runTest012 tests the runTest012 function indirectly via the analyzer.
//
// Params:
//   - t: testing context
func Test_runTest012(t *testing.T) {
	tests := []struct {
		name       string
		checkField string
		expected   any
	}{
		{"analyzer exists", "notNil", Analyzer012 != nil},
		{"run is defined", "notNil", Analyzer012.Run != nil},
		{"name is correct", "name", Analyzer012.Name},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification selon le type de check
			switch tt.checkField {
			// Check non-nil
			case "notNil":
				// Vérification booléenne
				if tt.expected != true {
					t.Errorf("expected true, got %v", tt.expected)
				}
			// Check name
			case "name":
				// Vérification du nom
				if tt.expected != "ktntest012" {
					t.Errorf("expected ktntest012, got %v", tt.expected)
				}
			}
		})
	}
}

// Test_isTestHelperCall tests the isTestHelperCall function.
//
// Params:
//   - t: testing context
func Test_isTestHelperCall(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "helper with t as first arg",
			code:     `package test; func f() { helper(t, 1, 2) }`,
			expected: true,
		},
		{
			name:     "function without t",
			code:     `package test; func f() { doSomething(1, 2) }`,
			expected: false,
		},
		{
			name:     "no arguments",
			code:     `package test; func f() { noArgs() }`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver le premier appel
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un appel
				if callExpr, ok := n.(*ast.CallExpr); ok {
					found = isTestHelperCall(callExpr)
					// Arrêter la traversée
					return false
				}
				// Continuer
				return true
			})

			// Vérification résultat
			if found != tt.expected {
				t.Errorf("isTestHelperCall() = %v, want %v", found, tt.expected)
			}
		})
	}
}
