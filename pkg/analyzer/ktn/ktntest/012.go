// Analyzer 012 for the ktntest package.
package ktntest

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	// assertionMethods contient uniquement les méthodes d'assertion de testing.T.
	// Log, Logf, Run, Skip*, Parallel, Helper, Cleanup ne sont PAS des assertions.
	assertionMethods = []string{
		"Error", "Errorf",
		"Fatal", "Fatalf",
		"Fail", "FailNow",
	}

	// subTestMethods contient les méthodes qui lancent des sous-tests.
	subTestMethods = []string{
		"Run",
	}

	// Analyzer012 detects passthrough tests that don't test anything.
	Analyzer012 = &analysis.Analyzer{
		Name:     "ktntest012",
		Doc:      "KTN-TEST-012: Les tests doivent contenir des assertions et vraiment tester quelque chose",
		Run:      runTest012,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

// runTest012 exécute l'analyse KTN-TEST-012.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest012(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Analyser les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification fichier de test
		if !shared.IsTestFile(filename) {
			return
		}

		// Skip exempt test files
		if shared.IsExemptTestFile(filename) {
			return
		}

		// Skip mock files
		if shared.IsMockFile(filename) {
			return
		}

		// Vérifier si c'est une fonction de test
		if !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Skip exempt test names
		if shared.IsExemptTestName(funcDecl.Name.Name) {
			return
		}

		// Vérifier si le test contient des assertions
		if isPassthroughTest(funcDecl) {
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-012: le test '%s' est un test passthrough - "+
					"il ne contient pas d'assertions et ne teste rien de substantiel",
				funcDecl.Name.Name,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isPassthroughTest vérifie si un test est passthrough.
// Un test est passthrough s'il ne contient aucun signal de validation:
//   - Pas d'assertion (t.Error, t.Fatal, assert.Equal, etc.)
//   - Pas de comparaison logique (==, !=, <, >, etc.)
//   - Pas de sous-test (t.Run)
//   - Pas d'appel à un helper de test
//
// Params:
//   - funcDecl: déclaration de fonction de test
//
// Returns:
//   - bool: true si le test est passthrough
func isPassthroughTest(funcDecl *ast.FuncDecl) bool {
	// Corps vide = passthrough
	if funcDecl.Body == nil || len(funcDecl.Body.List) == 0 {
		return true
	}

	// Vérifier si le test contient des signaux de validation
	hasValidation := false

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Déjà trouvé, pas besoin de continuer
		if hasValidation || n == nil {
			return false
		}

		// Vérifier les différents signaux
		if checkForValidationSignal(n) {
			hasValidation = true
			return false
		}

		// Continuer la traversée
		return true
	})

	// Passthrough si pas de signal de validation
	return !hasValidation
}

// checkForValidationSignal vérifie si un nœud contient un signal de validation.
// Cela inclut: assertions, comparaisons, sous-tests, helpers.
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - bool: true si signal trouvé
func checkForValidationSignal(n ast.Node) bool {
	// Vérifier selon le type de nœud
	switch node := n.(type) {
	case *ast.CallExpr:
		// Vérifier les appels de fonction
		return checkCallForValidation(node)
	case *ast.BinaryExpr:
		// Vérifier les comparaisons
		return isComparisonOperator(node.Op)
	}

	// Pas de signal trouvé
	return false
}

// checkCallForValidation vérifie si un appel est une validation.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est une validation
func checkCallForValidation(callExpr *ast.CallExpr) bool {
	// Méthodes sur testing.T : t.Error, t.Fatal, t.Fail, ...
	if isTestingAssertionCall(callExpr) {
		return true
	}

	// Sous-tests : t.Run
	if isSubTestCall(callExpr) {
		return true
	}

	// Bibliothèques d'assertion (assert, require, ...)
	if isAssertLibraryCall(callExpr) {
		return true
	}

	// Helper de test prenant t en premier argument
	return isTestHelperCall(callExpr)
}

// isComparisonOperator vérifie si un opérateur est une comparaison.
//
// Params:
//   - op: opérateur token
//
// Returns:
//   - bool: true si c'est une comparaison
func isComparisonOperator(op token.Token) bool {
	// Vérifier les opérateurs de comparaison
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return true
	}
	// Pas une comparaison
	return false
}

// isTestingAssertionCall vérifie si c'est un appel t.Error, t.Fatal, etc.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est une assertion testing.T
func isTestingAssertionCall(callExpr *ast.CallExpr) bool {
	// Vérifier si c'est un sélecteur (x.Method)
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Pas un sélecteur
	if !ok {
		return false
	}

	// Vérifier si c'est une méthode d'assertion
	return slices.Contains(assertionMethods, sel.Sel.Name)
}

// isSubTestCall vérifie si c'est un appel t.Run (sous-test).
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est un sous-test
func isSubTestCall(callExpr *ast.CallExpr) bool {
	// Vérifier si c'est un sélecteur (t.Run)
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Pas un sélecteur
	if !ok {
		return false
	}

	// Vérifier si le receiver est t
	recv, recvOk := sel.X.(*ast.Ident)
	// Pas un identifiant
	if !recvOk || recv.Name != "t" {
		return false
	}

	// Vérifier si c'est Run
	return slices.Contains(subTestMethods, sel.Sel.Name)
}

// isAssertLibraryCall vérifie si c'est un appel à une bibliothèque d'assertion.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est une assertion de bibliothèque
func isAssertLibraryCall(callExpr *ast.CallExpr) bool {
	// Vérifier si c'est un sélecteur (assert.Equal, require.NoError, etc.)
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Pas un sélecteur
	if !ok {
		return false
	}

	// Vérifier le receiver
	recv, recvOk := sel.X.(*ast.Ident)
	// Pas un identifiant simple
	if !recvOk {
		return false
	}

	// Vérifier si c'est une bibliothèque d'assertion connue
	receiverName := strings.ToLower(recv.Name)
	// Retour du résultat
	return receiverName == "assert" || receiverName == "require"
}

// isTestHelperCall vérifie si c'est un appel à un helper de test.
// Un helper de test est une fonction qui reçoit t (*testing.T) comme argument.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est un appel de helper de test
func isTestHelperCall(callExpr *ast.CallExpr) bool {
	// Vérifier si l'appel a des arguments
	if len(callExpr.Args) == 0 {
		return false
	}

	// Vérifier si le premier argument est t (le testing.T)
	firstArg, ok := callExpr.Args[0].(*ast.Ident)
	// Pas un identifiant
	if !ok {
		return false
	}

	// Si le premier argument s'appelle t, c'est probablement un helper
	return firstArg.Name == "t"
}
