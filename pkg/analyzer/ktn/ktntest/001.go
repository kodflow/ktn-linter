package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 checks that test files use external test packages (xxx_test)
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest001",
	Doc:      "KTN-TEST-001: Les fichiers de test doivent utiliser le package xxx_test",
	Run:      runTest001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest001 exécute l'analyse KTN-TEST-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest001(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Vérifier si on est dans un fichier de test
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename
		// Vérification de la condition
		if !strings.HasSuffix(filename, "_test.go") {
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
				"KTN-TEST-001: le fichier de test doit utiliser le package '%s' au lieu de '%s'",
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

	// Parcours des packages exemptés
	for _, exempt := range exemptPkgs {
		// Vérification de la condition
		if pkgName == exempt {
			// Package exempté
			return true
		}
	}

	// Package non exempté
	return false
}
