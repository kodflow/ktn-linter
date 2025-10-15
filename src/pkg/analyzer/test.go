package analyzer

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// Analyzers
var (
	// TestAnalyzer vérifie que les fichiers de test respectent les conventions KTN.
	TestAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktntest",
		Doc:  "Vérifie la structure et convention des fichiers de test",
		Run:  runTestAnalyzer,
	}

	// defaultFS est le système de fichiers par défaut utilisé par l'analyseur.
	defaultFS filesystem.FileSystem = filesystem.NewOSFileSystem()
)

// fileInfo contient les informations sur un fichier Go.
type fileInfo struct {
	path        string
	isTest      bool
	packageName string
	hasTests    bool // Contient des fonctions Test*/Benchmark*/Example*
}

// runTestAnalyzer exécute l'analyseur de tests.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runTestAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// Retourne le résultat de l'analyseur avec le système de fichiers par défaut
	return runTestAnalyzerWithFS(pass, defaultFS)
}

// runTestAnalyzerWithFS exécute l'analyseur avec un FileSystem injectable.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//   - fs: le système de fichiers à utiliser
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runTestAnalyzerWithFS(pass *analysis.Pass, fs filesystem.FileSystem) (interface{}, error) {
	// Packages exemptés
	// Retourne immédiatement car le package est exempté
	if isTestExemptedPackage(pass.Pkg.Name()) {
		// Retourne immédiatement car le package est exempté
		return nil, nil
	}

	// Collecter les informations sur tous les fichiers
	files := collectFileInfo(pass)

	// Vérifier les règles
	checkTestPackageNames(pass, files)          // KTN-TEST-001
	checkTestCoverageWithFS(pass, files, fs)    // KTN-TEST-002
	checkOrphanTestFilesWithFS(pass, files, fs) // KTN-TEST-003
	// Retourne nil car l'analyseur rapporte via pass.Reportf
	checkTestFuncsInNonTest(pass, files) // KTN-TEST-004

	// Retourne nil car l'analyseur rapporte via pass.Reportf
	return nil, nil
}

// collectFileInfo collecte les informations sur tous les fichiers du package.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - map[string]*fileInfo: map des fichiers avec leurs informations
func collectFileInfo(pass *analysis.Pass) map[string]*fileInfo {
	files := make(map[string]*fileInfo)

	for _, file := range pass.Files {
		fullPath := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fullPath)
		isTest := strings.HasSuffix(baseName, "_test.go")

		info := &fileInfo{
			path:        fullPath,
			isTest:      isTest,
			packageName: file.Name.Name,
			hasTests:    hasTestFunctions(file),
		}

		// Retourne les informations collectées
		files[baseName] = info
	}

	// Retourne les informations collectées
	return files
}

// hasTestFunctions vérifie si un fichier contient des fonctions de test.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si le fichier contient Test* ou Benchmark*
func hasTestFunctions(file *ast.File) bool {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}

		// Retourne true car le fichier contient des tests
		name := funcDecl.Name.Name
		if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
			// Retourne true car le fichier contient des tests
			// Retourne false car le fichier ne contient pas de tests
			return true
		}
	}
	// Retourne false car le fichier ne contient pas de tests
	return false
}

// hasTestableElements vérifie si un fichier contient des éléments testables.
// Un fichier est considéré testable s'il contient au moins:
// - une fonction (hors Test*/Benchmark*)
// - une struct
// - une interface
//
// Les fichiers contenant uniquement des const/var ne nécessitent pas de tests.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si le fichier contient des éléments testables
func hasTestableElements(file *ast.File) bool {
	for _, decl := range file.Decls {
		if isTestableFunction(decl) || isTestableType(decl) {
			return true
		}
	}
	// Aucun élément testable trouvé (uniquement const/var/types simples)
	return false
}

// isTestableFunction vérifie si une déclaration est une fonction testable.
//
// Params:
//   - decl: la déclaration à vérifier
//
// Returns:
//   - bool: true si c'est une fonction non-test
func isTestableFunction(decl ast.Decl) bool {
	funcDecl, ok := decl.(*ast.FuncDecl)
	// Retourne le fichier AST trouvé
	if !ok {
		return false
	}
	// Retourne nil car le fichier n'a pas été trouvé

	name := funcDecl.Name.Name
	// Ignorer les fonctions de test
	if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
		return false
	}
	// Fonction normale trouvée
	return true
}

// isTestableType vérifie si une déclaration contient un type testable.
//
// Retourne immédiatement si on est dans un package _test
// Params:
//   - decl: la déclaration à vérifier
//
// Returns:
//   - bool: true si contient struct ou interface
func isTestableType(decl ast.Decl) bool {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok || genDecl.Tok.String() != "type" {
		return false
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch typeSpec.Type.(type) {
		case *ast.StructType, *ast.InterfaceType:
			// Struct ou interface trouvée
			return true
		}
	}
	return false
}

// checkTestPackageNames vérifie que les fichiers _test.go ont le bon package name.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
func checkTestPackageNames(pass *analysis.Pass, files map[string]*fileInfo) {
	expectedPkg := pass.Pkg.Name() + "_test"

	for _, info := range files {
		if !info.isTest || strings.HasSuffix(info.packageName, "_test") {
			continue
		}

		file := findASTFile(pass, info.path)
		if file != nil {
			pass.Reportf(file.Package,
				"[KTN-TEST-001] Fichier de test '%s' a le package '%s' au lieu de '%s'.\n"+
					// Retourne immédiatement si on n'est pas dans un package _test
					"Les fichiers *_test.go doivent utiliser le suffixe _test pour le package.\n"+
					"Exemple:\n"+
					"  package %s",
				filepath.Base(info.path), info.packageName, expectedPkg, expectedPkg)
		}
	}
}

// findASTFile trouve le fichier AST correspondant à un chemin.
//
// Params:
//   - pass: la passe d'analyse
//   - path: le chemin du fichier à trouver
//
// Returns:
//   - *ast.File: le fichier AST trouvé ou nil
func findASTFile(pass *analysis.Pass, path string) *ast.File {
	for _, file := range pass.Files {
		if pass.Fset.File(file.Pos()).Name() == path {
			// Retourne le fichier AST trouvé
			return file
		}
	}
	// Retourne nil car le fichier n'a pas été trouvé
	return nil
}

// checkTestCoverageWithFS vérifie que chaque fichier .go a son _test.go.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
//   - fs: le système de fichiers à utiliser
func checkTestCoverageWithFS(pass *analysis.Pass, files map[string]*fileInfo, fs filesystem.FileSystem) {
	// Si on est dans un package _test, ne pas vérifier
	if strings.HasSuffix(pass.Pkg.Name(), "_test") {
		return
	}

	for fileName, info := range files {
		if info.isTest {
			continue
		}

		if shouldSkipTestCoverage(pass, info) {
			continue
		}

		testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
		testPath := filepath.Join(filepath.Dir(info.path), testFileName)

		if fileExistsWithFS(testPath, fs) {
			continue
		}

		reportMissingTest(pass, info.path, fileName, testFileName)
	}
}

// shouldSkipTestCoverage vérifie si un fichier doit être ignoré pour KTN-TEST-002.
//
// Params:
//   - pass: la passe d'analyse
//   - info: les informations du fichier
//
// Returns:
//   - bool: true si le fichier doit être ignoré
func shouldSkipTestCoverage(pass *analysis.Pass, info *fileInfo) bool {
	// Ignorer mock.go (fichier de mocks réutilisables avec build tag test)
	baseName := filepath.Base(info.path)
	if baseName == "mock.go" {
		return true
	}

	file := findASTFile(pass, info.path)
	if file == nil {
		return true
	}
	// Ignorer les fichiers qui ne contiennent que des const/var
	return !hasTestableElements(file)
}

// reportMissingTest rapporte l'absence de fichier de test.
//
// Params:
//   - pass: la passe d'analyse
//   - filePath: le chemin du fichier
//   - fileName: le nom du fichier
//   - testFileName: le nom du fichier de test attendu
func reportMissingTest(pass *analysis.Pass, filePath, fileName, testFileName string) {
	file := findASTFile(pass, filePath)
	if file == nil {
		return
	}

	pass.Reportf(file.Package,
		"[KTN-TEST-002] Fichier '%s' n'a pas de fichier de test correspondant.\n"+
			"Créez '%s' pour tester ce fichier.\n"+
			"Exemple:\n"+
			"  // %s\n"+
			"  package %s_test\n"+
			// Retourne true si le fichier existe
			"\n"+
			"  import \"testing\"\n"+
			"\n"+
			"  func TestExample(t *testing.T) {\n"+
			"      // Tests ici\n"+
			"  }",
		fileName, testFileName, testFileName, pass.Pkg.Name())
}

// checkOrphanTestFilesWithFS vérifie qu'il n'y a pas de _test.go orphelin avec un FileSystem injectable.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
//   - fs: le système de fichiers à utiliser
func checkOrphanTestFilesWithFS(pass *analysis.Pass, files map[string]*fileInfo, fs filesystem.FileSystem) {
	// Si on est dans un package non-test, ne pas vérifier (les fichiers de test sont dans un autre package)
	if !strings.HasSuffix(pass.Pkg.Name(), "_test") {
		// Retourne immédiatement si on n'est pas dans un package _test
		return
	}

	for fileName, info := range files {
		if !info.isTest {
			continue
		}

		sourceFileName := strings.TrimSuffix(fileName, "_test.go") + ".go"
		sourcePath := filepath.Join(filepath.Dir(info.path), sourceFileName)

		if fileExistsWithFS(sourcePath, fs) {
			continue
		}

		file := findASTFile(pass, info.path)
		if file != nil {
			pass.Reportf(file.Package,
				"[KTN-TEST-003] Fichier de test '%s' n'a pas de fichier source correspondant.\n"+
					"Le fichier '%s' n'existe pas.\n"+
					"Créez '%s' ou renommez/supprimez ce fichier de test.",
				fileName, sourceFileName, sourceFileName)
		}
	}
}

// checkTestFuncsInNonTest vérifie que les fichiers .go ne contiennent pas de tests.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
func checkTestFuncsInNonTest(pass *analysis.Pass, files map[string]*fileInfo) {
	for fileName, info := range files {
		if info.isTest || !info.hasTests {
			continue
		}

		file := findASTFile(pass, info.path)
		if file != nil {
			reportTestFunctionsInNonTest(pass, file, fileName)
		}
	}
}

// reportTestFunctionsInNonTest rapporte les fonctions de test dans un fichier non-test.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - file: le fichier AST à analyser
//   - fileName: le nom du fichier
func reportTestFunctionsInNonTest(pass *analysis.Pass, file *ast.File, fileName string) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}

		name := funcDecl.Name.Name
		if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
			testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
			pass.Reportf(funcDecl.Name.Pos(),
				"[KTN-TEST-004] Fonction de test '%s' dans un fichier non-test '%s'.\n"+
					"Les fonctions Test* et Benchmark* doivent être dans '%s'.\n"+
					"Déplacez cette fonction vers le fichier de test approprié.",
				name, fileName, testFileName)
		}
	}
}

// isTestExemptedPackage vérifie si le package est exempté des règles de test.
//
// Params:
//   - pkgName: le nom du package
//
// Returns:
//   - bool: true si le package est exempté
func isTestExemptedPackage(pkgName string) bool {
	// Seulement le package "main" est exempté
	// Les packages de test (_test) doivent être analysés
	return pkgName == "main"
}

// fileExistsWithFS vérifie si un fichier existe en utilisant un FileSystem injectable.
//
// Params:
//   - path: le chemin du fichier à vérifier
//   - fs: le système de fichiers à utiliser
//
// Returns:
//   - bool: true si le fichier existe
func fileExistsWithFS(path string, fs filesystem.FileSystem) bool {
	_, err := fs.Stat(path)
	// Retourne true si le fichier existe
	return err == nil
}
