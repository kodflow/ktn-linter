// Package ktntest implements KTN linter rules.
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
var Analyzer008 = &analysis.Analyzer{
	Name: "ktntest008",
	Doc:  "KTN-TEST-008: Chaque fichier .go doit avoir les fichiers de test appropriés (xxx_internal_test.go si fonctions privées, xxx_external_test.go si fonctions publiques)",
	Run:  runTest008,
}

// fileAnalysisResult contient le résultat de l'analyse d'un fichier.
type fileAnalysisResult struct {
	hasPublic  bool
	hasPrivate bool
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
			continue
		}

		// Ignorer le package main (ne peut pas avoir de tests externes)
		if file.Name.Name == "main" {
			continue
		}

		// Analyser le contenu du fichier
		result := analyzeFileFunctions(file)

		// Si le fichier n'a ni fonctions publiques ni privées, pas de test requis
		if !result.hasPublic && !result.hasPrivate {
			continue
		}

		// Vérifier les fichiers de test existants
		status := checkTestFilesExist(filename)

		// Reporter les problèmes
		reportTestFileIssues(pass, file, result, status)
	}

	// Retour de la fonction
	return nil, nil
}

// analyzeFileFunctions analyse un fichier pour détecter les fonctions/variables publiques/privées.
//
// Params:
//   - file: fichier AST à analyser
//
// Returns:
//   - fileAnalysisResult: résultat de l'analyse
func analyzeFileFunctions(file *ast.File) fileAnalysisResult {
	result := fileAnalysisResult{}

	// Parcourir l'AST du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier les déclarations de fonctions
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name != nil {
			funcName := funcDecl.Name.Name
			// Ignorer les fonctions exemptées (init, main)
			if !isExemptFunction(funcName) {
				// Vérifier si la fonction est publique ou privée
				if len(funcName) > 0 && unicode.IsUpper(rune(funcName[0])) {
					result.hasPublic = true
				} else {
					// Fonction privée (commence par minuscule)
					result.hasPrivate = true
				}
			}
			// Continuer l'itération
			return true
		}

		// Vérifier les déclarations de variables (publiques et privées)
		if genDecl, ok := n.(*ast.GenDecl); ok {
			checkVariables(genDecl, &result)
		}

		// Continuer l'itération
		return true
	})

	// Retour du résultat d'analyse
	return result
}

// checkVariables vérifie si une déclaration contient des variables publiques ou privées.
//
// Params:
//   - genDecl: déclaration générique
//   - result: résultat de l'analyse à mettre à jour
func checkVariables(genDecl *ast.GenDecl, result *fileAnalysisResult) {
	// Parcourir les spécifications de la déclaration
	for _, spec := range genDecl.Specs {
		// Vérifier si c'est une spécification de valeur (var ou const)
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Si ce n'est pas une valeur, continuer
		if !ok {
			continue
		}
		// Vérifier chaque nom dans la spécification
		for _, name := range valueSpec.Names {
			varName := name.Name
			// Ignorer les variables blank (_)
			if varName == "_" {
				continue
			}
			// Classifier la variable comme publique ou privée
			if len(varName) > 0 && unicode.IsUpper(rune(varName[0])) {
				// Variable publique (commence par majuscule)
				result.hasPublic = true
			} else {
				// Variable privée (commence par minuscule)
				result.hasPrivate = true
			}
		}
	}
}

// checkTestFilesExist vérifie l'existence des fichiers de test pour un fichier source.
//
// Params:
//   - filename: chemin du fichier source
//
// Returns:
//   - testFilesStatus: état des fichiers de test
func checkTestFilesExist(filename string) testFilesStatus {
	dir := filepath.Dir(filename)
	baseName := filepath.Base(filename)
	fileBase := strings.TrimSuffix(baseName, ".go")

	// Retour de l'état des fichiers de test
	return testFilesStatus{
		baseName:    baseName,
		fileBase:    fileBase,
		hasInternal: fileExistsOnDisk(filepath.Join(dir, fileBase+"_internal_test.go")),
		hasExternal: fileExistsOnDisk(filepath.Join(dir, fileBase+"_external_test.go")),
	}
}

// reportTestFileIssues reporte les problèmes de fichiers de test manquants ou superflus.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - result: résultat de l'analyse des fonctions
//   - status: état des fichiers de test
func reportTestFileIssues(pass *analysis.Pass, file *ast.File, result fileAnalysisResult, status testFilesStatus) {
	// Vérification selon le type de fichier
	switch {
	// Cas fichier mixte avec fonctions publiques ET privées
	case result.hasPublic && result.hasPrivate:
		reportMixedFunctionsIssues(pass, file, status)
	// Cas fichier avec fonctions publiques uniquement
	case result.hasPublic:
		reportPublicOnlyIssues(pass, file, status)
	// Cas fichier avec fonctions privées uniquement
	case result.hasPrivate:
		reportPrivateOnlyIssues(pass, file, status)
	}
}

// reportMixedFunctionsIssues reporte les problèmes pour fichiers avec fonctions publiques ET privées.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - status: état des fichiers de test
func reportMixedFunctionsIssues(pass *analysis.Pass, file *ast.File, status testFilesStatus) {
	// Fichier avec fonctions publiques ET privées → besoin des deux fichiers
	if !status.hasInternal && !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient des fonctions publiques ET privées. Il doit avoir DEUX fichiers de test : '%s_internal_test.go' (white-box) ET '%s_external_test.go' (black-box)",
			status.baseName, status.fileBase, status.fileBase)
		return
	}

	// Vérifier les fichiers manquants individuellement
	if !status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient des fonctions privées. Il doit avoir un fichier '%s_internal_test.go' (white-box)",
			status.baseName, status.fileBase)
	}
	// Vérification du fichier externe
	if !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient des fonctions publiques. Il doit avoir un fichier '%s_external_test.go' (black-box)",
			status.baseName, status.fileBase)
	}
}

// reportPublicOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions publiques.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - status: état des fichiers de test
func reportPublicOnlyIssues(pass *analysis.Pass, file *ast.File, status testFilesStatus) {
	// Vérification du fichier externe manquant
	if !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient des fonctions publiques. Il doit avoir un fichier '%s_external_test.go' (black-box)",
			status.baseName, status.fileBase)
		return
	}

	// Vérification du fichier interne superflu
	if status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient UNIQUEMENT des fonctions publiques. Le fichier '%s_internal_test.go' est inutile et doit être supprimé (utilisez '%s_external_test.go' pour tester l'API publique)",
			status.baseName, status.fileBase, status.fileBase)
	}
}

// reportPrivateOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions privées.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - status: état des fichiers de test
func reportPrivateOnlyIssues(pass *analysis.Pass, file *ast.File, status testFilesStatus) {
	// Vérification du fichier interne manquant
	if !status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient des fonctions privées. Il doit avoir un fichier '%s_internal_test.go' (white-box)",
			status.baseName, status.fileBase)
		return
	}

	// Vérification du fichier externe superflu
	if status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: le fichier '%s' contient UNIQUEMENT des fonctions privées. Le fichier '%s_external_test.go' est inutile et doit être supprimé (utilisez '%s_internal_test.go' pour tester les fonctions privées)",
			status.baseName, status.fileBase, status.fileBase)
	}
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
