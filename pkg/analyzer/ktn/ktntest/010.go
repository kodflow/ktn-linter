// Analyzer 010 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// INITIAL_PRIVATE_FUNCTIONS_CAP initial cap for private funcs map
	INITIAL_PRIVATE_FUNCTIONS_CAP int = 32
)

// Analyzer010 checks private function tests are in internal test files
var Analyzer010 = &analysis.Analyzer{
	Name:     "ktntest010",
	Doc:      "KTN-TEST-010: Les tests de fonctions privées (non-exportées) doivent être dans _internal_test.go uniquement (white-box testing)",
	Run:      runTest010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest010 exécute l'analyse KTN-TEST-010.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest010(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter toutes les fonctions privées
	privateFunctions := collectPrivateFunctions(pass, insp)

	// Vérifier les tests dans les fichiers _external_test.go
	checkExternalTestsForPrivateFunctions(pass, insp, privateFunctions)

	// Retour de la fonction
	return nil, nil
}

// collectPrivateFunctions collecte les fonctions privées du package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des noms de fonctions privées
func collectPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	privateFunctions := make(map[string]bool, INITIAL_PRIVATE_FUNCTIONS_CAP)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir toutes les déclarations de fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Ignorer les fichiers mock
		if shared.IsMockFile(filename) {
			return
		}

		// Ajouter la fonction si elle est privée
		addPrivateFunction(funcDecl, privateFunctions)
	})

	// Retour de la map
	return privateFunctions
}

// addPrivateFunction ajoute une fonction privée à la map.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - privateFunctions: map des fonctions privées
func addPrivateFunction(funcDecl *ast.FuncDecl, privateFunctions map[string]bool) {
	// Vérifier le nom de la fonction
	if funcDecl.Name == nil || len(funcDecl.Name.Name) == 0 {
		return
	}

	// Skip mock functions
	if shared.IsMockName(funcDecl.Name.Name) {
		return
	}

	// Use shared helper to classify function
	meta := shared.ClassifyFunc(funcDecl)

	// Skip mock receiver types
	if meta.ReceiverName != "" && shared.IsMockName(meta.ReceiverName) {
		return
	}

	// Only add private functions
	if meta.Visibility != shared.VIS_PRIVATE {
		return
	}

	// Add lookup key
	key := shared.BuildTestLookupKey(meta)
	// Vérification clé valide
	if key != "" {
		privateFunctions[key] = true
	}
}

// checkExternalTestsForPrivateFunctions vérifie les tests de fonctions privées.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - privateFunctions: map des fonctions privées
func checkExternalTestsForPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector, privateFunctions map[string]bool) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un test dans un fichier external
		if !strings.HasSuffix(baseName, "_external_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Skip exempt test files
		if shared.IsExemptTestFile(filename) {
			return
		}

		// Vérifier si c'est un test de fonction privée
		checkAndReportPrivateFunctionTest(pass, funcDecl, baseName, privateFunctions)
	})
}

// checkAndReportPrivateFunctionTest vérifie et reporte un test de fonction privée mal placé.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction de test
//   - baseName: nom de base du fichier
//   - privateFunctions: map des fonctions privées
func checkAndReportPrivateFunctionTest(pass *analysis.Pass, funcDecl *ast.FuncDecl, baseName string, privateFunctions map[string]bool) {
	testName := funcDecl.Name.Name

	// Skip exempt test names
	if shared.IsExemptTestName(testName) {
		return
	}

	// Use shared helper to parse test name
	target, ok := shared.ParseTestName(testName)
	// Vérification parsing réussi
	if !ok {
		return
	}

	// Only check tests targeting private functions
	if !target.IsPrivate {
		return
	}

	// Build lookup key
	key := shared.BuildTestTargetKey(target)
	// Vérification clé vide
	if key == "" {
		return
	}

	// Vérifier si c'est un test de fonction privée
	if privateFunctions[key] {
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-010: le test '%s' dans '%s' teste une fonction privée '%s'. Les tests de fonctions privées doivent être dans '%s' (white-box testing avec package xxx)",
			testName, baseName, target.FuncName,
			strings.Replace(baseName, "_external_test.go", "_internal_test.go", 1),
		)
	}
}
