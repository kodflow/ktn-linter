// Package testhelper implements KTN linter rules.
package testhelper

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	// INITIAL_TYPE_MAP_CAP définit la capacité initiale pour Types, Defs, Uses
	INITIAL_TYPE_MAP_CAP int = 64
	// INITIAL_SELECTOR_MAP_CAP définit la capacité initiale pour Implicits, Selections, Scopes
	INITIAL_SELECTOR_MAP_CAP int = 16
	// INITIAL_ANALYZER_MAP_CAP définit la capacité initiale pour ResultOf
	INITIAL_ANALYZER_MAP_CAP int = 8
)

// createTypeInfo crée une nouvelle structure types.Info pour le type checking.
//
// Returns:
//   - *types.Info: structure d'information de types
func createTypeInfo() *types.Info {
	// Retour de la structure types.Info avec capacités initiales estimées
	return &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue, INITIAL_TYPE_MAP_CAP),
		Defs:       make(map[*ast.Ident]types.Object, INITIAL_TYPE_MAP_CAP),
		Uses:       make(map[*ast.Ident]types.Object, INITIAL_TYPE_MAP_CAP),
		Implicits:  make(map[ast.Node]types.Object, INITIAL_SELECTOR_MAP_CAP),
		Selections: make(map[*ast.SelectorExpr]*types.Selection, INITIAL_SELECTOR_MAP_CAP),
		Scopes:     make(map[ast.Node]*types.Scope, INITIAL_SELECTOR_MAP_CAP),
	}
}

// createPass crée un analysis.Pass pour l'exécution d'un analyzer.
//
// Params:
//   - fset: ensemble de fichiers
//   - file: fichier AST à analyser
//   - pkg: package typé
//   - info: informations de types
//   - diagnostics: pointeur vers la liste des diagnostics
//
// Returns:
//   - *analysis.Pass: contexte d'analyse créé
func createPass(fset *token.FileSet, file *ast.File, pkg *types.Package, info *types.Info, diagnostics *[]analysis.Diagnostic) *analysis.Pass {
	// Retour du pass d'analyse
	return &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any, INITIAL_ANALYZER_MAP_CAP),
		ReadFile: func(filename string) ([]byte, error) {
			// Lecture du fichier pour les analyzers qui en ont besoin
			return os.ReadFile(filename)
		},
	}
}

// RunAnalyzer exécute un analyzer sur un fichier et retourne les diagnostics.
//
// Params:
//   - t: contexte de test
//   - analyzer: l'analyzer à exécuter
//   - filename: chemin du fichier à analyser
//
// Returns:
//   - []analysis.Diagnostic: liste des diagnostics trouvés
func RunAnalyzer(t TestingT, analyzer *analysis.Analyzer, filename string) []analysis.Diagnostic {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	// Vérification d'erreur de parsing
	if err != nil {
		// Retour de la fonction avec erreur fatale
		t.Fatalf("failed to parse %s: %v", filename, err)
		// Retour anticipé pour les mocks qui ne terminent pas le test
		return nil
	}

	// Création de la configuration de type checking
	conf := &types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // Ignorer les erreurs de type pour les tests
	}
	info := createTypeInfo()
	pkg, _ := conf.Check(file.Name.Name, fset, []*ast.File{file}, info)

	// Initialisation de la liste des diagnostics
	var diagnostics []analysis.Diagnostic
	pass := createPass(fset, file, pkg, info, &diagnostics)

	// Exécution des analyzers requis (comme inspect.Analyzer)
	for _, req := range analyzer.Requires {
		var result any
		result, err = req.Run(pass)
		// Vérification d'erreur de l'analyzer requis
		if err != nil {
			// Retour de la fonction avec erreur fatale
			t.Fatalf("required analyzer %s failed: %v", req.Name, err)
			// Retour anticipé pour les mocks qui ne terminent pas le test
			return nil
		}
		pass.ResultOf[req] = result
	}

	_, err = analyzer.Run(pass)
	// Vérification d'erreur de l'analyzer principal
	if err != nil {
		// Retour de la fonction avec erreur fatale
		t.Fatalf("analyzer failed: %v", err)
		// Retour anticipé pour les mocks qui ne terminent pas le test
		return nil
	}

	// Retour de la fonction avec les diagnostics
	return diagnostics
}

// TestGoodBad teste que good.go a 0 erreurs et bad.go le nombre attendu.
//
// Params:
//   - t: contexte de test
//   - analyzer: l'analyzer à tester
//   - testDir: nom du répertoire de test
//   - expectedBadErrors: nombre d'erreurs attendues dans bad.go
func TestGoodBad(t TestingT, analyzer *analysis.Analyzer, testDir string, expectedBadErrors int) {
	TestGoodBadWithFiles(t, analyzer, testDir, "good.go", "bad.go", expectedBadErrors)
}

// TestGoodBadWithFiles teste avec des noms de fichiers personnalisés.
//
// Params:
//   - t: contexte de test
//   - analyzer: l'analyzer à tester
//   - testDir: nom du répertoire de test
//   - goodFilename: nom du fichier "good"
//   - badFilename: nom du fichier "bad"
//   - expectedBadErrors: nombre d'erreurs attendues dans le fichier bad
func TestGoodBadWithFiles(t TestingT, analyzer *analysis.Analyzer, testDir string, goodFilename string, badFilename string, expectedBadErrors int) {
	goodFile := "testdata/src/" + testDir + "/" + goodFilename
	badFile := "testdata/src/" + testDir + "/" + badFilename

	// Test du fichier good - doit avoir 0 erreur
	diags := RunAnalyzer(t, analyzer, goodFile)
	// Vérification du nombre de diagnostics
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodFile, len(diags))
		// Itération sur les diagnostics pour les afficher
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test du fichier bad - doit avoir le nombre attendu d'erreurs
	diags = RunAnalyzer(t, analyzer, badFile)
	// Vérification du nombre de diagnostics
	if len(diags) != expectedBadErrors {
		t.Errorf("%s should have %d errors, got %d", badFile, expectedBadErrors, len(diags))
		// Itération sur les diagnostics pour les afficher
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}
}

// parsePackageFiles parse tous les fichiers .go d'un répertoire.
//
// Params:
//   - t: contexte de test
//   - dir: répertoire contenant les fichiers
//   - fset: ensemble de fichiers pour le parsing
//
// Returns:
//   - []*ast.File: liste des fichiers parsés
func parsePackageFiles(t TestingT, dir string, fset *token.FileSet) []*ast.File {
	entries, err := os.ReadDir(dir)
	// Vérification d'erreur de lecture du répertoire
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
		// Retour anticipé pour les mocks qui ne terminent pas le test
		return nil
	}

	var files []*ast.File
	// Parcours des fichiers du répertoire
	for _, entry := range entries {
		// Vérification si c'est un fichier .go
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			// Ignorer les répertoires et fichiers non-.go
			continue
		}

		filepath := dir + "/" + entry.Name()
		file, parseErr := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
		// Vérification d'erreur de parsing
		if parseErr != nil {
			t.Fatalf("failed to parse %s: %v", filepath, parseErr)
			// Retour anticipé pour les mocks qui ne terminent pas le test
			return nil
		}
		files = append(files, file)
	}

	// Vérification qu'au moins un fichier a été trouvé
	if len(files) == 0 {
		t.Fatalf("no .go files found in %s", dir)
		// Retour anticipé pour les mocks qui ne terminent pas le test
		return nil
	}

	// Retour de la liste des fichiers parsés
	return files
}

// createPassForPackage crée un analysis.Pass pour un package.
//
// Params:
//   - fset: ensemble de fichiers
//   - files: fichiers AST du package
//   - diagnostics: pointeur vers la liste des diagnostics
//
// Returns:
//   - *analysis.Pass: pass d'analyse créé
func createPassForPackage(fset *token.FileSet, files []*ast.File, diagnostics *[]analysis.Diagnostic) *analysis.Pass {
	conf := &types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // Ignorer les erreurs de type pour les tests
	}
	info := createTypeInfo()
	pkg, _ := conf.Check(files[0].Name.Name, fset, files, info)

	// Retour du pass d'analyse
	return &analysis.Pass{
		Fset:      fset,
		Files:     files,
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any, INITIAL_ANALYZER_MAP_CAP),
		ReadFile: func(filename string) ([]byte, error) {
			// Lecture du fichier pour les analyzers qui en ont besoin
			return os.ReadFile(filename)
		},
	}
}

// RunAnalyzerOnPackage exécute un analyzer sur tous les fichiers d'un package.
//
// Params:
//   - t: contexte de test
//   - analyzer: l'analyzer à exécuter
//   - dir: répertoire contenant les fichiers du package
//
// Returns:
//   - []analysis.Diagnostic: liste des diagnostics trouvés
func RunAnalyzerOnPackage(t TestingT, analyzer *analysis.Analyzer, dir string) []analysis.Diagnostic {
	fset := token.NewFileSet()
	files := parsePackageFiles(t, dir, fset)

	var diagnostics []analysis.Diagnostic
	pass := createPassForPackage(fset, files, &diagnostics)

	// Exécution des analyzers requis
	for _, req := range analyzer.Requires {
		var result any
		var err error
		result, err = req.Run(pass)
		// Vérification d'erreur de l'analyzer requis
		if err != nil {
			t.Fatalf("required analyzer %s failed: %v", req.Name, err)
			// Retour anticipé pour les mocks qui ne terminent pas le test
			return nil
		}
		pass.ResultOf[req] = result
	}

	_, err := analyzer.Run(pass)
	// Vérification d'erreur de l'analyzer principal
	if err != nil {
		t.Fatalf("analyzer failed: %v", err)
		// Retour anticipé pour les mocks qui ne terminent pas le test
		return nil
	}

	// Retour de la liste des diagnostics
	return diagnostics
}

// TestGoodBadPackage teste un analyzer sur des packages complets.
//
// Params:
//   - t: contexte de test
//   - analyzer: l'analyzer à tester
//   - testDir: nom du répertoire de test
//   - expectedBadErrors: nombre d'erreurs attendues dans le package bad
func TestGoodBadPackage(t TestingT, analyzer *analysis.Analyzer, testDir string, expectedBadErrors int) {
	goodDir := "testdata/src/" + testDir + "/good"
	badDir := "testdata/src/" + testDir + "/bad"

	// Test du package good - doit avoir 0 erreur
	diags := RunAnalyzerOnPackage(t, analyzer, goodDir)
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodDir, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test du package bad - doit avoir le nombre attendu d'erreurs
	diags = RunAnalyzerOnPackage(t, analyzer, badDir)
	if len(diags) != expectedBadErrors {
		t.Errorf("%s should have %d errors, got %d", badDir, expectedBadErrors, len(diags))
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}
}
