package ktntest

import (
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Analyzer008 checks that each source file has both internal and external test files
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest008",
	Doc:  "KTN-TEST-008: Chaque fichier .go doit avoir deux fichiers de test (xxx_internal_test.go pour white-box et xxx_external_test.go pour black-box)",
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
	// Maps pour tracker les fichiers
	sourceFiles := make(map[string]bool)        // xxx.go -> true
	internalTestFiles := make(map[string]bool)  // xxx -> true (a xxx_internal_test.go)
	externalTestFiles := make(map[string]bool)  // xxx -> true (a xxx_external_test.go)

	// Premier passage : collecter tous les fichiers
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		baseName := filepath.Base(filename)

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Extraire le nom de base du fichier testé
			var testedFileName string
			// Vérification du type de test
			if strings.HasSuffix(baseName, "_internal_test.go") {
				testedFileName = strings.TrimSuffix(baseName, "_internal_test.go")
				internalTestFiles[testedFileName] = true
			} else if strings.HasSuffix(baseName, "_external_test.go") {
				testedFileName = strings.TrimSuffix(baseName, "_external_test.go")
				externalTestFiles[testedFileName] = true
			}
			// Continuer à la prochaine itération
			continue
		}

		// Fichier source Go
		fileBase := strings.TrimSuffix(baseName, ".go")
		sourceFiles[fileBase] = true
	}

	// Deuxième passage : vérifier que chaque fichier source a ses deux fichiers de test
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		baseName := filepath.Base(filename)

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Ignorer les fichiers de test
			continue
		}

		fileBase := strings.TrimSuffix(baseName, ".go")

		// Vérifier si le fichier internal existe
		hasInternal := internalTestFiles[fileBase]
		// Vérifier si le fichier external existe
		hasExternal := externalTestFiles[fileBase]

		// Vérification des fichiers manquants
		if !hasInternal && !hasExternal {
			pass.Reportf(
				file.Name.Pos(),
				"KTN-TEST-008: le fichier '%s' doit avoir DEUX fichiers de test : '%s_internal_test.go' (white-box: teste fonctions privées) ET '%s_external_test.go' (black-box: teste API publique)",
				baseName,
				fileBase,
				fileBase,
			)
		} else if !hasInternal {
			pass.Reportf(
				file.Name.Pos(),
				"KTN-TEST-008: le fichier '%s' n'a pas de fichier '%s_internal_test.go' (white-box: pour tester les fonctions privées non-exportées)",
				baseName,
				fileBase,
			)
		} else if !hasExternal {
			pass.Reportf(
				file.Name.Pos(),
				"KTN-TEST-008: le fichier '%s' n'a pas de fichier '%s_external_test.go' (black-box: pour tester l'API publique exportée)",
				baseName,
				fileBase,
			)
		}
	}

	// Retour de la fonction
	return nil, nil
}
