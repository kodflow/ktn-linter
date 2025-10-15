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
	// Retourne le résultat de l'exécution
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
	// Retourne le fichier AST trouvé ou nil
	return findASTFile(pass, path)
}

// fileInfoForTest expose le type fileInfo pour les tests (type privé).
type fileInfoForTest struct {
	Path        string
	IsTest      bool
	PackageName string
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
	// Retourne une fileInfo pour utilisation interne
	return &fileInfo{
		path:        fi.Path,
		isTest:      fi.IsTest,
		packageName: fi.PackageName,
		hasTests:    fi.HasTests,
	}
}
