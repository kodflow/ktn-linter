// Analyzer 010 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	ruleCodeTest010 string = "KTN-TEST-010"
	// initialPrivateFunctionsCap initial cap for private funcs map
	initialPrivateFunctionsCap int = 32
)

// Analyzer010 checks private function tests are in internal test files
var Analyzer010 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest010) {
		// Règle désactivée
		return nil, nil
	}

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
	privateFunctions := make(map[string]bool, initialPrivateFunctionsCap)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification fichier test
		if shared.IsTestFile(filename) {
			// Retour si test
			return
		}

		// Vérification fichier mock
		if shared.IsMockFile(filename) {
			// Retour si mock
			return
		}

		// Ajouter fonction privée
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
	// Vérification nom fonction
	if funcDecl.Name == nil || len(funcDecl.Name.Name) == 0 {
		// Retour si pas de nom
		return
	}

	// Vérification fonction mock
	if shared.IsMockName(funcDecl.Name.Name) {
		// Retour si mock
		return
	}

	// Classifier la fonction
	meta := shared.ClassifyFunc(funcDecl)

	// Vérification receiver mock
	if meta.ReceiverName != "" && shared.IsMockName(meta.ReceiverName) {
		// Retour si mock
		return
	}

	// Vérification visibilité privée
	if meta.Visibility != shared.VisPrivate {
		// Retour si pas privée
		return
	}

	// Construire clé recherche
	key := shared.BuildTestLookupKey(meta)
	// Vérification clé non vide
	if key != "" {
		// Ajouter à la map
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
	// Récupération de la configuration
	cfg := config.Get()

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les tests
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest010, filename) {
			// Fichier exclu
			return
		}

		baseName := filepath.Base(filename)

		// Vérification fichier external et test unitaire
		if !strings.HasSuffix(baseName, "_external_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			// Retour si pas external ou pas unitaire
			return
		}

		// Vérification fichier exempté
		if shared.IsExemptTestFile(filename) {
			// Retour si exempté
			return
		}

		// Vérifier et reporter
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

	// Vérification nom exempté
	if shared.IsExemptTestName(testName) {
		// Retour si exempté
		return
	}

	// Parser nom test
	target, ok := shared.ParseTestName(testName)
	// Vérification parsing
	if !ok {
		// Retour si échec parsing
		return
	}

	// Vérification fonction privée
	if !target.IsPrivate {
		// Retour si pas privée
		return
	}

	// Construire clé
	key := shared.BuildTestTargetKey(target)
	// Vérification clé
	if key == "" {
		// Retour si clé vide
		return
	}

	// Vérification fonction privée
	if privateFunctions[key] {
		// Signaler erreur
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-010: le test '%s' dans '%s' teste une fonction privée '%s'. Les tests de fonctions privées doivent être dans '%s' (white-box testing avec package xxx)",
			testName, baseName, target.FuncName,
			strings.Replace(baseName, "_external_test.go", "_internal_test.go", 1),
		)
	}
}
