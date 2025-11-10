package ktntest

import (
	"os"
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
	// Parcourir tous les fichiers sources (non-test)
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Ignorer les fichiers de test
			continue
		}

		// Extraire le répertoire et le nom de base
		dir := filepath.Dir(filename)
		baseName := filepath.Base(filename)
		fileBase := strings.TrimSuffix(baseName, ".go")

		// Construire les chemins attendus pour les fichiers de test
		internalTestPath := filepath.Join(dir, fileBase+"_internal_test.go")
		externalTestPath := filepath.Join(dir, fileBase+"_external_test.go")

		// Vérifier si les fichiers existent sur le disque (détecte aussi XTestGoFiles)
		hasInternal := fileExistsOnDisk(internalTestPath)
		hasExternal := fileExistsOnDisk(externalTestPath)

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

// fileExistsOnDisk vérifie si un fichier existe sur le disque.
// Cette fonction permet de détecter les fichiers _external_test.go même s'ils utilisent
// package xxx_test et sont dans XTestGoFiles (package séparé de Go).
//
// Params:
//   - path: chemin absolu du fichier à vérifier
//
// Returns:
//   - bool: true si le fichier existe, false sinon
func fileExistsOnDisk(path string) bool {
	info, err := os.Stat(path)
	// Vérification de la condition
	if err != nil {
		// Erreur ou fichier n'existe pas
		return false
	}
	// Retour du résultat (fichier existe et n'est pas un répertoire)
	return !info.IsDir()
}
