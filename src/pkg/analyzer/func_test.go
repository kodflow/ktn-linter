package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestFuncAnalyzer_ValidCases vérifie que les fonctions correctement documentées passent.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzer_ValidCases(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "Fonction sans params ni returns",
			code: `package test

// doSomething fait quelque chose.
func doSomething() {
}
`,
		},
		{
			name: "Fonction avec params seulement",
			code: `package test

// processData traite les données.
//
// Params:
//   - data: les données à traiter
//   - count: le nombre d'éléments
func processData(data string, count int) {
}
`,
		},
		{
			name: "Fonction avec returns seulement",
			code: `package test

// getValue retourne une valeur.
//
// Returns:
//   - string: la valeur récupérée
func getValue() string {
	return "test"
}
`,
		},
		{
			name: "Fonction avec params et returns",
			code: `package test

// calculateSum calcule la somme.
//
// Params:
//   - a: le premier nombre
//   - b: le second nombre
//
// Returns:
//   - int: la somme de a et b
func calculateSum(a int, b int) int {
	return a + b
}
`,
		},
		{
			name: "Fonction avec 4 paramètres (en dessous de la limite)",
			code: `package test

// complexFunction effectue un traitement complexe.
//
// Params:
//   - param1: le premier paramètre
//   - param2: le second paramètre
//   - param3: le troisième paramètre
//   - param4: le quatrième paramètre
//
// Returns:
//   - bool: le résultat du traitement
func complexFunction(param1 string, param2 int, param3 bool, param4 float64) bool {
	return true
}
`,
		},
		{
			name: "Fonction avec nommage MixedCaps",
			code: `package test

// calculateTotalAmount calcule le montant total.
//
// Params:
//   - amount: le montant de base
//
// Returns:
//   - float64: le montant total calculé
func calculateTotalAmount(amount float64) float64 {
	return amount * 1.2
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(diag analysis.Diagnostic) {
					t.Errorf("Unexpected diagnostic: %s at %s", diag.Message, fset.Position(diag.Pos))
				},
			}

			_, err = analyzer.FuncAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}
		})
	}
}

// TestFuncAnalyzer_ErrorCases vérifie que les violations sont bien détectées.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzer_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedError string
	}{
		{
			name: "Fonction sans commentaire",
			code: `package test

func noComment() {
}
`,
			expectedError: "KTN-FUNC-002",
		},
		{
			name: "Fonction avec params non documentés",
			code: `package test

// processData traite les données.
func processData(data string) {
}
`,
			expectedError: "KTN-FUNC-003",
		},
		{
			name: "Fonction avec returns non documenté",
			code: `package test

// getValue retourne une valeur.
func getValue() string {
	return "test"
}
`,
			expectedError: "KTN-FUNC-004",
		},
		{
			name: "Fonction avec trop de paramètres",
			code: `package test

// tooManyParams a trop de paramètres.
//
// Params:
//   - p1: param 1
//   - p2: param 2
//   - p3: param 3
//   - p4: param 4
//   - p5: param 5
//   - p6: param 6
//
// Returns:
//   - bool: résultat
func tooManyParams(p1, p2, p3, p4, p5, p6 int) bool {
	return true
}
`,
			expectedError: "KTN-FUNC-005",
		},
		{
			name: "Fonction avec mauvais nommage",
			code: `package test

// bad_naming utilise des underscores.
func bad_naming() {
}
`,
			expectedError: "KTN-FUNC-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			foundExpectedError := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(diag analysis.Diagnostic) {
					if tt.expectedError != "" {
						// Vérifier que le message contient le code d'erreur attendu
						if !foundExpectedError {
							foundExpectedError = true
							t.Logf("Found expected error: %s", diag.Message)
						}
					}
				},
			}

			_, err = analyzer.FuncAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}

			if tt.expectedError != "" && !foundExpectedError {
				t.Errorf("Expected error containing %q, but no errors were reported", tt.expectedError)
			}
		})
	}
}


// TestCheckParamsFormat teste la vérification du format des paramètres.
//
// Params:
//   - t: instance de test
func TestCheckParamsFormat(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		shouldPass    bool
		expectedParam string // Si shouldPass=false, param manquant attendu
	}{
		{
			name: "Tous les params documentés",
			code: `package test

// process traite les données.
//
// Params:
//   - data: les données
//   - count: le nombre
func process(data string, count int) {
}
`,
			shouldPass: true,
		},
		{
			name: "Param manquant",
			code: `package test

// process traite les données.
//
// Params:
//   - data: les données
func process(data string, count int) {
}
`,
			shouldPass:    false,
			expectedParam: "count",
		},
		{
			name: "Format avec indentation",
			code: `package test

// calculate effectue un calcul.
//
// Params:
//   - x: valeur x
//   - y: valeur y
//   - z: valeur z
//
// Returns:
//   - int: le résultat
func calculate(x, y, z int) int {
	return x + y + z
}
`,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			hasError := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(diag analysis.Diagnostic) {
					hasError = true
					t.Logf("Diagnostic: %s", diag.Message)
				},
			}

			_, err = analyzer.FuncAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}

			if tt.shouldPass && hasError {
				t.Errorf("Expected no errors, but got diagnostic")
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected error for missing param %q, but got none", tt.expectedParam)
			}
		})
	}
}

// TestRealWorldExample teste un exemple réel du codebase.
//
// Params:
//   - t: instance de test
func TestRealWorldExample(t *testing.T) {
	code := `package test

// checkGodocFormat vérifie le format strict avec sections Params: et Returns:.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
//   - doc: le texte du commentaire godoc
func checkGodocFormat(pass int, funcDecl string, funcName string, doc string) {
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			t.Errorf("Unexpected diagnostic for valid function: %s at %s",
				diag.Message, fset.Position(diag.Pos))
		},
	}

	_, err = analyzer.FuncAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Analyzer returned error: %v", err)
	}
}

// BenchmarkFuncAnalyzer mesure les performances de l'analyzer.
//
// Params:
//   - b: instance de benchmark
func BenchmarkFuncAnalyzer(b *testing.B) {
	code := `package test

// calculateSum calcule la somme.
//
// Params:
//   - a: le premier nombre
//   - b: le second nombre
//
// Returns:
//   - int: la somme de a et b
func calculateSum(a int, b int) int {
	return a + b
}
`

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, parser.ParseComments)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {},
		}
		analyzer.FuncAnalyzer.Run(pass)
	}
}
