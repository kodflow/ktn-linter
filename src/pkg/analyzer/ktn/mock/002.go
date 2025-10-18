package ktn_mock

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule002 vérifie que chaque interface a un mock correspondant.
//
// KTN-MOCK-002: Chaque interface dans interfaces.go doit avoir un mock correspondant.
// Le mock doit être nommé Mock<InterfaceName> et être dans mock.go.
//
// Incorrect: Interface "Service" sans mock
// Correct: Interface "Service" avec mock "MockService" dans mock.go
var Rule002 = &analysis.Analyzer{
	Name: "KTN_MOCK_002",
	Doc:  "Vérifie que chaque interface a un mock correspondant dans mock.go",
	Run:  runRule002,
}

// runRule002 exécute la vérification KTN-MOCK-002.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule002(pass *analysis.Pass) (any, error) {
	// Ignorer les packages exemptés
	if isExemptedPackage002(pass.Pkg.Name()) {
		return nil, nil
	}

	// Chercher interfaces.go
	interfacesFile, interfacesPath := findInterfacesFile002(pass)
	if interfacesFile == nil {
		return nil, nil
	}

	// Extraire les interfaces avec leurs positions
	interfaces := extractInterfaceNamesWithPos(interfacesFile)
	if len(interfaces) == 0 {
		return nil, nil
	}

	// Chercher mock.go
	mockFile := findMockFile(pass, interfacesPath)
	if mockFile == nil {
		// KTN-MOCK-001 devrait déjà l'avoir rapporté
		return nil, nil
	}

	// Extraire les mocks
	mocks := extractMockNames(mockFile)

	// Vérifier que chaque interface a un mock
	for name, pos := range interfaces {
		expectedMockName := "Mock" + name
		if !contains(mocks, expectedMockName) {
			reportMissingMock(pass, name, expectedMockName, pos)
		}
	}

	return nil, nil
}

// isExemptedPackage002 vérifie si le package est exempté.
//
// Params:
//   - pkgName: le nom du package
//
// Returns:
//   - bool: true si exempté
func isExemptedPackage002(pkgName string) bool {
	return pkgName == "main" || strings.HasSuffix(pkgName, "_test")
}

// findInterfacesFile002 cherche le fichier interfaces.go.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - *ast.File: le fichier interfaces.go ou nil
//   - string: le chemin du fichier
func findInterfacesFile002(pass *analysis.Pass) (*ast.File, string) {
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

// findMockFile cherche le fichier mock.go.
//
// Params:
//   - pass: la passe d'analyse
//   - interfacesPath: le chemin de interfaces.go
//
// Returns:
//   - *ast.File: le fichier mock.go ou nil
func findMockFile(pass *analysis.Pass, interfacesPath string) *ast.File {
	dir := filepath.Dir(interfacesPath)

	for _, file := range pass.Files {
		filePos := pass.Fset.File(file.Pos())
		if filePos == nil {
			continue
		}
		path := filePos.Name()
		if filepath.Dir(path) == dir && filepath.Base(path) == "mock.go" {
			return file
		}
	}
	return nil
}

// extractInterfaceNamesWithPos extrait les noms de toutes les interfaces avec leurs positions.
//
// Params:
//   - file: le fichier AST
//
// Returns:
//   - map[string]token.Pos: map des noms d'interfaces vers leurs positions
func extractInterfaceNamesWithPos(file *ast.File) map[string]token.Pos {
	interfaces := make(map[string]token.Pos)

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
				interfaces[typeSpec.Name.Name] = typeSpec.Pos()
			}
		}
	}

	return interfaces
}

// extractMockNames extrait les noms de tous les mocks (structs commençant par "Mock").
//
// Params:
//   - file: le fichier AST
//
// Returns:
//   - []string: liste des noms de mocks
func extractMockNames(file *ast.File) []string {
	var mocks []string

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

			if isMockStruct(typeSpec) {
				mocks = append(mocks, typeSpec.Name.Name)
			}
		}
	}

	return mocks
}

// isMockStruct vérifie si un TypeSpec est un struct Mock.
//
// Params:
//   - typeSpec: spécification de type
//
// Returns:
//   - bool: true si c'est un struct commençant par "Mock"
func isMockStruct(typeSpec *ast.TypeSpec) bool {
	_, isStruct := typeSpec.Type.(*ast.StructType)
	return isStruct && strings.HasPrefix(typeSpec.Name.Name, "Mock")
}

// contains vérifie si une slice contient une chaîne.
//
// Params:
//   - slice: la slice à vérifier
//   - item: l'élément à chercher
//
// Returns:
//   - bool: true si l'élément est trouvé
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// reportMissingMock rapporte une violation KTN-MOCK-002.
//
// Params:
//   - pass: la passe d'analyse
//   - interfaceName: nom de l'interface sans mock
//   - expectedMockName: nom du mock attendu
//   - pos: position de l'interface
func reportMissingMock(pass *analysis.Pass, interfaceName, expectedMockName string, pos token.Pos) {
	pass.Reportf(pos,
		"[KTN-MOCK-002] L'interface '%s' n'a pas de mock correspondant dans 'mock.go'.\n"+
			"Créez le mock '%s' dans 'mock.go'.\n"+
			"Exemple:\n"+
			"  // %s est le mock de %s.\n"+
			"  type %s struct {\n"+
			"      // Champs pour simuler le comportement\n"+
			"  }\n"+
			"\n"+
			"  // Implémentez toutes les méthodes de l'interface %s\n",
		interfaceName, expectedMockName, expectedMockName, interfaceName, expectedMockName, interfaceName)
}
