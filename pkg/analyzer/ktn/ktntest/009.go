// Analyzer 009 for the ktntest package.
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
	// INITIAL_PUBLIC_FUNCS_CAP initial capacity for public funcs map
	INITIAL_PUBLIC_FUNCS_CAP int = 32
)

// Analyzer009 checks that public function tests are in external test files
var Analyzer009 = &analysis.Analyzer{
	Name:     "ktntest009",
	Doc:      "KTN-TEST-009: Les tests de fonctions publiques (exportées) doivent être dans _external_test.go uniquement (black-box testing)",
	Run:      runTest009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest009 exécute l'analyse KTN-TEST-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest009(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter toutes les fonctions publiques
	publicFunctions := collectPublicFunctions(pass, insp)

	// Vérifier les tests dans les fichiers _internal_test.go
	checkInternalTestsForPublicFunctions(pass, insp, publicFunctions)

	// Retour de la fonction
	return nil, nil
}

// collectPublicFunctions collecte les fonctions publiques du package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des noms de fonctions publiques
func collectPublicFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	publicFunctions := make(map[string]bool, INITIAL_PUBLIC_FUNCS_CAP)
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

		// Ajouter la fonction si elle est publique
		addPublicFunction(funcDecl, publicFunctions)
	})

	// Retour de la map
	return publicFunctions
}

// addPublicFunction ajoute une fonction publique à la map.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - publicFunctions: map des fonctions publiques
func addPublicFunction(funcDecl *ast.FuncDecl, publicFunctions map[string]bool) {
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

	// Only add public functions
	if meta.Visibility != shared.VisPublic {
		return
	}

	// Add lookup key
	key := shared.BuildTestLookupKey(meta)
	// Vérification clé valide
	if key != "" {
		publicFunctions[key] = true
	}
}

// checkInternalTestsForPublicFunctions vérifie les tests de fonctions publiques.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - publicFunctions: map des fonctions publiques
func checkInternalTestsForPublicFunctions(pass *analysis.Pass, insp *inspector.Inspector, publicFunctions map[string]bool) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un test dans un fichier internal
		if !strings.HasSuffix(baseName, "_internal_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Skip exempt test files
		if shared.IsExemptTestFile(filename) {
			return
		}

		// Vérifier si c'est un test de fonction publique
		checkAndReportPublicFunctionTest(pass, funcDecl, baseName, publicFunctions)
	})
}

// checkAndReportPublicFunctionTest vérifie et reporte un test de fonction publique mal placé.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction de test
//   - baseName: nom de base du fichier
//   - publicFunctions: map des fonctions publiques
func checkAndReportPublicFunctionTest(pass *analysis.Pass, funcDecl *ast.FuncDecl, baseName string, publicFunctions map[string]bool) {
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

	// Build lookup key
	key := shared.BuildTestTargetKey(target)
	// Vérification clé vide
	if key == "" {
		return
	}

	// Vérifier si c'est un test de fonction publique
	if publicFunctions[key] {
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-009: le test '%s' dans '%s' teste une fonction publique '%s'. Les tests de fonctions publiques doivent être dans '%s' (black-box testing avec package xxx_test)",
			testName, baseName, key,
			strings.Replace(baseName, "_internal_test.go", "_external_test.go", 1),
		)
	}
}
