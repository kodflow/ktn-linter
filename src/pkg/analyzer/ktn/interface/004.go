package ktn_interface

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Rule004 = &analysis.Analyzer{
	Name: "KTN_INTERFACE_004",
	Doc:  "Vérifie la présence des constructeurs pour les interfaces",
	Run:  RunRule004,
}

// RunRule004 exécute la règle 004.
func RunRule004(pass *analysis.Pass) (any, error) {
	// Packages exemptés
	if IsExemptedPackage004(pass.Pkg.Name()) {
		return nil, nil
	}

	// Collecter les interfaces publiques
	publicInterfaces := make(map[string]*interfaceInfo)
	constructors := make(map[string]bool)

	for _, file := range pass.Files {
		// Collecter les interfaces
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

				ifaceType, isInterface := typeSpec.Type.(*ast.InterfaceType)
				if !isInterface {
					continue
				}

				methodCount := 0
				if ifaceType.Methods != nil {
					methodCount = ifaceType.Methods.NumFields()
				}

				publicInterfaces[name] = &interfaceInfo{
					name:        name,
					pos:         typeSpec.Pos(),
					methodCount: methodCount,
				}
			}
		}

		// Collecter les constructeurs
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Recv != nil {
				continue
			}

			name := funcDecl.Name.Name
			if strings.HasPrefix(name, "New") {
				typeName := strings.TrimPrefix(name, "New")
				constructors[typeName] = true
			}
		}
	}

	// Vérifier que chaque interface a un constructeur
	for _, iface := range publicInterfaces {
		if !constructors[iface.name] && iface.methodCount > 0 {
			pass.Reportf(iface.pos,
				"[KTN_INTERFACE_004] Interface '%s' sans constructeur.\n"+
					"Ajoutez un constructeur qui retourne l'interface.\n"+
					"Exemple:\n"+
					"  func New%s(deps ...) %s {\n"+
					"      return &%s{...}\n"+
					"  }",
				iface.name, iface.name, iface.name, strings.ToLower(iface.name[:1])+iface.name[1:]+"Impl")
		}
	}

	return nil, nil
}

type interfaceInfo struct {
	name        string
	pos         token.Pos
	methodCount int
}

// IsExemptedPackage004 vérifie si un package est exempté pour Rule004.
func IsExemptedPackage004(pkgName string) bool {
	exempted := []string{"main", "main_test"}
	for _, exempt := range exempted {
		if pkgName == exempt || strings.HasSuffix(pkgName, "_test") {
			return true
		}
	}
	return false
}
