package ktntest

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 checks that test files have corresponding source files
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest002",
	Doc:      "KTN-TEST-002: Chaque fichier _test.go doit avoir un fichier .go correspondant",
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
		sourceFile := strings.TrimSuffix(filename, "_test.go") + ".go"

		// Vérifier si le fichier source existe
		if !fileExists(sourceFile) && !isExemptTestFile(filename) {
			// Signaler l'erreur
			pass.Reportf(
				f.Name.Pos(),
				"KTN-TEST-002: fichier de test '%s' sans fichier source correspondant '%s'",
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

	// Parcours des patterns exemptés
	for _, pattern := range exemptPatterns {
		// Vérification de la condition
		if baseName == pattern {
			// Fichier exempté
			return true
		}
	}

	// Fichier non exempté
	return false
}
