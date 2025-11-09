package ktntest

import (
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Analyzer008 checks that test files follow strict internal/external naming convention
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest008",
	Doc:  "KTN-TEST-008: Les fichiers de test doivent se terminer par _internal_test.go (package xxx) ou _external_test.go (package xxx_test)",
	Run:  runTest008,
}

// runTest008 exécute l'analyse KTN-TEST-008.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest008(pass *analysis.Pass) (any, error) {
	// Parcourir tous les fichiers du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		packageName := file.Name.Name

		// Vérifier si c'est un fichier de test
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test, continuer
			continue
		}

		baseName := filepath.Base(filename)

		// Vérifier la convention de nommage
		hasInternalSuffix := strings.HasSuffix(
			baseName,
			"_internal_test.go",
		)
		hasExternalSuffix := strings.HasSuffix(
			baseName,
			"_external_test.go",
		)

		// Le fichier doit avoir l'un des deux suffixes
		if !hasInternalSuffix && !hasExternalSuffix {
			pass.Reportf(
				file.Name.Pos(),
				"KTN-TEST-008: le fichier de test '%s' doit se terminer par '_internal_test.go' ou '_external_test.go'",
				baseName,
			)
			// Continuer au fichier suivant
			continue
		}

		// Vérifier la cohérence package/suffixe
		if hasInternalSuffix {
			expectedPackage := extractExpectedPackage(baseName)
			// Le package doit être le même que le fichier source
			if strings.HasSuffix(packageName, "_test") {
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' doit utiliser 'package %s' (pas '%s') pour les tests internes",
					baseName,
					expectedPackage,
					packageName,
				)
			}
		} else if hasExternalSuffix {
			// Le package doit se terminer par _test
			if !strings.HasSuffix(packageName, "_test") {
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' doit utiliser 'package xxx_test' (pas '%s') pour les tests externes",
					baseName,
					packageName,
				)
			}
		}
	}

	// Retour de la fonction
	return nil, nil
}

// extractExpectedPackage extrait le nom de package attendu d'un fichier _internal_test.go.
//
// Params:
//   - filename: nom du fichier de test
//
// Returns:
//   - string: nom de package attendu
func extractExpectedPackage(filename string) string {
	// Retirer le suffixe _internal_test.go
	baseName := strings.TrimSuffix(filename, "_internal_test.go")
	// Le package devrait être le nom du dossier parent
	// mais pour simplifier on retourne juste le base sans le suffixe
	// Early return from function.
	return baseName
}
