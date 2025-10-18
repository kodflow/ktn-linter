package ktn_interface

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Rule005 analyzer for KTN linter.
var Rule005 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_INTERFACE_005",
	Doc:  "Vérifie que interfaces.go n'est pas vide",
	Run:  RunRule005,
}

// RunRule005 exécute la règle 005.
func RunRule005(pass *analysis.Pass) (any, error) {
	// Packages exemptés
	if IsExemptedPackage005(pass.Pkg.Name()) {
		// Analysis completed successfully.
		return nil, nil
	}

	hasInterfacesFile := false
	publicInterfaceCount := 0
	var interfacesFilePos token.Pos

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		if baseName == "interfaces.go" {
			hasInterfacesFile = true
			interfacesFilePos = file.Package

			// Compter les interfaces publiques dans ce fichier
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
					if isInterface {
						publicInterfaceCount++
					}
				}
			}
		}
	}

	if hasInterfacesFile && publicInterfaceCount == 0 {
		pass.Reportf(interfacesFilePos,
			"[KTN_INTERFACE_005] Fichier interfaces.go existe mais ne contient aucune interface publique.\n"+
				"Supprimez ce fichier car le package '%s' ne définit aucune interface.\n"+
				"Les fichiers interfaces.go ne doivent exister que s'ils contiennent au moins une interface publique.",
			pass.Pkg.Name())
	}

	// Analysis completed successfully.
	return nil, nil
}

// IsExemptedPackage005 vérifie si un package est exempté pour Rule005.
func IsExemptedPackage005(pkgName string) bool {
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
