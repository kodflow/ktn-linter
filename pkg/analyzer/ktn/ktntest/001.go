// Package ktntest provides analyzers for test file lint rules.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest001 string = "KTN-TEST-001"
)

// Analyzer001 detects test files not following naming convention.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest001",
	Doc:  "KTN-TEST-001: fichier de test doit se terminer par _internal_test.go, _external_test.go, _bench_test.go ou _integration_test.go",
	Run:  runTest001,
}

// runTest001 analyzes test file naming conventions.
//
// Params:
//   - pass: Analysis pass
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runTest001(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest001) {
		// Règle désactivée
		return nil, nil
	}

	// Parcourir tous les fichiers du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest001, filename) {
			// Fichier exclu
			continue
		}

		base := filepath.Base(filename)

		// Vérification si fichier de test
		if !strings.HasSuffix(base, "_test.go") {
			// Pas un fichier de test, continuer
			continue
		}

		// Vérifier les suffixes valides
		validSuffix := hasValidTestSuffix(base)

		// Vérification du suffixe valide
		if validSuffix {
			// Vérification si fichier bench
			if strings.HasSuffix(base, "_bench_test.go") {
				verifyBenchFile(pass, file, base)
			}
			// Fichier valide, continuer
			continue
		}

		// Signalement de l'erreur
		pass.Reportf(
			file.Pos(),
			"KTN-TEST-001: le fichier '%s' doit être renommé en '%s' (white-box), '%s' (black-box), '%s' (benchmarks) ou '%s' (tests d'intégration)",
			base,
			strings.Replace(base, "_test.go", "_internal_test.go", 1),
			strings.Replace(base, "_test.go", "_external_test.go", 1),
			strings.Replace(base, "_test.go", "_bench_test.go", 1),
			strings.Replace(base, "_test.go", "_integration_test.go", 1),
		)
	}

	// Retour de la fonction
	return nil, nil
}

// hasValidTestSuffix checks if a test file has a valid suffix.
//
// Params:
//   - base: base filename to check
//
// Returns:
//   - bool: true if suffix is valid
func hasValidTestSuffix(base string) bool {
	// Check all valid suffixes
	return strings.HasSuffix(base, "_internal_test.go") ||
		strings.HasSuffix(base, "_external_test.go") ||
		strings.HasSuffix(base, "_bench_test.go") ||
		strings.HasSuffix(base, "_integration_test.go")
}

// verifyBenchFile verifies that a _bench_test.go file contains only Benchmark functions.
//
// Params:
//   - pass: Analysis pass
//   - file: AST file to check
//   - base: base filename
func verifyBenchFile(pass *analysis.Pass, file *ast.File, base string) {
	// Parcourir les déclarations du fichier
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		// Vérification si fonction
		if !ok {
			// Pas une fonction, continuer
			continue
		}

		name := funcDecl.Name.Name
		// Vérification si fonction de test
		if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Benchmark") {
			// Pas une fonction de test, continuer
			continue
		}

		// Vérification si fonction Test dans bench
		if strings.HasPrefix(name, "Test") {
			// Signalement de l'erreur
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-001: le fichier '%s' ne doit contenir que des fonctions Benchmark*, pas '%s'. Déplacer vers _internal_test.go ou _external_test.go",
				base,
				name,
			)
		}
	}
}
