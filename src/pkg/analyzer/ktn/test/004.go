package ktn_test

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule004 analyzer for KTN linter.
var Rule004 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_TEST_004",
	Doc:  "Vérifie que les fichiers .go ne contiennent pas de tests",
	Run:  runRule004,
}

func runRule004(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Ignorer les fichiers _test.go
		if strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Recv != nil {
				continue
			}

			name := funcDecl.Name.Name
			if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
				testFileName := strings.TrimSuffix(baseName, ".go") + "_test.go"
				pass.Reportf(funcDecl.Name.Pos(),
					"[KTN_TEST_004] Fonction de test '%s' dans un fichier non-test '%s'.\n"+
						"Les fonctions Test* et Benchmark* doivent être dans '%s'.\n"+
						"Déplacez cette fonction vers le fichier de test approprié.",
					name, baseName, testFileName)
			}
		}
	}

	// Analysis completed successfully.
	return nil, nil
}
