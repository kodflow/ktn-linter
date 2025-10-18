package ktn_test

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule002 analyzer for KTN linter.
var Rule002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_TEST_002",
	Doc:  "Vérifie que chaque fichier .go a son _test.go",
	Run:  runRule002,
}

func runRule002(pass *analysis.Pass) (any, error) {
	// Package main exempté
	if pass.Pkg.Name() == "main" {
		// Analysis completed successfully.
		return nil, nil
	}

	// Package de test exempté
	if strings.HasSuffix(pass.Pkg.Name(), "_test") {
		// Analysis completed successfully.
		return nil, nil
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Ignorer les fichiers _test.go
		if strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		// Ignorer les fixtures de test
		if strings.Contains(fileName, "/tests/target/") || strings.Contains(fileName, "\\tests\\target\\") ||
			strings.Contains(fileName, "/tests/bad_usage/") || strings.Contains(fileName, "\\tests\\bad_usage\\") ||
			strings.Contains(fileName, "/tests/good_usage/") || strings.Contains(fileName, "\\tests\\good_usage\\") {
			continue
		}

		// Ignorer mock.go
		if baseName == "mock.go" {
			continue
		}

		// Ignorer interfaces.go SI ET SEULEMENT SI il contient uniquement des interfaces
		if baseName == "interfaces.go" && containsOnlyInterfaces002(file) {
			continue
		}

		// Ignorer les fichiers qui ne contiennent que des const/var
		if !hasTestableElements002(file) {
			continue
		}

		// Vérifier la présence du fichier de test
		testFileName := strings.TrimSuffix(baseName, ".go") + "_test.go"
		testPath := filepath.Join(filepath.Dir(fileName), testFileName)

		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			pass.Reportf(file.Package,
				"[KTN_TEST_002] Fichier '%s' n'a pas de fichier de test correspondant.\n"+
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
				baseName, testFileName, testFileName, pass.Pkg.Name())
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

func containsOnlyInterfaces002(file *ast.File) bool {
	hasInterface := false

	for _, decl := range file.Decls {
		_, isFunc := decl.(*ast.FuncDecl)
		if isFunc {
			// Condition not met, return false.
			return false
		}

		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			_, isInterface := typeSpec.Type.(*ast.InterfaceType)
			if isInterface {
				hasInterface = true
			} else {
				// Condition not met, return false.
				return false
			}
		}
	}

	// Early return from function.
	return hasInterface
}

func hasTestableElements002(file *ast.File) bool {
	for _, decl := range file.Decls {
		if isTestableFunction002(decl) || isTestableType002(decl) {
			// Continue traversing AST nodes.
			return true
		}
	}
	// Condition not met, return false.
	return false
}

func isTestableFunction002(decl ast.Decl) bool {
	funcDecl, ok := decl.(*ast.FuncDecl)
	if !ok {
		// Condition not met, return false.
		return false
	}

	name := funcDecl.Name.Name
	if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
		// Condition not met, return false.
		return false
	}
	// Continue traversing AST nodes.
	return true
}

func isTestableType002(decl ast.Decl) bool {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok || genDecl.Tok.String() != "type" {
		// Condition not met, return false.
		return false
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch typeSpec.Type.(type) {
		case *ast.StructType, *ast.InterfaceType:
			// Continue traversing AST nodes.
			return true
		}
	}
	// Condition not met, return false.
	return false
}
