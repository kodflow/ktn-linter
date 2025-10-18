package ktn_test

import (
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Rule001 = &analysis.Analyzer{
	Name: "KTN_TEST_001",
	Doc:  "Vérifie que les fichiers _test.go ont le bon package name",
	Run:  runRule001,
}

func runRule001(pass *analysis.Pass) (any, error) {
	// Package main exempté
	if pass.Pkg.Name() == "main" {
		return nil, nil
	}

	expectedPkg := pass.Pkg.Name() + "_test"

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Vérifier uniquement les fichiers _test.go
		if !strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		// Vérifier si le package a déjà le suffixe _test
		if strings.HasSuffix(file.Name.Name, "_test") {
			continue
		}

		pass.Reportf(file.Package,
			"[KTN_TEST_001] Fichier de test '%s' a le package '%s' au lieu de '%s'.\n"+
				"Les fichiers *_test.go doivent utiliser le suffixe _test pour le package.\n"+
				"Exemple:\n"+
				"  package %s",
			baseName, file.Name.Name, expectedPkg, expectedPkg)
	}

	return nil, nil
}
