package ktn_test

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule003 analyzer for KTN linter.
var Rule003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_TEST_003",
	Doc:  "Vérifie qu'il n'y a pas de _test.go orphelin",
	Run:  runRule003,
}

func runRule003(pass *analysis.Pass) (any, error) {
	// Si on est dans un package non-test, ne pas vérifier
	if !strings.HasSuffix(pass.Pkg.Name(), "_test") {
		// Analysis completed successfully.
		return nil, nil
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Vérifier uniquement les fichiers _test.go
		if !strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		sourceFileName := strings.TrimSuffix(baseName, "_test.go") + ".go"
		sourcePath := filepath.Join(filepath.Dir(fileName), sourceFileName)

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			pass.Reportf(file.Package,
				"[KTN_TEST_003] Fichier de test '%s' n'a pas de fichier source correspondant.\n"+
					"Le fichier '%s' n'existe pas.\n"+
					"Créez '%s' ou renommez/supprimez ce fichier de test.",
				baseName, sourceFileName, sourceFileName)
		}
	}

	// Analysis completed successfully.
	return nil, nil
}
