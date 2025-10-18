package ktn_interface

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Rule004 analyzer for KTN linter.
var Rule004 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_INTERFACE_004",
	Doc:  "Vérifie la présence des constructeurs pour les interfaces",
	Run:  RunRule004,
}

// RunRule004 exécute la règle 004.
func RunRule004(pass *analysis.Pass) (any, error) {
	if IsExemptedPackage004(pass.Pkg.Name()) {
		return nil, nil
	}

	publicInterfaces := make(map[string]*interfaceInfo)
	constructors := make(map[string]bool)

	collectInterfacesAndConstructors(pass, publicInterfaces, constructors)
	checkConstructorsExist(pass, publicInterfaces, constructors)

	return nil, nil
}

// collectInterfacesAndConstructors collecte les interfaces publiques et les constructeurs.
func collectInterfacesAndConstructors(pass *analysis.Pass, interfaces map[string]*interfaceInfo, constructors map[string]bool) {
	for _, file := range pass.Files {
		collectPublicInterfaces(file, interfaces)
		collectConstructors(file, constructors)
	}
}

// collectPublicInterfaces collecte toutes les interfaces publiques.
func collectPublicInterfaces(file *ast.File, interfaces map[string]*interfaceInfo) {
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
			if !unicode.IsUpper(rune(name[0])) {
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

			interfaces[name] = &interfaceInfo{
				name:        name,
				pos:         typeSpec.Pos(),
				methodCount: methodCount,
			}
		}
	}
}

// collectConstructors collecte tous les constructeurs (fonctions New*).
func collectConstructors(file *ast.File, constructors map[string]bool) {
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

// checkConstructorsExist vérifie que chaque interface a un constructeur.
func checkConstructorsExist(pass *analysis.Pass, interfaces map[string]*interfaceInfo, constructors map[string]bool) {
	for _, iface := range interfaces {
		if !constructors[iface.name] && iface.methodCount > 0 {
			reportMissingConstructor(pass, iface)
		}
	}
}

// reportMissingConstructor signale l'absence d'un constructeur pour une interface.
func reportMissingConstructor(pass *analysis.Pass, iface *interfaceInfo) {
	implName := strings.ToLower(iface.name[:1]) + iface.name[1:] + "Impl"
	pass.Reportf(iface.pos,
		"[KTN_INTERFACE_004] Interface '%s' sans constructeur.\n"+
			"Ajoutez un constructeur qui retourne l'interface.\n"+
			"Exemple:\n"+
			"  func New%s(deps ...) %s {\n"+
			"      return &%s{...}\n"+
			"  }",
		iface.name, iface.name, iface.name, implName)
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
