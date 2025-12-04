// Analyzer 003 for the ktntest package.
package ktntest

import (
	"go/ast"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks that test files have corresponding source files
var Analyzer003 = &analysis.Analyzer{
	Name:     "ktntest003",
	Doc:      "KTN-TEST-003: Chaque fichier _test.go doit avoir un fichier .go correspondant",
	Run:      runTest003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest003 exécute l'analyse KTN-TEST-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest003(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Vérifier chaque fichier de test
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename
		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test
			continue
		}

		// Extraire le nom du fichier source correspondant
		var sourceFile string
		var base string
		var ok bool
		// Vérification du suffixe (convention internal/external)
		if base, ok = strings.CutSuffix(filename, "_internal_test.go"); ok {
			// Fichier _internal_test.go → chercher .go
			sourceFile = base + ".go"
			// Cas alternatif: external test
		} else if base, ok = strings.CutSuffix(filename, "_external_test.go"); ok {
			// Fichier _external_test.go → chercher .go
			sourceFile = base + ".go"
			// Cas par défaut: standard test
		} else {
			// Fichier _test.go standard → chercher .go
			sourceFile = strings.TrimSuffix(filename, "_test.go") + ".go"
		}

		// Vérifier si le fichier source existe
		if !fileExists(sourceFile) && !isExemptTestFile(filename) {
			// Signaler l'erreur
			pass.Reportf(
				f.Name.Pos(),
				"KTN-TEST-003: fichier de test '%s' sans fichier source correspondant '%s'",
				filepath.Base(filename),
				filepath.Base(sourceFile),
			)
		}
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Nothing to do in preorder
	})

	// Retour de la fonction
	return nil, nil
}

// fileExists vérifie si un fichier existe.
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - bool: true si le fichier existe
func fileExists(path string) bool {
	info, err := os.Stat(path)
	// Vérification de la condition
	if err != nil {
		// Erreur ou fichier n'existe pas
		return false
	}
	// Retour du résultat
	return !info.IsDir()
}

// isExemptTestFile vérifie si un fichier de test est exempté.
//
// Params:
//   - filename: nom du fichier
//
// Returns:
//   - bool: true si le fichier est exempté
func isExemptTestFile(filename string) bool {
	baseName := filepath.Base(filename)
	// Fichiers de test exemptés (integration tests, helper tests, etc.)
	exemptPatterns := []string{
		"helper_test.go",
		"integration_test.go",
		"suite_test.go",
		"main_test.go",
	}

	// Vérifier si le fichier est exempté
	return slices.Contains(exemptPatterns, baseName)
}
