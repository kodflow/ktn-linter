package analyzer

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// InterfaceStrictAnalyzer vérifie que interfaces.go contient UNIQUEMENT des interfaces.
var InterfaceStrictAnalyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktninterfacestrict",
	Doc:  "Vérifie que interfaces.go contient uniquement des interfaces et que mock.go existe si nécessaire",
	Run:  runInterfaceStrictAnalyzer,
}

// runInterfaceStrictAnalyzer exécute l'analyseur strict pour interfaces.go et mock.go.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runInterfaceStrictAnalyzer(pass *analysis.Pass) (interface{}, error) {
	return runInterfaceStrictAnalyzerWithFS(pass, defaultFS)
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
	// Ignorer les packages exemptés
	if isInterfaceExemptedPackage(pass.Pkg.Name()) {
		return nil, nil
	}

	// Chercher interfaces.go
	var interfacesFile *ast.File
	var interfacesPath string
	
	for _, file := range pass.Files {
		path := pass.Fset.File(file.Pos()).Name()
		if filepath.Base(path) == "interfaces.go" {
			interfacesFile = file
			interfacesPath = path
			break
		}
	}

	// Si pas de interfaces.go, rien à vérifier
	if interfacesFile == nil {
		return nil, nil
	}

	// Vérifier que interfaces.go contient UNIQUEMENT des interfaces
	hasInterfaces := false
	hasNonInterfaces := false
	var firstNonInterface token.Pos

	for _, decl := range interfacesFile.Decls {
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

	// KTN-INTERFACE-008: interfaces.go doit contenir UNIQUEMENT des interfaces
	if hasNonInterfaces {
		pass.Reportf(firstNonInterface,
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

	// Vérifier mock.go
	if hasInterfaces {
		checkMockFileExists(pass, interfacesPath, fs)
	}

	return nil, nil
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
