// Analyzer 003 for the ktntest package.
package ktntest

import (
	"go/ast"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest003 string = "KTN-TEST-003"
)

// Analyzer003 checks that test files have corresponding source files
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest003) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Vérifier chaque fichier de test
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest003, filename) {
			// Fichier exclu
			continue
		}

		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test
			continue
		}

		// Extraire le nom du fichier source correspondant
		sourceFile := getSourceFileForTest(filename)

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

// getSourceFileForTest retourne le chemin du fichier source correspondant à un fichier de test.
//
// Params:
//   - filename: chemin du fichier de test
//
// Returns:
//   - string: chemin du fichier source correspondant
func getSourceFileForTest(filename string) string {
	var base string
	var ok bool

	// Vérification du suffixe internal
	if base, ok = strings.CutSuffix(filename, "_internal_test.go"); ok {
		// Retour du fichier source
		return base + ".go"
	}
	// Vérification du suffixe external
	if base, ok = strings.CutSuffix(filename, "_external_test.go"); ok {
		// Retour du fichier source
		return base + ".go"
	}
	// Vérification du suffixe bench
	if base, ok = strings.CutSuffix(filename, "_bench_test.go"); ok {
		// Retour du fichier source
		return base + ".go"
	}
	// Vérification du suffixe integration
	if base, ok = strings.CutSuffix(filename, "_integration_test.go"); ok {
		// Retour du fichier source
		return base + ".go"
	}
	// Retour pour fichier test standard
	return strings.TrimSuffix(filename, "_test.go") + ".go"
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
