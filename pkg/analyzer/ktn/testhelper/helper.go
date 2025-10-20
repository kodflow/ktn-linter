package testhelper

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/analysis"
)

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
	conf := types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // Ignorer les erreurs de type pour les tests
	}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	pkg, _ := conf.Check(file.Name.Name, fset, []*ast.File{file}, info)

	// Initialisation de la liste des diagnostics
	var diagnostics []analysis.Diagnostic

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
		ReadFile: func(filename string) ([]byte, error) {
			// Lecture du fichier pour les analyzers qui en ont besoin
			return os.ReadFile(filename)
		},
	}

	// Exécution des analyzers requis (comme inspect.Analyzer)
	for _, req := range analyzer.Requires {
		result, err := req.Run(pass)
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
	goodFile := "testdata/src/" + testDir + "/good.go"
	badFile := "testdata/src/" + testDir + "/bad.go"

	// Test de good.go - doit avoir 0 erreur
	diags := RunAnalyzer(t, analyzer, goodFile)
	// Vérification du nombre de diagnostics
	if len(diags) != 0 {
		t.Errorf("%s should have 0 errors, got %d", goodFile, len(diags))
		// Itération sur les diagnostics pour les afficher
		for _, d := range diags {
			t.Logf("  %v", d.Message)
		}
	}

	// Test de bad.go - doit avoir le nombre attendu d'erreurs
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
