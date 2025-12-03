// Package ktntest implements KTN linter rules.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
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
	// Ignore xxx_test packages (test entire package)
	packageName := pass.Pkg.Name()
	// Vérification de la condition
	if strings.HasSuffix(packageName, "_test") {
		// Package xxx_test - pas de contrainte 1:1
		return nil, nil
	}

	// Collecter tous les fichiers du package
	sourceFiles := make(map[string]bool, 0)
	// Map baseName -> file AST node
	testFiles := make(map[string]*testFileInfo, 0)

	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		basename := filepath.Base(filename)

		// Vérification si fichier de test
		if shared.IsTestFile(basename) {
			// Extract base name without test suffix
			var baseName string
			var base string
			var ok bool
			// Vérification du suffixe
			if base, ok = strings.CutSuffix(basename, "_internal_test.go"); ok {
				// Fichier _internal_test.go → chercher .go
				baseName = base
				// Cas alternatif: external test
			} else if base, ok = strings.CutSuffix(basename, "_external_test.go"); ok {
				// Fichier _external_test.go → chercher .go
				baseName = base
				// Cas par défaut: standard test
			} else {
				// Fichier _test.go standard → chercher .go
				baseName = strings.TrimSuffix(basename, "_test.go")
			}
			testFiles[baseName] = &testFileInfo{
				basename: basename,
				filename: filename,
				fileNode: file,
			}
		} else {
			// Fichier source .go (non-test)
			baseName := strings.TrimSuffix(basename, ".go")
			sourceFiles[baseName] = true
		}
	}

	// Vérifier chaque fichier de test
	for baseName, info := range testFiles {
		// Vérifier si le fichier source correspondant existe
		if !sourceFiles[baseName] {
			// Fichier de test orphelin - reporter à la position du fichier
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
