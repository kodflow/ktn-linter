// Package ktntest implements KTN linter rules.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

const (
	ruleCodeTest006 string = "KTN-TEST-006"
)

// Analyzer006 checks that each test file has a corresponding source file (1:1 pattern)
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest006",
	Doc:  "KTN-TEST-006: Chaque fichier _test.go doit correspondre à un fichier source (pattern 1:1)",
	Run:  runTest006,
}

// runTest006 exécute l'analyse KTN-TEST-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest006(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest006) {
		// Règle désactivée
		return nil, nil
	}

	// Ignore xxx_test packages (test entire package)
	packageName := pass.Pkg.Name()
	// Vérification de la condition
	if strings.HasSuffix(packageName, "_test") {
		// Package xxx_test - pas de contrainte 1:1
		// Retour si package test
		return nil, nil
	}

	// Collecter tous les fichiers du package
	sourceFiles := make(map[string]bool, 0)
	// Map baseName -> file AST node
	testFiles := make(map[string]*testFileInfo, 0)

	// Itération sur les fichiers
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeTest006, filename) {
			// Fichier exclu
			continue
		}

		basename := filepath.Base(filename)

		// Vérification si test
		if shared.IsTestFile(basename) {
			// Extraire le nom de base
			var baseName string
			var base string
			var ok bool
			// Vérification suffixe internal
			if base, ok = strings.CutSuffix(basename, "_internal_test.go"); ok {
				// Définir le nom de base
				baseName = base
			} else {
				// Vérification suffixe external
				if base, ok = strings.CutSuffix(basename, "_external_test.go"); ok {
					// Cas alternatif: external
					// Définir le nom de base
					baseName = base
				} else {
					// Cas alternatif: test standard
					// Définir le nom de base
					baseName = strings.TrimSuffix(basename, "_test.go")
				}
			}
			// Ajouter le fichier de test
			testFiles[baseName] = &testFileInfo{
				basename: basename,
				filename: filename,
				fileNode: file,
			}
		} else {
			// Cas alternatif: fichier source
			// Ajouter le fichier source
			baseName := strings.TrimSuffix(basename, ".go")
			sourceFiles[baseName] = true
		}
	}

	// Itération sur les tests
	for baseName, info := range testFiles {
		// Vérification si source existe
		if !sourceFiles[baseName] {
			// Signaler l'erreur
			pass.Reportf(
				info.fileNode.Pos(),
				"KTN-TEST-006: fichier de test '%s' n'a pas de fichier source correspondant '%s.go' dans le même package. Dispatcher son contenu dans les fichiers de test appropriés puis le supprimer",
				info.basename,
				baseName,
			)
		}
	}

	// Retour de la fonction
	return nil, nil
}

// testFileInfo stores information about a test file
type testFileInfo struct {
	basename string
	filename string
	fileNode *ast.File
}
