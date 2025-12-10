// Analyzer 009 for the ktntest package.
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
	ruleCodeTest009 string = "KTN-TEST-009"
	// initialPublicFuncsCap initial capacity for public funcs map
	initialPublicFuncsCap int = 32
)

// Analyzer009 checks that public function tests are in external test files
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest009) {
		// Règle désactivée
		return nil, nil
	}

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
	publicFunctions := make(map[string]bool, initialPublicFuncsCap)
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

		// Ajouter fonction publique
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

	// Vérification visibilité publique
	if meta.Visibility != shared.VisPublic {
		// Retour si pas publique
		return
	}

	// Construire clé recherche
	key := shared.BuildTestLookupKey(meta)
	// Vérification clé non vide
	if key != "" {
		// Ajouter à la map
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
	// Récupération de la configuration
	cfg := config.Get()

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les tests
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest009, filename) {
			// Fichier exclu
			return
		}

		baseName := filepath.Base(filename)

		// Vérification fichier internal et test unitaire
		if !strings.HasSuffix(baseName, "_internal_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			// Retour si pas internal ou pas unitaire
			return
		}

		// Vérification fichier exempté
		if shared.IsExemptTestFile(filename) {
			// Retour si exempté
			return
		}

		// Vérifier et reporter
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

	// Construire clé
	key := shared.BuildTestTargetKey(target)
	// Vérification clé
	if key == "" {
		// Retour si clé vide
		return
	}

	// Vérification fonction publique
	if publicFunctions[key] {
		// Signaler erreur
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-009: le test '%s' dans '%s' teste une fonction publique '%s'. Les tests de fonctions publiques doivent être dans '%s' (black-box testing avec package xxx_test)",
			testName, baseName, key,
			strings.Replace(baseName, "_internal_test.go", "_external_test.go", 1),
		)
	}
}
