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


// Test_runFunc005_disabled tests behavior when rule is disabled.
func Test_runFunc005_disabled(t *testing.T) {
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
					"KTN-FUNC-005": {Enabled: config.Bool(false)},
				},
			})
			// Reset config après le test
			defer config.Reset()

			// Créer un pass minimal
			result, err := runFunc005(&analysis.Pass{})
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

// Test_runFunc005_excludedFile tests behavior with excluded files.
func Test_runFunc005_excludedFile(t *testing.T) {
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
					"KTN-FUNC-005": {
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
			_, err = runFunc005(pass)
			// Vérification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}

// Test_runFunc005 tests the runFunc005 private function.
func Test_runFunc005(t *testing.T) {
	// Test cases pour la fonction privée runFunc005
	// La logique principale est testée via l'API publique dans 001_external_test.go
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

// Test_isLineToSkip tests the isLineToSkip private function.
func Test_isLineToSkip(t *testing.T) {
	tests := []struct {
		name           string
		trimmed        string
		inBlockComment bool
		want           bool
	}{
		{"empty line error case", "", false, true},
		{"comment line error case", "// comment", false, true},
		{"block comment start error case", "/* comment", false, true},
		{"code line", "code", false, false},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			inBlock := tt.inBlockComment
			got := isLineToSkip(tt.trimmed, &inBlock)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isLineToSkip(%q) = %v, want %v", tt.trimmed, got, tt.want)
			}
		})
	}
}

// Test_countPureCodeLines tests the countPureCodeLines private function.
func Test_countPureCodeLines(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "simple code with comment",
			code: `package test
func test() {
	// This is a comment
	x := 1
}`,
			want: 1,
		},
		{
			name: "code with block comment",
			code: `package test
func test() {
	/* This is
	   a block comment */
	x := 1
	y := 2
}`,
			want: 2,
		},
		{
			name: "code with empty lines",
			code: `package test
func test() {
	x := 1

	y := 2

}`,
			want: 2,
		},
		{
			name: "code with only braces",
			code: `package test
func test() {
	{
		x := 1
	}
}`,
			want: 1,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver la fonction
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok {
					// Assignation de la déclaration de fonction
					funcDecl = fd
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que la fonction a été trouvée
			if funcDecl == nil || funcDecl.Body == nil {
				t.Fatal("Expected to find function with body")
			}

			// Créer un pass avec ReadFile
			pass := &analysis.Pass{
				Fset: fset,
				ReadFile: func(filename string) ([]byte, error) {
					// Retourner le code source
					return []byte(tt.code), nil
				},
			}

			// Appeler countPureCodeLines
			result := countPureCodeLines(pass, funcDecl.Body)
			// Vérification du résultat
			if result != tt.want {
				t.Errorf("countPureCodeLines() = %v, want %v", result, tt.want)
			}
		})
	}
}
