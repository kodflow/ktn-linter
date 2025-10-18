package ktn_interface

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Rule003 analyzer for KTN linter.
var Rule003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_INTERFACE_003",
	Doc:  "Vérifie que les interfaces publiques sont dans interfaces.go",
	Run:  RunRule003,
}

// RunRule003 exécute la règle 003.
func RunRule003(pass *analysis.Pass) (any, error) {
	// Packages exemptés
	if IsExemptedPackage003(pass.Pkg.Name()) {
		// Analysis completed successfully.
		return nil, nil
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Ignorer interfaces.go lui-même
		if baseName == "interfaces.go" {
			continue
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				name := typeSpec.Name.Name
				isPublic := unicode.IsUpper(rune(name[0]))

				if !isPublic {
					continue
				}

				_, isInterface := typeSpec.Type.(*ast.InterfaceType)
				if !isInterface {
					continue
				}

				pass.Reportf(typeSpec.Pos(),
					"[KTN_INTERFACE_003] Interface '%s' définie dans %s.\n"+
						"Les interfaces publiques doivent être dans interfaces.go.\n"+
						"Déplacez cette interface vers interfaces.go.",
					name, baseName)
			}
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

// IsExemptedPackage003 vérifie si un package est exempté pour Rule003.
func IsExemptedPackage003(pkgName string) bool {
	exempted := []string{"main", "main_test"}
	for _, exempt := range exempted {
		if pkgName == exempt || strings.HasSuffix(pkgName, "_test") {
			// Continue traversing AST nodes.
			return true
		}
	}
	// Condition not met, return false.
	return false
}
