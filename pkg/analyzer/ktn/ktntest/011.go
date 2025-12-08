// Analyzer 011 for the ktntest package.
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

// Analyzer011 checks package naming convention for internal/external test files
var Analyzer011 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest011",
	Doc:      "KTN-TEST-011: Les fichiers _internal_test.go doivent utiliser 'package xxx', les fichiers _external_test.go doivent utiliser 'package xxx_test'",
	Run:      runTest011,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest011 exécute l'analyse KTN-TEST-011.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest011(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Parcourir tous les fichiers de test
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		basename := filepath.Base(filename)

		// Vérification si fichier de test
		if !shared.IsTestFile(basename) {
			// Pas un fichier de test
			continue
		}

		actualPkg := file.Name.Name

		// Vérifier les conventions internal/external
		if strings.HasSuffix(basename, "_internal_test.go") {
			// Fichier _internal_test.go → doit utiliser package xxx (sans _test)
			if basePkg, ok := strings.CutSuffix(actualPkg, "_test"); ok {
				// Erreur : package xxx_test dans _internal_test.go
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-011: le fichier '%s' doit utiliser 'package %s' (white-box testing) au lieu de 'package %s'. Les fichiers _internal_test.go testent les fonctions privées et doivent partager le même package",
					basename,
					basePkg,
					actualPkg,
				)
			}
			// Verification de la condition
			// Alternative path handling
		} else if strings.HasSuffix(basename, "_external_test.go") {
			// Fichier _external_test.go → doit utiliser package xxx_test
			if !strings.HasSuffix(actualPkg, "_test") {
				// Extraire le nom du package attendu depuis le nom de fichier
				expectedPkg := extractExpectedPackageFromFilename(basename)
				// Erreur : package xxx dans _external_test.go
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-011: le fichier '%s' doit utiliser 'package %s_test' (black-box testing) au lieu de 'package %s'. Les fichiers _external_test.go testent l'API publique et doivent utiliser un package externe",
					basename,
					expectedPkg,
					actualPkg,
				)
			}
		}
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Rien à faire dans preorder
	})

	// Retour de la fonction
	return nil, nil
}

// extractExpectedPackageFromFilename extrait le nom de package attendu depuis le nom de fichier.
//
// Params:
//   - filename: nom du fichier (ex: calculator_external_test.go)
//
// Returns:
//   - string: nom du package attendu (ex: calculator)
func extractExpectedPackageFromFilename(filename string) string {
	// Retirer _external_test.go
	baseName := strings.TrimSuffix(filename, "_external_test.go")
	// Retour du nom de base
	return baseName
}
