package ktn_interface

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Rule001 = &analysis.Analyzer{
	Name: "KTN_INTERFACE_001",
	Doc:  "Vérifie la présence du fichier interfaces.go",
	Run:  runRule001,
}

func runRule001(pass *analysis.Pass) (any, error) {
	// Collecter les informations du package
	pkgInfo := collectPackageInfo(pass)

	// Packages exemptés
	if isExemptedPackage(pkgInfo.name) {
		return nil, nil
	}

	// Vérifier la présence de interfaces.go
	if !pkgInfo.hasInterfacesFile {
		// Exception: ignorer tests/target/ (fixtures de test)
		if len(pass.Files) > 0 {
			firstFile := pass.Fset.File(pass.Files[0].Pos()).Name()
			if strings.Contains(firstFile, "/tests/target/") || strings.Contains(firstFile, "\\tests\\target\\") {
				return nil, nil
			}
		}

		// Exception: pas besoin d'interfaces.go si le package ne contient que des fonctions
		if needsInterfacesFile(pkgInfo) {
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
					pkgInfo.name, pkgInfo.name)
			}
		}
	}

	return nil, nil
}

type packageInfo struct {
	name              string
	hasInterfacesFile bool
	publicStructs     []publicStruct
	publicInterfaces  []publicInterface
}

type publicStruct struct {
	name     string
	fileName string
}

type publicInterface struct {
	name string
}

func collectPackageInfo(pass *analysis.Pass) *packageInfo {
	info := &packageInfo{
		name:             pass.Pkg.Name(),
		publicStructs:    []publicStruct{},
		publicInterfaces: []publicInterface{},
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		if baseName == "interfaces.go" {
			info.hasInterfacesFile = true
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
					info.publicInterfaces = append(info.publicInterfaces, publicInterface{name: name})
				case *ast.StructType:
					info.publicStructs = append(info.publicStructs, publicStruct{
						name:     name,
						fileName: fileName,
					})
				}
			}
		}
	}

	return info
}

func isExemptedPackage(pkgName string) bool {
	exempted := []string{"main", "main_test"}
	for _, exempt := range exempted {
		if pkgName == exempt || strings.HasSuffix(pkgName, "_test") {
			return true
		}
	}
	return false
}

func needsInterfacesFile(info *packageInfo) bool {
	if len(info.publicStructs) > 0 {
		for _, ps := range info.publicStructs {
			if !isAllowedPublicType(ps.name) {
				return true
			}
		}
	}

	if len(info.publicInterfaces) > 0 {
		return true
	}

	return false
}

func isAllowedPublicType(typeName string) bool {
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
			return true
		}
	}

	return false
}
