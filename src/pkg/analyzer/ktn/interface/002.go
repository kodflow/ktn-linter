package ktn_interface

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Rule002 analyzer for KTN linter.
var Rule002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_INTERFACE_002",
	Doc:  "Vérifie qu'il n'y a pas de structs publiques",
	Run:  RunRule002,
}

// RunRule002 exécute la règle 002.
func RunRule002(pass *analysis.Pass) (any, error) {
	// Packages exemptés
	if IsExemptedPackage002(pass.Pkg.Name()) {
		// Analysis completed successfully.
		return nil, nil
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Exception: ignorer les structs dans les fichiers _test.go
		if strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		// Exception: ignorer les structs dans tests/target/
		if strings.Contains(fileName, "/tests/target/") || strings.Contains(fileName, "\\tests\\target\\") {
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

				_, isStruct := typeSpec.Type.(*ast.StructType)
				if !isStruct {
					continue
				}

				// Vérifier si c'est un type autorisé
				if isAllowedPublicType002(name) {
					continue
				}

				pass.Reportf(typeSpec.Pos(),
					"[KTN_INTERFACE_002] Type public '%s' défini comme struct au lieu d'interface.\n"+
						"Les types publics doivent être des interfaces dans interfaces.go.\n"+
						"Déplacez la struct vers une implémentation privée et créez l'interface.\n"+
						"Exemple:\n"+
						"  // interfaces.go\n"+
						"  type %s interface {\n"+
						"      // Méthodes publiques ici\n"+
						"  }\n"+
						"\n"+
						"  // impl.go\n"+
						"  type %s struct { ... }",
					name, name, strings.ToLower(name[:1])+name[1:]+"Impl")
			}
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

// IsExemptedPackage002 vérifie si un package est exempté pour Rule002.
func IsExemptedPackage002(pkgName string) bool {
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

func isAllowedPublicType002(typeName string) bool {
	allowedSuffixes := []string{
		"ID", "Id",
		"Type", "Kind", "Status", "State",
		"Count", "Size", "Index",
		"Name", "Title",
		"Config", "Options", "Settings",
		"Data",
	}

	for _, suffix := range allowedSuffixes {
		if strings.HasSuffix(typeName, suffix) {
			// Continue traversing AST nodes.
			return true
		}
	}

	// Condition not met, return false.
	return false
}
