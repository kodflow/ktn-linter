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
	assertionMethods []string = []string{
		"Error", "Errorf",
		"Fatal", "Fatalf",
		"Fail", "FailNow",
	}

	// subTestMethods contient les méthodes qui lancent des sous-tests.
	subTestMethods []string = []string{
		"Run",
	}

	// Analyzer012 detects passthrough tests that don't test anything.
	Analyzer012 *analysis.Analyzer = &analysis.Analyzer{
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

	// Parcourir les fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Conversion en FuncDecl
		funcDecl := n.(*ast.FuncDecl)
		// Obtenir le chemin du fichier
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification si test
		if !shared.IsTestFile(filename) {
			// Retour si pas test
			return
		}

		// Vérification fichier exempté
		if shared.IsExemptTestFile(filename) {
			// Retour si exempté
			return
		}

		// Vérification fichier mock
		if shared.IsMockFile(filename) {
			// Retour si mock
			return
		}

		// Vérification si test unitaire
		if !shared.IsUnitTestFunction(funcDecl) {
			// Retour si pas unitaire
			return
		}

		// Vérification nom exempté
		if shared.IsExemptTestName(funcDecl.Name.Name) {
			// Retour si exempté
			return
		}

		// Vérification si passthrough
		if isPassthroughTest(funcDecl) {
			// Signaler test passthrough
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
	// Vérification corps vide
	if funcDecl.Body == nil || len(funcDecl.Body.List) == 0 {
		// Retour passthrough si vide
		return true
	}

	// Initialiser validation
	hasValidation := false

	// Parcourir le corps
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérification si déjà trouvé
		if hasValidation || n == nil {
			// Arrêter la recherche
			return false
		}

		// Vérification signal validation
		if checkForValidationSignal(n) {
			// Marquer comme trouvé
			hasValidation = true
			// Arrêter la recherche
			return false
		}

		// Continuer la recherche
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
	// Sélection selon le type
	switch node := n.(type) {
	// Vérification appel de fonction
	case *ast.CallExpr:
		// Cas appel de fonction
		// Vérifier appel validation
		return checkCallForValidation(node)
	// Vérification expression binaire
	case *ast.BinaryExpr:
		// Cas expression binaire
		// Vérifier opérateur comparaison
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
	// Vérification assertion testing.T
	if isTestingAssertionCall(callExpr) {
		// Retour si assertion trouvée
		return true
	}

	// Vérification sous-test
	if isSubTestCall(callExpr) {
		// Retour si sous-test trouvé
		return true
	}

	// Vérification bibliothèque assertion
	if isAssertLibraryCall(callExpr) {
		// Retour si assertion trouvée
		return true
	}

	// Retour vérification helper
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
	// Sélection selon opérateur
	switch op {
	// Vérification opérateurs de comparaison
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		// Cas opérateur comparaison
		// Retour trouvé
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
	// Conversion en SelectorExpr
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Vérification si sélecteur
	if !ok {
		// Retour si pas sélecteur
		return false
	}

	// Retour vérification méthode
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
	// Conversion en SelectorExpr
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Vérification si sélecteur
	if !ok {
		// Retour si pas sélecteur
		return false
	}

	// Vérifier le receiver
	recv, recvOk := sel.X.(*ast.Ident)
	// Vérification receiver t
	if !recvOk || recv.Name != "t" {
		// Retour si pas t
		return false
	}

	// Retour vérification Run
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
	// Conversion en SelectorExpr
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Vérification si sélecteur
	if !ok {
		// Retour si pas sélecteur
		return false
	}

	// Vérifier le receiver
	recv, recvOk := sel.X.(*ast.Ident)
	// Vérification identifiant
	if !recvOk {
		// Retour si pas identifiant
		return false
	}

	// Vérifier bibliothèque assertion
	receiverName := strings.ToLower(recv.Name)
	// Retour vérification nom
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
	// Vérification arguments présents
	if len(callExpr.Args) == 0 {
		// Retour si pas d'arguments
		return false
	}

	// Vérifier premier argument
	firstArg, ok := callExpr.Args[0].(*ast.Ident)
	// Vérification identifiant
	if !ok {
		// Retour si pas identifiant
		return false
	}

	// Retour vérification nom t
	return firstArg.Name == "t"
}
