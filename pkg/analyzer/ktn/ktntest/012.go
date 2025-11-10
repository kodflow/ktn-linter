package ktntest

import (
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer012 detects test files not following _internal_test.go or _external_test.go naming convention.
var Analyzer012 = &analysis.Analyzer{
	Name: "ktntest012",
	Doc:  "KTN-TEST-012: fichier de test doit se terminer par _internal_test.go ou _external_test.go",
	Run:  runTest012,
}

// runTest012 analyzes test file naming conventions.
//
// Params:
//   - pass: Analysis pass
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runTest012(pass *analysis.Pass) (any, error) {
	// Check each file in the package
	// Itération sur les éléments
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		base := filepath.Base(filename)

		// Skip if not a test file
		if !strings.HasSuffix(base, "_test.go") {
			continue
		}

		// Check if it follows the correct naming convention
		hasInternalSuffix := strings.HasSuffix(base, "_internal_test.go")
		hasExternalSuffix := strings.HasSuffix(base, "_external_test.go")

		// Report error if it's a plain _test.go file
		// Vérification de la condition
		if !hasInternalSuffix && !hasExternalSuffix {
			pass.Reportf(
				file.Pos(),
				"KTN-TEST-012: le fichier '%s' doit être renommé en '%s' (white-box) ou '%s' (black-box), ou son contenu doit être dispatché dans ces fichiers. Les tests publics doivent aller dans _external_test.go (package xxx_test), les tests de fonctions privées dans _internal_test.go (package xxx)",
				base,
				strings.Replace(base, "_test.go", "_internal_test.go", 1),
				strings.Replace(base, "_test.go", "_external_test.go", 1),
			)
		}
	}

	// Early return from function.
	return nil, nil
}
