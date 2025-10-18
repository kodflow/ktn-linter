package ktn_interface

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Rule001 analyzer for KTN linter.
var Rule001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_INTERFACE_001",
	Doc:  "Vérifie la présence du fichier interfaces.go",
	Run:  RunRule001,
}

// RunRule001 exécute la règle 001.
func RunRule001(pass *analysis.Pass) (any, error) {
	// Collecter les informations du package
	pkgInfo := CollectPackageInfo(pass)

	// Packages exemptés
	if IsExemptedPackage(pkgInfo.Name) {
		// Analysis completed successfully.
		return nil, nil
	}

	// Vérifier la présence de interfaces.go
	if !pkgInfo.HasInterfacesFile {
		// Exception: ignorer tests/target/ (fixtures de test)
		if len(pass.Files) > 0 {
			firstFile := pass.Fset.File(pass.Files[0].Pos()).Name()
			if strings.Contains(firstFile, "/tests/target/") || strings.Contains(firstFile, "\\tests\\target\\") {
				// Analysis completed successfully.
				return nil, nil
			}
		}

		// Exception: pas besoin d'interfaces.go si le package ne contient que des fonctions
		if NeedsInterfacesFile(pkgInfo) {
			if len(pass.Files) > 0 {
				pass.Reportf(pass.Files[0].Package,
					"[KTN_INTERFACE_001] Package '%s' sans fichier interfaces.go.\n"+
						"Créez interfaces.go pour définir les interfaces publiques du package.\n"+
						"Exemple:\n"+
						"  // interfaces.go\n"+
						"  package %s\n"+
						"\n"+
						"  // MyService définit le contrat du service.\n"+
						"  type MyService interface {\n"+
						"      DoSomething(input string) (string, error)\n"+
						"  }",
					pkgInfo.Name, pkgInfo.Name)
			}
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

// PackageInfo contient les informations collectées sur un package.
type PackageInfo struct {
	Name              string
	HasInterfacesFile bool
	PublicStructs     []PublicStruct
	PublicInterfaces  []PublicInterface
}

// PublicStruct représente une structure publique dans un package.
type PublicStruct struct {
	Name     string
	FileName string
}

// PublicInterface représente une interface publique dans un package.
type PublicInterface struct {
	Name string
}

// CollectPackageInfo collecte les informations sur un package.
func CollectPackageInfo(pass *analysis.Pass) *PackageInfo {
	info := &PackageInfo{
		Name:             pass.Pkg.Name(),
		PublicStructs:    []PublicStruct{},
		PublicInterfaces: []PublicInterface{},
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		if baseName == "interfaces.go" {
			info.HasInterfacesFile = true
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

				switch typeSpec.Type.(type) {
				case *ast.InterfaceType:
					info.PublicInterfaces = append(info.PublicInterfaces, PublicInterface{Name: name})
				case *ast.StructType:
					info.PublicStructs = append(info.PublicStructs, PublicStruct{
						Name:     name,
						FileName: fileName,
					})
				}
			}
		}
	}

	// Early return from function.
	return info
}

// IsExemptedPackage vérifie si un package est exempté.
func IsExemptedPackage(pkgName string) bool {
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

// NeedsInterfacesFile détermine si un package a besoin d'un fichier interfaces.go.
func NeedsInterfacesFile(info *PackageInfo) bool {
	if len(info.PublicStructs) > 0 {
		for _, ps := range info.PublicStructs {
			if !IsAllowedPublicType(ps.Name) {
				// Continue traversing AST nodes.
				return true
			}
		}
	}

	if len(info.PublicInterfaces) > 0 {
		// Continue traversing AST nodes.
		return true
	}

	// Condition not met, return false.
	return false
}

// IsAllowedPublicType vérifie si un type public est autorisé.
func IsAllowedPublicType(typeName string) bool {
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
