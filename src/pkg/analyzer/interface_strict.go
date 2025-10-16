package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// Analyzers
var (
	// InterfaceStrictAnalyzer vérifie que interfaces.go contient UNIQUEMENT des interfaces.
	InterfaceStrictAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktninterfacestrict",
		Doc:  "Vérifie que interfaces.go contient uniquement des interfaces et que mock.go existe si nécessaire",
		Run:  runInterfaceStrictAnalyzer,
	}

	// interfaceDefaultFS est le système de fichiers par défaut pour interface strict.
	interfaceDefaultFS filesystem.FileSystem = filesystem.NewOSFileSystem()
)

// runInterfaceStrictAnalyzer exécute l'analyseur strict pour interfaces.go et mock.go.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runInterfaceStrictAnalyzer(pass *analysis.Pass) (interface{}, error) {
	return runInterfaceStrictAnalyzerWithFS(pass, interfaceDefaultFS)
}

// runInterfaceStrictAnalyzerWithFS exécute l'analyseur avec un FileSystem injectable.
//
// Params:
//   - pass: la passe d'analyse
//   - fs: le système de fichiers à utiliser
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runInterfaceStrictAnalyzerWithFS(pass *analysis.Pass, fs filesystem.FileSystem) (interface{}, error) {
	if isInterfaceExemptedPackage(pass.Pkg.Name()) {
		return nil, nil
	}

	interfacesFile, interfacesPath := findInterfacesFile(pass)
	if interfacesFile == nil {
		return nil, nil
	}

	hasInterfaces, hasNonInterfaces, firstNonInterface := analyzeInterfacesFile(interfacesFile)

	if hasNonInterfaces {
		reportNonInterfaceTypes(pass, firstNonInterface)
	}

	if hasInterfaces && fs != nil {
		checkMockFileExists(pass, interfacesPath, fs)
		// Vérifier que mock.go contient des mocks pour toutes les interfaces
		checkMockCompleteness(pass, interfacesFile, interfacesPath, fs)
	}

	return nil, nil
}

// findInterfacesFile cherche le fichier interfaces.go dans le package.
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

// analyzeInterfacesFile analyse le contenu de interfaces.go.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si contient des interfaces
//   - bool: true si contient des non-interfaces (structs, etc.)
//   - token.Pos: position du premier type non-interface
func analyzeInterfacesFile(file *ast.File) (bool, bool, token.Pos) {
	hasInterfaces := false
	hasNonInterfaces := false
	var firstNonInterface token.Pos

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

			switch typeSpec.Type.(type) {
			case *ast.InterfaceType:
				hasInterfaces = true
			case *ast.StructType:
				if !hasNonInterfaces {
					firstNonInterface = typeSpec.Pos()
				}
				hasNonInterfaces = true
			}
		}
	}

	return hasInterfaces, hasNonInterfaces, firstNonInterface
}

// reportNonInterfaceTypes rapporte la violation KTN-INTERFACE-008.
//
// Params:
//   - pass: la passe d'analyse
//   - pos: la position du premier type non-interface
func reportNonInterfaceTypes(pass *analysis.Pass, pos token.Pos) {
	pass.Reportf(pos,
		"[KTN-INTERFACE-008] Le fichier 'interfaces.go' contient des types non-interface (struct, alias, etc.).\n"+
			"Ce fichier doit contenir UNIQUEMENT des interfaces.\n"+
			"Déplacez les implémentations (struct) vers impl.go ou un autre fichier.\n"+
			"Exemple:\n"+
			"  // interfaces.go - UNIQUEMENT interfaces\n"+
			"  type Service interface {\n"+
			"      Process() error\n"+
			"  }\n"+
			"\n"+
			"  // impl.go - implémentations\n"+
			"  type service struct {}\n"+
			"  func (s *service) Process() error { return nil }")
}

// checkMockFileExists vérifie que mock.go existe si interfaces.go contient des interfaces.
//
// Params:
//   - pass: la passe d'analyse
//   - interfacesPath: le chemin de interfaces.go
//   - fs: le système de fichiers
func checkMockFileExists(pass *analysis.Pass, interfacesPath string, fs filesystem.FileSystem) {
	dir := filepath.Dir(interfacesPath)
	mockPath := filepath.Join(dir, "mock.go")

	_, err := fs.Stat(mockPath)
	if err != nil {
		// mock.go n'existe pas
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

// checkMockCompleteness vérifie que mock.go contient des mocks pour toutes les interfaces.
//
// Params:
//   - pass: la passe d'analyse
//   - interfacesFile: le fichier interfaces.go
//   - interfacesPath: le chemin de interfaces.go
//   - fs: le système de fichiers
func checkMockCompleteness(pass *analysis.Pass, interfacesFile *ast.File, interfacesPath string, fs filesystem.FileSystem) {
	dir := filepath.Dir(interfacesPath)
	mockPath := filepath.Join(dir, "mock.go")

	// Vérifier que mock.go existe
	_, err := fs.Stat(mockPath)
	if err != nil {
		return // KTN-MOCK-001 déjà rapporté par checkMockFileExists
	}

	// Collecter toutes les interfaces de interfaces.go
	interfaces := extractInterfaceNames(interfacesFile)
	if len(interfaces) == 0 {
		return
	}

	// Parser le fichier mock.go depuis le filesystem
	mockFile, err := parseMockFile(mockPath)
	if err != nil {
		return // Erreur de parsing, on ne rapporte pas
	}

	// Collecter tous les mocks (structs commençant par "Mock")
	mocks := extractMockNames(mockFile)

	// Vérifier que chaque interface a un mock
	for _, interfaceName := range interfaces {
		expectedMockName := "Mock" + interfaceName
		if !contains(mocks, expectedMockName) {
			reportMissingMock(pass, interfaceName, expectedMockName)
		}
	}
}

// extractInterfaceNames extrait les noms de toutes les interfaces d'un fichier.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - []string: liste des noms d'interfaces
func extractInterfaceNames(file *ast.File) []string {
	var interfaces []string

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
				interfaces = append(interfaces, typeSpec.Name.Name)
			}
		}
	}

	return interfaces
}

// parseMockFile parse le fichier mock.go depuis le filesystem.
//
// Params:
//   - mockPath: chemin vers mock.go
//
// Returns:
//   - *ast.File: le fichier AST parsé
//   - error: erreur de parsing
func parseMockFile(mockPath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, mockPath, nil, parser.ParseComments)
}

// extractMockNames extrait les noms de tous les mocks (structs commençant par "Mock").
//
// Params:
//   - file: le fichier AST à analyser
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

		mocks = extractMocksFromGenDecl(genDecl, mocks)
	}

	return mocks
}

// extractMocksFromGenDecl extrait les mocks d'une déclaration générique.
//
// Params:
//   - genDecl: déclaration générique
//   - mocks: liste de mocks existante
//
// Returns:
//   - []string: liste mise à jour des mocks
func extractMocksFromGenDecl(genDecl *ast.GenDecl, mocks []string) []string {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		if isMockStruct(typeSpec) {
			mocks = append(mocks, typeSpec.Name.Name)
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

// reportMissingMock rapporte la violation KTN-MOCK-002.
//
// Params:
//   - pass: la passe d'analyse
//   - interfaceName: nom de l'interface sans mock
//   - expectedMockName: nom du mock attendu
func reportMissingMock(pass *analysis.Pass, interfaceName, expectedMockName string) {
	pass.Reportf(token.Pos(1),
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

// isInterfaceExemptedPackage vérifie si le package est exempté des règles interface strictes.
//
// Params:
//   - pkgName: le nom du package
//
// Returns:
//   - bool: true si le package est exempté
func isInterfaceExemptedPackage(pkgName string) bool {
	return pkgName == "main" || strings.HasSuffix(pkgName, "_test")
}
