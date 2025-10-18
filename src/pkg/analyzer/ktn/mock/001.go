package ktn_mock

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que mock.go existe si interfaces.go contient des interfaces.
//
// KTN-MOCK-001: Si interfaces.go contient des interfaces, mock.go doit exister.
// Le fichier mock.go doit contenir les mocks réutilisables pour les tests.
//
// Incorrect: interfaces.go avec interfaces mais pas de mock.go
// Correct: interfaces.go avec interfaces ET mock.go avec mocks
var Rule001 = &analysis.Analyzer{
	Name: "KTN_MOCK_001",
	Doc:  "Vérifie que mock.go existe si interfaces.go contient des interfaces",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-MOCK-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	// Ignorer les packages exemptés
	if isExemptedPackage(pass.Pkg.Name()) {
		return nil, nil
	}

	// Chercher interfaces.go
	interfacesFile, interfacesPath := findInterfacesFile(pass)
	if interfacesFile == nil {
		return nil, nil
	}

	// Vérifier si interfaces.go contient des interfaces
	if !hasInterfaces(interfacesFile) {
		return nil, nil
	}

	// Vérifier que mock.go existe
	checkMockFileExists(pass, interfacesPath)

	return nil, nil
}

// isExemptedPackage vérifie si le package est exempté.
//
// Params:
//   - pkgName: le nom du package
//
// Returns:
//   - bool: true si exempté
func isExemptedPackage(pkgName string) bool {
	return pkgName == "main" || strings.HasSuffix(pkgName, "_test")
}

// findInterfacesFile cherche le fichier interfaces.go.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - *ast.File: le fichier interfaces.go ou nil
//   - string: le chemin du fichier
func findInterfacesFile(pass *analysis.Pass) (*ast.File, string) {
	for _, file := range pass.Files {
		filePos := pass.Fset.File(file.Pos())
		if filePos == nil {
			continue
		}
		path := filePos.Name()
		if filepath.Base(path) == "interfaces.go" {
			return file, path
		}
	}
	return nil, ""
}

// hasInterfaces vérifie si un fichier contient des interfaces.
//
// Params:
//   - file: le fichier AST
//
// Returns:
//   - bool: true si contient des interfaces
func hasInterfaces(file *ast.File) bool {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
				return true
			}
		}
	}
	return false
}

// checkMockFileExists vérifie que mock.go existe.
//
// Params:
//   - pass: la passe d'analyse
//   - interfacesPath: le chemin de interfaces.go
func checkMockFileExists(pass *analysis.Pass, interfacesPath string) {
	// Vérifier dans les fichiers du pass si mock.go existe
	dir := filepath.Dir(interfacesPath)
	hasMockFile := false

	for _, file := range pass.Files {
		filePos := pass.Fset.File(file.Pos())
		if filePos == nil {
			continue
		}
		path := filePos.Name()
		if filepath.Dir(path) == dir && filepath.Base(path) == "mock.go" {
			hasMockFile = true
			break
		}
	}

	if !hasMockFile {
		pass.Reportf(token.Pos(1),
			"[KTN-MOCK-001] Le fichier 'interfaces.go' contient des interfaces mais 'mock.go' n'existe pas.\n"+
				"Créez 'mock.go' avec les mocks réutilisables pour les tests.\n"+
				"Template:\n"+
				"  //go:build test\n"+
				"  // +build test\n"+
				"\n"+
				"  package %s\n"+
				"\n"+
				"  // Mock* structs ici\n",
			pass.Pkg.Name())
	}
}
