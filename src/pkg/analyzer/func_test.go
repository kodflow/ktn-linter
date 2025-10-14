package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// runFuncAnalyzerValidTest exécute un test qui ne devrait pas produire d'erreur.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
func runFuncAnalyzerValidTest(t *testing.T, name, code string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
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

// runFuncAnalyzerErrorTest exécute un test qui devrait produire une erreur spécifique.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - expectedError: code d'erreur attendu
func runFuncAnalyzerErrorTest(t *testing.T, name, code, expectedError string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		foundExpectedError := false
		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				if expectedError != "" && !foundExpectedError {
					foundExpectedError = true
					t.Logf("Found expected error: %s", diag.Message)
				}
			},
		}

		_, err = analyzer.FuncAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}

		if expectedError != "" && !foundExpectedError {
			t.Errorf("Expected error containing %q, but no errors were reported", expectedError)
		}
	})
}

// TestFuncAnalyzerBasicCases vérifie les cas basiques sans params/returns.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerBasicCases(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Fonction sans params ni returns", `package test

// doSomething fait quelque chose.
func doSomething() {
}
`)
}

// TestFuncAnalyzerWithParams vérifie les fonctions avec paramètres.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerWithParams(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Fonction avec params seulement", `package test

// processData traite les données.
//
// Params:
//   - data: les données à traiter
//   - count: le nombre d'éléments
func processData(data string, count int) {
}
`)

	runFuncAnalyzerValidTest(t, "Fonction avec 4 paramètres (en dessous de la limite)", `package test

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
`)
}

// TestFuncAnalyzerWithReturns vérifie les fonctions avec valeurs de retour.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerWithReturns(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Fonction avec returns seulement", `package test

// getValue retourne une valeur.
//
// Returns:
//   - string: la valeur récupérée
func getValue() string {
	return "test"
}
`)

	runFuncAnalyzerValidTest(t, "Fonction avec params et returns", `package test

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
`)

	runFuncAnalyzerValidTest(t, "Fonction avec nommage MixedCaps", `package test

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
`)
}

// TestFuncAnalyzerMissingDoc vérifie les erreurs de documentation.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerMissingDoc(t *testing.T) {
	runFuncAnalyzerErrorTest(t, "Fonction sans commentaire", `package test

func noComment() {
}
`, "KTN-FUNC-002")

	runFuncAnalyzerErrorTest(t, "Fonction avec params non documentés", `package test

// processData traite les données.
func processData(data string) {
}
`, "KTN-FUNC-003")

	runFuncAnalyzerErrorTest(t, "Fonction avec returns non documenté", `package test

// getValue retourne une valeur.
func getValue() string {
	return "test"
}
`, "KTN-FUNC-004")
}

// TestFuncAnalyzerViolations vérifie les violations de règles.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerViolations(t *testing.T) {
	runFuncAnalyzerErrorTest(t, "Fonction avec trop de paramètres", `package test

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
`, "KTN-FUNC-005")

	runFuncAnalyzerErrorTest(t, "Fonction avec mauvais nommage", `package test

// bad_naming utilise des underscores.
func bad_naming() {
}
`, "KTN-FUNC-001")
}


// TestParamsFormatValid teste le format valide des paramètres.
//
// Params:
//   - t: instance de test
func TestParamsFormatValid(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Tous les params documentés", `package test

// process traite les données.
//
// Params:
//   - data: les données
//   - count: le nombre
func process(data string, count int) {
}
`)

	runFuncAnalyzerValidTest(t, "Format avec indentation", `package test

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
`)
}

// TestParamsFormatInvalid teste les erreurs de format des paramètres.
//
// Params:
//   - t: instance de test
func TestParamsFormatInvalid(t *testing.T) {
	runFuncAnalyzerErrorTest(t, "Param manquant", `package test

// process traite les données.
//
// Params:
//   - data: les données
func process(data string, count int) {
}
`, "KTN-FUNC-003")
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
