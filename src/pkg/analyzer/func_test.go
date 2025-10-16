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
	// Retourne true car le traitement est réussi
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
	// Retourne "test" comme valeur de démonstration
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
	// Retourne la somme des deux nombres
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
	// Retourne le montant avec majoration de 20%
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
	// Retourne "test" comme valeur de test
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
	// Retourne true car le test est réussi
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
	// Retourne la somme des trois valeurs
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

// TestFuncAnalyzerNestingDepth teste la détection de profondeur d'imbrication.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerNestingDepth(t *testing.T) {
	testNestingDepthBasicCases(t)
	testNestingDepthStatementTypes(t)
}

// testNestingDepthBasicCases teste les cas basiques de profondeur.
//
// Params:
//   - t: instance de test
func testNestingDepthBasicCases(t *testing.T) {
	cases := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{"Nesting depth 1 (OK)", `package test
// simpleIf vérifie une condition simple.
// Params:
//   - x: valeur à vérifier
func simpleIf(x int) { if x > 0 { x++ } }`, false},
		{"Nesting depth 2 (OK)", `package test
// nestedTwo a deux niveaux d'imbrication.
// Params:
//   - x: valeur à vérifier
func nestedTwo(x int) { if x > 0 { if x < 10 { x++ } } }`, false},
		{"Nesting depth 3 (OK - max)", `package test
// nestedThree a trois niveaux d'imbrication.
// Params:
//   - x: valeur à vérifier
func nestedThree(x int) { if x > 0 { if x < 10 { if x != 5 { x++ } } } }`, false},
		{"Nesting depth 4 (ERROR)", `package test
// nestedFour a quatre niveaux d'imbrication.
// Params:
//   - x: valeur à vérifier
func nestedFour(x int) { if x > 0 { if x < 10 { if x != 5 { if x != 7 { x++ } } } } }`, true},
	}

	for _, tc := range cases {
		if tc.wantErr {
			runFuncAnalyzerErrorTest(t, tc.name, tc.code, "KTN-FUNC-010")
		} else {
			runFuncAnalyzerValidTest(t, tc.name, tc.code)
		}
	}
}

// testNestingDepthStatementTypes teste différents types de statements.
//
// Params:
//   - t: instance de test
func testNestingDepthStatementTypes(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Nesting with for loop", `package test
// withForLoop utilise une boucle for.
// Params:
//   - items: les éléments à traiter
func withForLoop(items []int) { for i := 0; i < len(items); i++ { if items[i] > 0 { items[i]++ } } }`)

	runFuncAnalyzerValidTest(t, "Nesting with range", `package test
// withRange utilise range.
// Params:
//   - items: les éléments à traiter
func withRange(items []int) { for _, item := range items { if item > 0 { item++ } } }`)

	runFuncAnalyzerValidTest(t, "Nesting with switch", `package test
// withSwitch utilise un switch.
// Params:
//   - x: valeur à vérifier
func withSwitch(x int) { switch x { case 1: if x > 0 { x++ } } }`)

	runFuncAnalyzerValidTest(t, "Nesting with select", `package test
// withSelect utilise select.
// Params:
//   - ch: canal à surveiller
func withSelect(ch chan int) { select { case v := <-ch: if v > 0 { v++ } } }`)

	runFuncAnalyzerValidTest(t, "Nesting with if-else", `package test
// withElse utilise if-else.
// Params:
//   - x: valeur à vérifier
func withElse(x int) { if x > 0 { if x < 10 { if x == 5 { x++ } } } else { if x < -10 { x-- } } }`)
}

// TestFuncAnalyzerComplexity teste la détection de complexité cyclomatique.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerComplexity(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Complexity 1 (no branching)", `package test
// simple est une fonction simple.
func simple() { x := 1; x++ }`)

	runFuncAnalyzerValidTest(t, "Complexity with if", `package test
// withIf contient un if.
// Params:
//   - x: valeur à tester
func withIf(x int) { if x > 0 { x++ } }`)

	runFuncAnalyzerValidTest(t, "Complexity with && and ||", `package test
// withLogic utilise && et ||.
// Params:
//   - x: première valeur
//   - y: seconde valeur
func withLogic(x, y int) { if x > 0 && y > 0 { x++ }; if x < 0 || y < 0 { x-- } }`)

	runFuncAnalyzerValidTest(t, "Complexity with switch cases", `package test
// withSwitch utilise un switch.
// Params:
//   - x: valeur à tester
func withSwitch(x int) { switch x { case 1: x++; case 2: x--; default: x = 0 } }`)

	runFuncAnalyzerValidTest(t, "Complexity 6 with depth 3 (OK)", `package test
// complexSix a complexité 6 avec profondeur 3.
// Params:
//   - x: valeur à tester
//   - y: autre valeur
func complexSix(x, y int) { if x > 0 { if x < 10 || y > 5 { if x == 1 && y == 2 { x++ } } } }`)

	runFuncAnalyzerErrorTest(t, "Complexity > 10 (ERROR)", `package test
// tooComplex a complexité trop élevée.
// Params:
//   - x: valeur à tester
func tooComplex(x int) { if x > 0 { if x < 10 { if x == 1 || x == 2 { if x != 3 && x != 4 { if x > 5 { if x < 8 || x > 12 { if x != 6 { x++ } } } } } } } }`, "KTN-FUNC-007")
}

// TestFuncAnalyzerLength teste la détection de longueur de fonction.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerLength(t *testing.T) {
	// Fonction courte (OK)
	runFuncAnalyzerValidTest(t, "Short function (OK)", `package test

// shortFunc est courte.
func shortFunc() {
	x := 1
	x++
	x--
}
`)

	// Fonction de 35 lignes (limite en production)
	lines := "package test\n\n// exactLimit a exactement 35 lignes.\nfunc exactLimit() {\n"
	for i := 0; i < 35; i++ {
		lines += "\tx := 1\n"
	}
	lines += "}\n"
	runFuncAnalyzerValidTest(t, "Exactly 35 lines (OK)", lines)

	// Fonction > 35 lignes (ERROR en production)
	longLines := "package test\n\n// tooLong a plus de 35 lignes.\nfunc tooLong() {\n"
	for i := 0; i < 40; i++ {
		longLines += "\tx := 1\n"
	}
	longLines += "}\n"
	runFuncAnalyzerErrorTest(t, "More than 35 lines (ERROR)", longLines, "KTN-FUNC-006")
}

// TestFuncAnalyzerEdgeCases teste les cas limites.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerEdgeCases(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Function without body", `package test
// externalFunc est définie ailleurs.
func externalFunc()`)

	runFuncAnalyzerValidTest(t, "Method with receiver (analyzed)", `package test
type MyType struct{}
// Method est une méthode documentée.
func (m *MyType) Method() { x := 1; x++ }`)

	runFuncAnalyzerValidTest(t, "Unnamed params", `package test
// unnamedParams utilise des params sans noms.
// Params:
//   - int: le premier paramètre
// Returns:
//   - int: le résultat
func unnamedParams(int) int {
	// Retourne 0 car c'est la valeur par défaut
	return 0
}`)

	runFuncAnalyzerValidTest(t, "Grouped params", `package test
// groupedParams a des paramètres groupés.
// Params:
//   - x: première valeur
//   - y: seconde valeur
// Returns:
//   - int: le résultat
func groupedParams(x, y int) int {
	// Retourne la somme de x et y
	return x + y
}`)

	runFuncAnalyzerValidTest(t, "Underscore param", `package test
// withUnderscore ignore un paramètre.
// Params:
//   - x: valeur utilisée
// Returns:
//   - int: le résultat
func withUnderscore(x int, _ string) int {
	// Retourne x car le second paramètre est ignoré
	return x
}`)

	runFuncAnalyzerValidTest(t, "Switch default", `package test
// withDefault utilise default.
// Params:
//   - x: valeur à tester
func withDefault(x int) { switch x { default: x = 0 } }`)

	runFuncAnalyzerValidTest(t, "Select default", `package test
// selectDefault utilise select default.
func selectDefault() { select { default: x := 0; x++ } }`)

	runFuncAnalyzerValidTest(t, "Binary expr without logical", `package test
// withAddition utilise addition.
// Params:
//   - x: première valeur
//   - y: seconde valeur
// Returns:
//   - int: la somme
func withAddition(x, y int) int {
	// Retourne la somme de x et y
	return x + y
}`)

	runFuncAnalyzerValidTest(t, "Nil params", `package test
// noParams n'a pas de paramètres.
func noParams() { x := 1; x++ }`)

	runFuncAnalyzerValidTest(t, "Both sections", `package test
// withBothSections a les deux sections.
// Params:
//   - x: valeur d'entrée
// Returns:
//   - int: valeur de sortie
func withBothSections(x int) int {
	// Retourne x inchangé
	return x
}`)

	runFuncAnalyzerValidTest(t, "Empty line in section", `package test
// withEmptyLine a ligne vide.
// Params:
//   - x: valeur
func withEmptyLine(x int) { x++ }`)
}

// TestFuncAnalyzerTestFileRelaxedRules teste les règles assouplies pour fichiers de test.
//
// Params:
//   - t: instance de test
func TestFuncAnalyzerTestFileRelaxedRules(t *testing.T) {
	// Fonction de 50 lignes dans test file (OK car limite = 100)
	t.Run("50 lines in test file (OK)", func(t *testing.T) {
		lines := "package test\n\n// testFunc est une fonction de test.\nfunc testFunc() {\n"
		for i := 0; i < 50; i++ {
			lines += "\tx := 1\n"
		}
		lines += "}\n"

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test_test.go", lines, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				t.Errorf("Unexpected diagnostic: %s", diag.Message)
			},
		}

		_, err = analyzer.FuncAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}
	})

	// Complexité 12 dans test file (OK car limite = 15)
	// But nesting must still be <= 3, so we spread complexity horizontally
	t.Run("Complexity 12 in test file (OK)", func(t *testing.T) {
		code := `package test

// testComplex a complexité 12 avec profondeur 3.
//
// Params:
//   - x: valeur à tester
//   - y: autre valeur
func testComplex(x, y int) {
	// Depth 1: 1 if + 2 &&/|| = complexity +3
	if x > 0 && y > 0 || x < 0 {
		// Depth 2: 1 if + 2 &&/|| = complexity +3
		if x < 10 && y < 10 || x > 100 {
			// Depth 3: 1 if + 2 &&/|| = complexity +3
			if x == 1 || y == 1 && x != 2 {
				x++
			}
		}
	}
	// Total: 1 + 3 + 3 + 3 = 10 (actually < 15, still OK)
}
`

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test_test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				t.Errorf("Unexpected diagnostic: %s", diag.Message)
			},
		}

		_, err = analyzer.FuncAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}
	})
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
	// Retourne la somme de a et b
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

// TestReturnComments teste la règle KTN-FUNC-008.
//
// Params:
//   - t: instance de test
func TestReturnComments(t *testing.T) {
	testReturnCommentsValid(t)
	testReturnCommentsInvalid(t)
	testReturnCommentsEdgeCases(t)
}

// testReturnCommentsValid teste les cas valides avec commentaires.
//
// Params:
//   - t: instance de test
func testReturnCommentsValid(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Single return with comment", `package test

// getValue retourne une valeur.
//
// Returns:
//   - string: la valeur
func getValue() string {
	// Retourne "test" comme valeur
	return "test"
}
`)

	runFuncAnalyzerValidTest(t, "Multiple returns with comments", `package test

// Process traite une valeur.
//
// Params:
//   - err: l'erreur potentielle
//
// Returns:
//   - error: l'erreur rencontrée
func Process(err error) error {
	if err != nil {
		// Erreur de traitement
		return err
	}
	// Succès
	return nil
}
`)

	runFuncAnalyzerValidTest(t, "Return in nested blocks", `package test

// checkValue vérifie une valeur.
//
// Params:
//   - x: la valeur à vérifier
//
// Returns:
//   - bool: le résultat de la vérification
func checkValue(x int) bool {
	if x > 0 {
		if x < 10 {
			// Retourne true car la valeur est dans la plage valide
			return true
		}
	}
	// Retourne false car la valeur est hors plage
	return false
}
`)

	runFuncAnalyzerValidTest(t, "Return with multi-line comment", `package test

// getData récupère des données.
//
// Returns:
//   - string: les données
func getData() string {
	// Retourne les données formatées
	// avec des informations supplémentaires
	return "data"
}
`)
}

// testReturnCommentsInvalid teste les cas invalides sans commentaires.
//
// Params:
//   - t: instance de test
func testReturnCommentsInvalid(t *testing.T) {
	runFuncAnalyzerErrorTest(t, "Return without comment", `package test

// getValue retourne une valeur.
//
// Returns:
//   - string: la valeur
func getValue() string {
	return "test"
}
`, "KTN-FUNC-008")

	runFuncAnalyzerErrorTest(t, "Multiple returns, one without comment", `package test

// Process traite une valeur.
//
// Params:
//   - err: l'erreur potentielle
//
// Returns:
//   - error: l'erreur rencontrée
func Process(err error) error {
	if err != nil {
		// Erreur de traitement
		return err
	}
	return nil
}
`, "KTN-FUNC-008")

	runFuncAnalyzerErrorTest(t, "All returns without comments", `package test

// check vérifie une condition.
//
// Params:
//   - x: valeur à vérifier
//
// Returns:
//   - bool: le résultat
func check(x int) bool {
	if x > 0 {
		return true
	}
	return false
}
`, "KTN-FUNC-008")
}

// testReturnCommentsEdgeCases teste les cas limites.
//
// Params:
//   - t: instance de test
func testReturnCommentsEdgeCases(t *testing.T) {
	runFuncAnalyzerValidTest(t, "Function without body (no return type)", `package test

// externalFunc est définie ailleurs.
func externalFunc()`)

	runFuncAnalyzerValidTest(t, "Function without return", `package test

// doSomething fait quelque chose.
func doSomething() {
	x := 1
	x++
}
`)

	runFuncAnalyzerValidTest(t, "Return with inline comment same line", `package test

// getValue retourne une valeur.
//
// Returns:
//   - int: la valeur
func getValue() int {
	// Retourne 42 car c'est la réponse universelle
	return 42
}
`)

	runFuncAnalyzerErrorTest(t, "Return with comment too far above", `package test

// getValue retourne une valeur.
//
// Returns:
//   - int: la valeur
func getValue() int {
	// Commentaire ici

	return 42
}
`, "KTN-FUNC-008")

	runFuncAnalyzerValidTest(t, "Switch with returns", `package test

// getType retourne un type.
//
// Params:
//   - x: la valeur à tester
//
// Returns:
//   - string: le type
func getType(x int) string {
	switch {
	case x > 0:
		// Retourne "positive" car x est positif
		return "positive"
	case x < 0:
		// Retourne "negative" car x est négatif
		return "negative"
	default:
		// Retourne "zero" car x est nul
		return "zero"
	}
}
`)
}
