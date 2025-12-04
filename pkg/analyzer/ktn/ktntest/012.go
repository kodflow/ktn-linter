// Analyzer 012 for the ktntest package.
package ktntest

import (
	"go/ast"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	// assertionMethods contient les méthodes d'assertion de testing.T.
	assertionMethods = []string{
		"Error", "Errorf",
		"Fatal", "Fatalf",
		"Fail", "FailNow",
		"Log", "Logf",
		"Run", "Skip", "Skipf", "SkipNow",
		"Parallel", "Helper", "Cleanup",
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

		// Vérifier si c'est une fonction de test
		if !shared.IsUnitTestFunction(funcDecl) {
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

	// Vérifier si le test contient des assertions
	hasAssertion := false

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérifier les appels de fonction
		if checkForAssertion(n) {
			hasAssertion = true
			// Pas besoin de continuer
			return false
		}
		// Continuer la traversée
		return true
	})

	// Passthrough si pas d'assertion
	return !hasAssertion
}

// checkForAssertion vérifie si un nœud contient une assertion.
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - bool: true si assertion trouvée
func checkForAssertion(n ast.Node) bool {
	// Vérifier les appels de fonction
	callExpr, ok := n.(*ast.CallExpr)
	// Pas un appel de fonction
	if !ok {
		return false
	}

	// Vérifier les différents types d'assertions
	if isTestingMethodCall(callExpr) {
		return true
	}

	// Vérifier les bibliothèques d'assertion
	if isAssertLibraryCall(callExpr) {
		return true
	}

	// Vérifier si c'est un appel de helper avec t
	return isTestHelperCall(callExpr)
}

// isTestingMethodCall vérifie si c'est un appel t.Error, t.Fatal, etc.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si c'est une méthode de testing
func isTestingMethodCall(callExpr *ast.CallExpr) bool {
	// Vérifier si c'est un sélecteur (x.Method)
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	// Pas un sélecteur
	if !ok {
		return false
	}

	// Vérifier si c'est une méthode d'assertion
	return slices.Contains(assertionMethods, sel.Sel.Name)
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
