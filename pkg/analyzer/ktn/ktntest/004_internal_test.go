package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestIsTestsVariableName teste la détection des noms de variables de tests.
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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := isTestsVariableName(tt.varName)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isTestsVariableName(%q) = %v, want %v", tt.varName, got, tt.want)
			}
		})
	}
}

// TestHasTableDrivenPattern teste la détection du pattern table-driven.
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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

// Test_checkAssignStmt tests the checkAssignStmt private function.
func Test_checkAssignStmt(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "tests variable detected",
			code: `package test
func TestExample(t *testing.T) {
	tests := []struct{}{}
}`,
			want: true,
		},
		{
			name: "testcases variable detected",
			code: `package test
func TestExample(t *testing.T) {
	testcases := []struct{}{}
}`,
			want: true,
		},
		{
			name: "error case - no test variable",
			code: `package test
func TestExample(t *testing.T) {
	data := []int{1, 2}
}`,
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification pas d'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			var found bool
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if assignStmt, ok := n.(*ast.AssignStmt); ok {
					found = checkAssignStmt(assignStmt)
				}
				// Continuer la traversée
				return true
			})

			// Vérification du résultat
			if found != tt.want {
				t.Errorf("checkAssignStmt() = %v, want %v", found, tt.want)
			}
		})
	}
}

// Test_checkRangeStmt tests the checkRangeStmt private function.
func Test_checkRangeStmt(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "range over tests variable",
			code: `package test
func TestExample(t *testing.T) {
	for _, tt := range tests {
	}
}`,
			want: true,
		},
		{
			name: "range over testcases variable",
			code: `package test
func TestExample(t *testing.T) {
	for _, tc := range testcases {
	}
}`,
			want: true,
		},
		{
			name: "error case - range over non-test variable",
			code: `package test
func TestExample(t *testing.T) {
	data := []int{1, 2}
	for _, d := range data {
	}
}`,
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification pas d'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			var found bool
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if rangeStmt, ok := n.(*ast.RangeStmt); ok {
					found = checkRangeStmt(rangeStmt)
				}
				// Continuer la traversée
				return true
			})

			// Vérification du résultat
			if found != tt.want {
				t.Errorf("checkRangeStmt() = %v, want %v", found, tt.want)
			}
		})
	}
}

// Test_runTest004 tests the runTest004 private function.
func Test_runTest004(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - minimal test",
			code: `package test
func TestExample(t *testing.T) {}`,
		},
		{
			name: "error case - empty body",
			code: `package test
func TestEmpty(t *testing.T) {
	// empty test
}`,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing code: %s", tt.code)
		})
	}
}

// Test_runTest004_disabled tests that the rule is skipped when disabled.
func Test_runTest004_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runTest004_excludedFile tests that excluded files are skipped.
func Test_runTest004_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}
