package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestIsTestingMethod teste la détection des méthodes du package testing
func TestIsTestingMethod(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		want       bool
	}{
		// Méthodes testing valides
		{name: "Error", methodName: "Error", want: true},
		{name: "Errorf", methodName: "Errorf", want: true},
		{name: "Fatal", methodName: "Fatal", want: true},
		{name: "Fatalf", methodName: "Fatalf", want: true},
		{name: "Fail", methodName: "Fail", want: true},
		{name: "FailNow", methodName: "FailNow", want: true},
		{name: "Log", methodName: "Log", want: true},
		{name: "Logf", methodName: "Logf", want: true},
		{name: "Skip", methodName: "Skip", want: true},
		{name: "Skipf", methodName: "Skipf", want: true},
		{name: "SkipNow", methodName: "SkipNow", want: true},
		// Méthodes non-testing
		{name: "Parallel", methodName: "Parallel", want: false},
		{name: "Run", methodName: "Run", want: false},
		{name: "Cleanup", methodName: "Cleanup", want: false},
		{name: "Helper", methodName: "Helper", want: false},
		{name: "Name", methodName: "Name", want: false},
		{name: "Random", methodName: "Random", want: false},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTestingMethod(tt.methodName)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isTestingMethod(%q) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}

// TestIsAssertMethod teste la détection des méthodes testify/assert
func TestIsAssertMethod(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		want       bool
	}{
		// Méthodes assert valides
		{name: "Equal", methodName: "Equal", want: true},
		{name: "NotEqual", methodName: "NotEqual", want: true},
		{name: "Nil", methodName: "Nil", want: true},
		{name: "NotNil", methodName: "NotNil", want: true},
		{name: "True", methodName: "True", want: true},
		{name: "False", methodName: "False", want: true},
		{name: "Empty", methodName: "Empty", want: true},
		{name: "NotEmpty", methodName: "NotEmpty", want: true},
		{name: "Len", methodName: "Len", want: true},
		{name: "Contains", methodName: "Contains", want: true},
		{name: "NotContains", methodName: "NotContains", want: true},
		{name: "Greater", methodName: "Greater", want: true},
		{name: "GreaterOrEqual", methodName: "GreaterOrEqual", want: true},
		{name: "Less", methodName: "Less", want: true},
		{name: "LessOrEqual", methodName: "LessOrEqual", want: true},
		{name: "NoError", methodName: "NoError", want: true},
		{name: "Error", methodName: "Error", want: true},
		{name: "ErrorIs", methodName: "ErrorIs", want: true},
		{name: "ErrorAs", methodName: "ErrorAs", want: true},
		{name: "ErrorContains", methodName: "ErrorContains", want: true},
		{name: "Panics", methodName: "Panics", want: true},
		{name: "NotPanics", methodName: "NotPanics", want: true},
		{name: "JSONEq", methodName: "JSONEq", want: true},
		{name: "YAMLEq", methodName: "YAMLEq", want: true},
		{name: "Regexp", methodName: "Regexp", want: true},
		{name: "NotRegexp", methodName: "NotRegexp", want: true},
		{name: "Zero", methodName: "Zero", want: true},
		{name: "NotZero", methodName: "NotZero", want: true},
		// Méthodes non-assert
		{name: "Random", methodName: "Random", want: false},
		{name: "NewAssertionFunc", methodName: "NewAssertionFunc", want: false},
		{name: "ObjectsAreEqual", methodName: "ObjectsAreEqual", want: false},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAssertMethod(tt.methodName)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isAssertMethod(%q) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}

// TestIsRequireMethod teste la détection des méthodes testify/require
func TestIsRequireMethod(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		want       bool
	}{
		// require utilise les mêmes méthodes que assert
		{name: "NoError", methodName: "NoError", want: true},
		{name: "Equal", methodName: "Equal", want: true},
		{name: "NotNil", methodName: "NotNil", want: true},
		{name: "True", methodName: "True", want: true},
		{name: "Random", methodName: "Random", want: false},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRequireMethod(tt.methodName)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isRequireMethod(%q) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}

// TestIsAssertionCall teste la détection des appels d'assertion
func TestIsAssertionCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		// Cas testing package
		{
			name: "t.Error",
			code: `package test
func TestExample(t *testing.T) {
	t.Error("message")
}`,
			want: true,
		},
		{
			name: "t.Fatal",
			code: `package test
func TestExample(t *testing.T) {
	t.Fatal("message")
}`,
			want: true,
		},
		{
			name: "t.Errorf",
			code: `package test
func TestExample(t *testing.T) {
	t.Errorf("message %s", "test")
}`,
			want: true,
		},
		// Cas testify/assert
		{
			name: "assert.Equal",
			code: `package test
func TestExample(t *testing.T) {
	assert.Equal(t, 1, 1)
}`,
			want: true,
		},
		{
			name: "assert.True",
			code: `package test
func TestExample(t *testing.T) {
	assert.True(t, true)
}`,
			want: true,
		},
		{
			name: "assert.NoError",
			code: `package test
func TestExample(t *testing.T) {
	assert.NoError(t, nil)
}`,
			want: true,
		},
		// Cas testify/require
		{
			name: "require.NoError",
			code: `package test
func TestExample(t *testing.T) {
	require.NoError(t, nil)
}`,
			want: true,
		},
		{
			name: "require.Equal",
			code: `package test
func TestExample(t *testing.T) {
	require.Equal(t, 1, 1)
}`,
			want: true,
		},
		// Cas négatifs
		{
			name: "t.Run (pas une assertion)",
			code: `package test
func TestExample(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {})
}`,
			want: false,
		},
		{
			name: "t.Parallel (pas une assertion)",
			code: `package test
func TestExample(t *testing.T) {
	t.Parallel()
}`,
			want: false,
		},
		{
			name: "autre fonction",
			code: `package test
func TestExample(t *testing.T) {
	fmt.Println("test")
}`,
			want: false,
		},
		{
			name: "appel sur sélecteur (pas ident)",
			code: `package test
func TestExample(t *testing.T) {
	myObj.assert.Equal(1, 1)
}`,
			want: false,
		},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parser le code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification pas d'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Trouver le premier appel de fonction
			var foundCall *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification si c'est un appel
				if call, ok := n.(*ast.CallExpr); ok {
					foundCall = call
					// Arrêter la recherche
					return false
				}
				// Continuer la recherche
				return true
			})

			// Vérifier qu'on a trouvé un appel
			if foundCall == nil {
				t.Fatal("no function call found in test code")
			}

			// Tester isAssertionCall
			got := isAssertionCall(foundCall)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isAssertionCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsTestsVariableName teste la détection des noms de variables de tests
func TestIsTestsVariableName(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		want    bool
	}{
		// Noms valides
		{name: "tests", varName: "tests", want: true},
		{name: "Tests", varName: "Tests", want: true},
		{name: "TESTS", varName: "TESTS", want: true},
		{name: "testcases", varName: "testcases", want: true},
		{name: "TestCases", varName: "TestCases", want: true},
		{name: "cases", varName: "cases", want: true},
		{name: "Cases", varName: "Cases", want: true},
		// Noms invalides
		{name: "test", varName: "test", want: false},
		{name: "testData", varName: "testData", want: false},
		{name: "tc", varName: "tc", want: false},
		{name: "data", varName: "data", want: false},
		{name: "items", varName: "items", want: false},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTestsVariableName(tt.varName)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isTestsVariableName(%q) = %v, want %v", tt.varName, got, tt.want)
			}
		})
	}
}

// TestHasTableDrivenPattern teste la détection du pattern table-driven
func TestHasTableDrivenPattern(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "table-driven avec tests",
			code: `package test
func TestExample(t *testing.T) {
	tests := []struct {
		name string
		input int
		want int
	}{
		{"case1", 1, 2},
		{"case2", 2, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}`,
			want: true,
		},
		{
			name: "table-driven avec testcases",
			code: `package test
func TestExample(t *testing.T) {
	testcases := []struct {
		input int
		want int
	}{
		{1, 2},
		{2, 4},
	}
	for _, tc := range testcases {
		// test
	}
}`,
			want: true,
		},
		{
			name: "table-driven avec cases",
			code: `package test
func TestExample(t *testing.T) {
	cases := []struct {
		input int
		want int
	}{
		{1, 2},
	}
	for _, c := range cases {
		// test
	}
}`,
			want: true,
		},
		{
			name: "pas de table-driven (variable différente)",
			code: `package test
func TestExample(t *testing.T) {
	data := []int{1, 2, 3}
	for _, d := range data {
		// test
	}
}`,
			want: false,
		},
		{
			name: "pas de table-driven (pas de range)",
			code: `package test
func TestExample(t *testing.T) {
	tests := []struct {
		input int
		want int
	}{
		{1, 2},
	}
	// pas de boucle range
}`,
			want: false,
		},
		{
			name: "pas de table-driven (pas de variable tests)",
			code: `package test
func TestExample(t *testing.T) {
	for i := 0; i < 10; i++ {
		// test
	}
}`,
			want: false,
		},
		{
			name: "pas de table-driven (range sur fonction, pas ident)",
			code: `package test
func TestExample(t *testing.T) {
	for _, item := range getItems() {
		// test
	}
}`,
			want: false,
		},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parser le code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification pas d'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Trouver la fonction TestExample
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification si c'est une fonction
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "TestExample" {
					funcDecl = fn
					// Arrêter la recherche
					return false
				}
				// Continuer la recherche
				return true
			})

			// Vérifier qu'on a trouvé la fonction
			if funcDecl == nil {
				t.Fatal("TestExample function not found")
			}

			// Tester hasTableDrivenPattern
			got := hasTableDrivenPattern(funcDecl)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("hasTableDrivenPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasMultipleAssertions teste le comptage des assertions
func TestHasMultipleAssertions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "3 assertions t.Error (devrait déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	t.Error("error 1")
	t.Error("error 2")
	t.Error("error 3")
}`,
			want: true,
		},
		{
			name: "2 assertions t.Error (ne devrait pas déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	t.Error("error 1")
	t.Error("error 2")
}`,
			want: false,
		},
		{
			name: "3 assertions assert.Equal (devrait déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	assert.Equal(t, 1, 1)
	assert.Equal(t, 2, 2)
	assert.Equal(t, 3, 3)
}`,
			want: true,
		},
		{
			name: "3 assertions require.NoError (devrait déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	require.NoError(t, nil)
	require.NoError(t, nil)
	require.NoError(t, nil)
}`,
			want: true,
		},
		{
			name: "mix assertions (3+, devrait déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	t.Error("error")
	assert.True(t, true)
	require.NoError(t, nil)
}`,
			want: true,
		},
		{
			name: "1 assertion (ne devrait pas déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	t.Error("error")
}`,
			want: false,
		},
		{
			name: "pas d'assertions (ne devrait pas déclencher)",
			code: `package test
func TestExample(t *testing.T) {
	result := Calculate(1, 2)
	fmt.Println(result)
}`,
			want: false,
		},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parser le code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification pas d'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Trouver la fonction TestExample
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification si c'est une fonction
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "TestExample" {
					funcDecl = fn
					// Arrêter la recherche
					return false
				}
				// Continuer la recherche
				return true
			})

			// Vérifier qu'on a trouvé la fonction
			if funcDecl == nil {
				t.Fatal("TestExample function not found")
			}

			// Tester hasMultipleAssertions
			got := hasMultipleAssertions(funcDecl)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("hasMultipleAssertions() = %v, want %v", got, tt.want)
			}
		})
	}
}
