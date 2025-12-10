// Analyzer 002 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"slices"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest002 string = "KTN-TEST-002"
)

// Analyzer002 checks that test files use external test packages (xxx_test)
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest002",
	Doc:      "KTN-TEST-002: Les fichiers de test doivent utiliser le package xxx_test",
	Run:      runTest002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest002 exécute l'analyse KTN-TEST-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest002) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Vérifier si on est dans un fichier de test
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest002, filename) {
			// Fichier exclu
			continue
		}

		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test, continuer
			continue
		}

		// Extraire le nom du package attendu
		dir := filepath.Dir(filename)
		expectedPkg := filepath.Base(dir) + "_test"

		// Vérifier le nom du package
		actualPkg := f.Name.Name
		// Vérification de la condition
		if actualPkg != expectedPkg && !isExemptPackage(actualPkg) {
			// Signaler l'erreur
			pass.Reportf(
				f.Name.Pos(),
				"KTN-TEST-002: le fichier de test doit utiliser le package '%s' au lieu de '%s'",
				expectedPkg,
				actualPkg,
			)
		}
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Nothing to do in preorder, already handled above
	})

	// Retour de la fonction
	return nil, nil
}

// isExemptPackage vérifie si un package est exempté de la règle.
//
// Params:
//   - pkgName: nom du package
//
// Returns:
//   - bool: true si le package est exempté
func isExemptPackage(pkgName string) bool {
	// Certains packages internes peuvent être exemptés
	// (ceux qui testent des fonctions privées)
	exemptPkgs := []string{
		"main",
		"testhelper",
		"cmd",       // tests de fonctions privées cmd
		"utils",     // tests de fonctions privées utils
		"formatter", // tests de fonctions privées formatter
		"ktn",       // registry tests need same package for KTN-TEST-003
		"ktnconst",  // registry tests need same package for KTN-TEST-003
		"ktnfunc",   // registry tests need same package for KTN-TEST-003
		"ktntest",   // registry tests need same package for KTN-TEST-003
		"ktnvar",    // registry tests need same package for KTN-TEST-003
	}

	// Vérifier si le package est exempté
	return slices.Contains(exemptPkgs, pkgName)
}
