package ktntest

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Analyzer008 checks that each source file has appropriate test files based on its content
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest008",
	Doc:  "KTN-TEST-008: Chaque fichier .go doit avoir les fichiers de test appropriés (xxx_internal_test.go si fonctions privées, xxx_external_test.go si fonctions publiques)",
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

		// Analyser le contenu du fichier pour détecter les fonctions
		hasPublicFuncs := false
		hasPrivateFuncs := false

		// Parcourir l'AST du fichier
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			// Si ce n'est pas une fonction, continuer
			if !ok {
				// Continue traversal
				return true
			}

			// Ignorer les fonctions sans nom
			if funcDecl.Name == nil {
				// Continue traversal
				return true
			}

			funcName := funcDecl.Name.Name
			// Ignorer les fonctions exemptées (init, main)
			if isExemptFunction(funcName) {
				// Continue traversal
				return true
			}

			// Vérifier si la fonction est publique ou privée
			if len(funcName) > 0 && unicode.IsUpper(rune(funcName[0])) {
				// Fonction publique
				hasPublicFuncs = true
			} else {
				// Fonction privée
				hasPrivateFuncs = true
			}

			// Continue traversal
			return true
		})

		// Si le fichier n'a ni fonctions publiques ni privées, pas de test requis
		if !hasPublicFuncs && !hasPrivateFuncs {
			// Pas de fonctions, pas de tests requis
			continue
		}

		// Extraire le répertoire et le nom de base
		dir := filepath.Dir(filename)
		baseName := filepath.Base(filename)
		fileBase := strings.TrimSuffix(baseName, ".go")

		// Construire les chemins attendus pour les fichiers de test
		internalTestPath := filepath.Join(dir, fileBase+"_internal_test.go")
		externalTestPath := filepath.Join(dir, fileBase+"_external_test.go")

		// Vérifier si les fichiers existent sur le disque
		hasInternal := fileExistsOnDisk(internalTestPath)
		hasExternal := fileExistsOnDisk(externalTestPath)

		// Vérifier les fichiers manquants selon le contenu
		if hasPublicFuncs && hasPrivateFuncs {
			// Fichier avec fonctions publiques ET privées → besoin des deux fichiers
			if !hasInternal && !hasExternal {
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient des fonctions publiques ET privées. Il doit avoir DEUX fichiers de test : '%s_internal_test.go' (white-box) ET '%s_external_test.go' (black-box)",
					baseName,
					fileBase,
					fileBase,
				)
			} else if !hasInternal {
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient des fonctions privées. Il doit avoir un fichier '%s_internal_test.go' (white-box)",
					baseName,
					fileBase,
				)
			} else if !hasExternal {
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient des fonctions publiques. Il doit avoir un fichier '%s_external_test.go' (black-box)",
					baseName,
					fileBase,
				)
			}
		} else if hasPublicFuncs {
			// Fichier avec UNIQUEMENT des fonctions publiques
			if !hasExternal {
				// Manque _external_test.go
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient des fonctions publiques. Il doit avoir un fichier '%s_external_test.go' (black-box)",
					baseName,
					fileBase,
				)
			} else if hasInternal {
				// A _external (correct) mais aussi _internal (inutile) → demander de supprimer _internal
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient UNIQUEMENT des fonctions publiques. Le fichier '%s_internal_test.go' est inutile et doit être supprimé (utilisez '%s_external_test.go' pour tester l'API publique)",
					baseName,
					fileBase,
					fileBase,
				)
			}
		} else if hasPrivateFuncs {
			// Fichier avec UNIQUEMENT des fonctions privées
			if !hasInternal {
				// Manque _internal_test.go
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient des fonctions privées. Il doit avoir un fichier '%s_internal_test.go' (white-box)",
					baseName,
					fileBase,
				)
			} else if hasExternal {
				// A _internal (correct) mais aussi _external (inutile) → demander de supprimer _external
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-008: le fichier '%s' contient UNIQUEMENT des fonctions privées. Le fichier '%s_external_test.go' est inutile et doit être supprimé (utilisez '%s_internal_test.go' pour tester les fonctions privées)",
					baseName,
					fileBase,
					fileBase,
				)
			}
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
