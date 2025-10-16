package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// RunTestAnalyzerWithFS expose runTestAnalyzerWithFS pour les tests.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//   - fs: le système de fichiers à utiliser
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func RunTestAnalyzerWithFS(pass *analysis.Pass, fs filesystem.FileSystem) (interface{}, error) {
	// Délègue à l'implémentation interne
	return runTestAnalyzerWithFS(pass, fs)
}

// FindASTFileForTest expose findASTFile pour les tests.
//
// Params:
//   - pass: la passe d'analyse
//   - path: le chemin du fichier à trouver
//
// Returns:
//   - *ast.File: le fichier AST trouvé ou nil
func FindASTFileForTest(pass *analysis.Pass, path string) *ast.File {
	// Délègue à l'implémentation interne
	return findASTFile(pass, path)
}

// fileInfoForTest expose le type fileInfo pour les tests (type privé).
type fileInfoForTest struct {
	// Path est le chemin du fichier
	Path        string
	// IsTest indique si c'est un fichier de test
	IsTest      bool
	// PackageName est le nom du package
	PackageName string
	// HasTests indique si le fichier contient des tests
	HasTests    bool
}

// ToInternalFileInfo convertit fileInfoForTest en fileInfo interne.
//
// Params:
//   - fi: le fileInfoForTest à convertir
//
// Returns:
//   - *fileInfo: une nouvelle instance de fileInfo avec les mêmes valeurs
func ToInternalFileInfo(fi fileInfoForTest) *fileInfo {
	// Convertit et retourne une instance de fileInfo interne
	return &fileInfo{
		path:        fi.Path,
		isTest:      fi.IsTest,
		packageName: fi.PackageName,
		hasTests:    fi.HasTests,
	}
}

// ContainsOnlyInterfacesForTest expose containsOnlyInterfaces pour les tests.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si le fichier contient uniquement des interfaces
func ContainsOnlyInterfacesForTest(file *ast.File) bool {
	// Délègue à l'implémentation interne
	return containsOnlyInterfaces(file)
}

// IsFunctionDeclForTest expose isFunctionDecl pour les tests.
//
// Params:
//   - decl: la déclaration à vérifier
//
// Returns:
//   - bool: true si c'est une fonction
func IsFunctionDeclForTest(decl ast.Decl) bool {
	// Délègue à l'implémentation interne
	return isFunctionDecl(decl)
}

// IsInterfaceTypeForTest expose isInterfaceType pour les tests.
//
// Params:
//   - typeSpec: la spécification de type
//
// Returns:
//   - bool: true si c'est une interface
func IsInterfaceTypeForTest(typeSpec *ast.TypeSpec) bool {
	// Délègue à l'implémentation interne
	return isInterfaceType(typeSpec)
}

// IsTestableTypeForTest expose isTestableType pour les tests.
//
// Params:
//   - decl: la déclaration à vérifier
//
// Returns:
//   - bool: true si contient struct ou interface
func IsTestableTypeForTest(decl ast.Decl) bool {
	// Délègue à l'implémentation interne
	return isTestableType(decl)
}

// HasTestableElementsForTest expose hasTestableElements pour les tests.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si le fichier contient des éléments testables
func HasTestableElementsForTest(file *ast.File) bool {
	// Délègue à l'implémentation interne
	return hasTestableElements(file)
}

// ShouldSkipTestCoverageForTest expose shouldSkipTestCoverage pour les tests via un wrapper.
//
// Params:
//   - pass: la passe d'analyse
//   - path: le chemin du fichier
//   - isTest: si c'est un fichier de test
//   - pkgName: nom du package
//   - hasTests: si le fichier contient des tests
//
// Returns:
//   - bool: true si le fichier doit être ignoré
func ShouldSkipTestCoverageForTest(pass *analysis.Pass, path string, isTest bool, pkgName string, hasTests bool) bool {
	info := &fileInfo{
		path:        path,
		isTest:      isTest,
		packageName: pkgName,
		hasTests:    hasTests,
	}
	// Délègue à l'implémentation interne avec une fileInfo construite
	return shouldSkipTestCoverage(pass, info)
}
