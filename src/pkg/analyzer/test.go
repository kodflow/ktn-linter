package analyzer

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// TestAnalyzer vérifie que les fichiers de test respectent les conventions KTN.
	TestAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktntest",
		Doc:  "Vérifie la structure et convention des fichiers de test",
		Run:  runTestAnalyzer,
	}
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
	// Packages exemptés
	if isTestExemptedPackage(pass.Pkg.Name()) {
		return nil, nil
	}

	// Collecter les informations sur tous les fichiers
	files := collectFileInfo(pass)

	// Vérifier les règles
	checkTestPackageNames(pass, files)    // KTN-TEST-001
	checkTestCoverage(pass, files)        // KTN-TEST-002
	checkOrphanTestFiles(pass, files)     // KTN-TEST-003
	checkTestFuncsInNonTest(pass, files)  // KTN-TEST-004

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

		files[baseName] = info
	}

	return files
}

// hasTestFunctions vérifie si un fichier contient des fonctions de test.
//
// Params:
//   - file: le fichier AST à analyser
//
// Returns:
//   - bool: true si le fichier contient Test*, Benchmark* ou Example*
func hasTestFunctions(file *ast.File) bool {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}

		name := funcDecl.Name.Name
		if strings.HasPrefix(name, "Test") ||
			strings.HasPrefix(name, "Benchmark") ||
			strings.HasPrefix(name, "Example") {
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
			return file
		}
	}
	return nil
}

// checkTestCoverage vérifie que chaque fichier .go a son _test.go.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
func checkTestCoverage(pass *analysis.Pass, files map[string]*fileInfo) {
	// Si on est dans un package _test, ne pas vérifier (les fichiers sont dans un autre package)
	if strings.HasSuffix(pass.Pkg.Name(), "_test") {
		return
	}

	for fileName, info := range files {
		if info.isTest {
			continue
		}

		testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
		testPath := filepath.Join(filepath.Dir(info.path), testFileName)

		if fileExists(testPath) {
			continue
		}

		file := findASTFile(pass, info.path)
		if file != nil {
			pass.Reportf(file.Package,
				"[KTN-TEST-002] Fichier '%s' n'a pas de fichier de test correspondant.\n"+
					"Créez '%s' pour tester ce fichier.\n"+
					"Exemple:\n"+
					"  // %s\n"+
					"  package %s_test\n"+
					"\n"+
					"  import \"testing\"\n"+
					"\n"+
					"  func TestExample(t *testing.T) {\n"+
					"      // Tests ici\n"+
					"  }",
				fileName, testFileName, testFileName, pass.Pkg.Name())
		}
	}
}

// checkOrphanTestFiles vérifie qu'il n'y a pas de _test.go orphelin.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - files: les informations sur tous les fichiers
func checkOrphanTestFiles(pass *analysis.Pass, files map[string]*fileInfo) {
	// Si on est dans un package non-test, ne pas vérifier (les fichiers de test sont dans un autre package)
	if !strings.HasSuffix(pass.Pkg.Name(), "_test") {
		return
	}

	for fileName, info := range files {
		if !info.isTest {
			continue
		}

		sourceFileName := strings.TrimSuffix(fileName, "_test.go") + ".go"
		sourcePath := filepath.Join(filepath.Dir(info.path), sourceFileName)

		if fileExists(sourcePath) {
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
		if strings.HasPrefix(name, "Test") ||
			strings.HasPrefix(name, "Benchmark") ||
			strings.HasPrefix(name, "Example") {

			testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
			pass.Reportf(funcDecl.Name.Pos(),
				"[KTN-TEST-004] Fonction de test '%s' dans un fichier non-test '%s'.\n"+
					"Les fonctions Test*, Benchmark* et Example* doivent être dans '%s'.\n"+
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

// fileExists vérifie si un fichier existe sur le disque.
//
// Params:
//   - path: le chemin du fichier à vérifier
//
// Returns:
//   - bool: true si le fichier existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
