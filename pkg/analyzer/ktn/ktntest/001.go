// Analyzer 001 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
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
	// Check each file in the package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		base := filepath.Base(filename)

		// Skip if not a test file
		if !strings.HasSuffix(base, "_test.go") {
			continue
		}

		// Check valid suffixes
		validSuffix := hasValidTestSuffix(base)

		// If valid suffix, verify content matches suffix
		if validSuffix {
			// Verify _bench_test.go contains only Benchmark* functions
			if strings.HasSuffix(base, "_bench_test.go") {
				verifyBenchFile(pass, file, base)
			}
			// _integration_test.go is accepted without content verification
			continue
		}

		// Report error if it's a plain _test.go file
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
	// Track non-benchmark test functions
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		// Skip non-functions
		if !ok {
			continue
		}

		name := funcDecl.Name.Name
		// Skip non-test functions (helpers, etc.)
		if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Benchmark") {
			continue
		}

		// Report if Test* function found in bench file
		if strings.HasPrefix(name, "Test") {
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-001: le fichier '%s' ne doit contenir que des fonctions Benchmark*, pas '%s'. Déplacer vers _internal_test.go ou _external_test.go",
				base,
				name,
			)
		}
	}
}
