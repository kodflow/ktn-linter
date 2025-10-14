package analyzer

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// InterfaceAnalyzer vérifie que les packages respectent les règles d'interfaces mockables
	InterfaceAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktninterface",
		Doc:  "Vérifie que chaque package a interfaces.go et que tous les types publics sont des interfaces",
		Run:  runInterfaceAnalyzer,
	}
)

// packageInfo contient les informations collectées sur un package.
type packageInfo struct {
	name              string
	path              string
	hasInterfacesFile bool
	publicStructs     []publicStruct
	publicInterfaces  []publicInterface
	privateStructs    []privateStruct
	constructors      map[string]bool
}

// publicStruct représente une struct publique trouvée.
type publicStruct struct {
	name     string
	pos      token.Pos
	fileName string
}

// publicInterface représente une interface publique trouvée.
type publicInterface struct {
	name       string
	pos        token.Pos
	fileName   string
	methodCount int
}

// privateStruct représente une struct privée (implémentation potentielle).
type privateStruct struct {
	name     string
	pos      token.Pos
	fileName string
}

// runInterfaceAnalyzer vérifie que toutes les interfaces respectent les règles KTN.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runInterfaceAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// Collecter les informations du package
	pkgInfo := collectPackageInfo(pass)

	// Packages exemptés
	if isExemptedPackage(pkgInfo.name) {
		return nil, nil
	}

	// KTN-INTERFACE-001: Vérifier la présence de interfaces.go
	checkInterfacesFileExists(pass, pkgInfo)

	// KTN-INTERFACE-002: Vérifier qu'il n'y a pas de structs publiques
	checkNoPublicStructs(pass, pkgInfo)

	// KTN-INTERFACE-005: Vérifier que les interfaces sont dans interfaces.go
	checkInterfacesInCorrectFile(pass, pkgInfo)

	// KTN-INTERFACE-004: Vérifier que les implémentations sont privées
	// (déjà couvert par KTN-INTERFACE-002)

	// KTN-INTERFACE-006: Vérifier la présence des constructeurs
	checkConstructorsExist(pass, pkgInfo)

	return nil, nil
}

// collectPackageInfo collecte les informations sur le package.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - *packageInfo: les informations collectées
func collectPackageInfo(pass *analysis.Pass) *packageInfo {
	info := &packageInfo{
		name:          pass.Pkg.Name(),
		path:          pass.Pkg.Path(),
		publicStructs: []publicStruct{},
		publicInterfaces: []publicInterface{},
		privateStructs: []privateStruct{},
		constructors:  make(map[string]bool),
	}

	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		baseName := filepath.Base(fileName)

		// Vérifier si c'est interfaces.go
		if baseName == "interfaces.go" {
			info.hasInterfacesFile = true
		}

		// Parcourir les déclarations
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok == token.TYPE {
					analyzeTypeDecl(d, fileName, info)
				}
			case *ast.FuncDecl:
				analyzeConstructor(d, info)
			}
		}
	}

	return info
}

// analyzeTypeDecl analyse une déclaration de type.
//
// Params:
//   - genDecl: la déclaration générale de type
//   - fileName: le nom du fichier contenant la déclaration
//   - info: les informations du package à mettre à jour
func analyzeTypeDecl(genDecl *ast.GenDecl, fileName string, info *packageInfo) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		name := typeSpec.Name.Name
		isPublic := unicode.IsUpper(rune(name[0]))

		switch typeSpec.Type.(type) {
		case *ast.InterfaceType:
			if isPublic {
				methodCount := countInterfaceMethods(typeSpec.Type.(*ast.InterfaceType))
				info.publicInterfaces = append(info.publicInterfaces, publicInterface{
					name:        name,
					pos:         typeSpec.Pos(),
					fileName:    fileName,
					methodCount: methodCount,
				})
			}
		case *ast.StructType:
			if isPublic {
				info.publicStructs = append(info.publicStructs, publicStruct{
					name:     name,
					pos:      typeSpec.Pos(),
					fileName: fileName,
				})
			} else {
				info.privateStructs = append(info.privateStructs, privateStruct{
					name:     name,
					pos:      typeSpec.Pos(),
					fileName: fileName,
				})
			}
		}
	}
}

// analyzeConstructor analyse une fonction pour détecter les constructeurs.
//
// Params:
//   - funcDecl: la déclaration de fonction
//   - info: les informations du package à mettre à jour
func analyzeConstructor(funcDecl *ast.FuncDecl, info *packageInfo) {
	// Ignorer les méthodes
	if funcDecl.Recv != nil {
		return
	}

	name := funcDecl.Name.Name
	// Constructeurs commencent par New
	if strings.HasPrefix(name, "New") {
		// Extraire le nom du type: NewMyService -> MyService
		typeName := strings.TrimPrefix(name, "New")
		info.constructors[typeName] = true
	}
}

// countInterfaceMethods compte le nombre de méthodes dans une interface.
//
// Params:
//   - iface: le type interface
//
// Returns:
//   - int: le nombre de méthodes
func countInterfaceMethods(iface *ast.InterfaceType) int {
	if iface.Methods == nil {
		return 0
	}
	return iface.Methods.NumFields()
}

// isExemptedPackage vérifie si le package est exempté des règles.
//
// Params:
//   - pkgName: le nom du package
//
// Returns:
//   - bool: true si le package est exempté
func isExemptedPackage(pkgName string) bool {
	exempted := []string{
		"main",
		"main_test",
	}

	for _, exempt := range exempted {
		if pkgName == exempt || strings.HasSuffix(pkgName, "_test") {
			return true
		}
	}

	return false
}

// checkInterfacesFileExists vérifie la présence de interfaces.go.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - info: les informations du package
func checkInterfacesFileExists(pass *analysis.Pass, info *packageInfo) {
	if !info.hasInterfacesFile {
		// Trouver la première position dans le package pour le diagnostic
		if len(pass.Files) > 0 {
			pass.Reportf(pass.Files[0].Package,
				"[KTN-INTERFACE-001] Package '%s' sans fichier interfaces.go.\n"+
					"Créez interfaces.go pour définir les interfaces publiques du package.\n"+
					"Exemple:\n"+
					"  // interfaces.go\n"+
					"  package %s\n"+
					"\n"+
					"  // MyService définit le contrat du service.\n"+
					"  type MyService interface {\n"+
					"      DoSomething(input string) (string, error)\n"+
					"  }",
				info.name, info.name)
		}
	}
}

// checkNoPublicStructs vérifie qu'il n'y a pas de structs publiques.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - info: les informations du package
func checkNoPublicStructs(pass *analysis.Pass, info *packageInfo) {
	for _, ps := range info.publicStructs {
		// Vérifier si c'est un type natif autorisé
		if isAllowedPublicType(ps.name) {
			continue
		}

		pass.Reportf(ps.pos,
			"[KTN-INTERFACE-002] Type public '%s' défini comme struct au lieu d'interface.\n"+
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
			ps.name, ps.name, strings.ToLower(ps.name[:1])+ps.name[1:]+"Impl")
	}
}

// checkInterfacesInCorrectFile vérifie que les interfaces sont dans interfaces.go.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - info: les informations du package
func checkInterfacesInCorrectFile(pass *analysis.Pass, info *packageInfo) {
	for _, pi := range info.publicInterfaces {
		baseName := filepath.Base(pi.fileName)
		if baseName != "interfaces.go" {
			pass.Reportf(pi.pos,
				"[KTN-INTERFACE-005] Interface '%s' définie dans %s.\n"+
					"Les interfaces publiques doivent être dans interfaces.go.\n"+
					"Déplacez cette interface vers interfaces.go.",
				pi.name, baseName)
		}
	}
}

// checkConstructorsExist vérifie la présence des constructeurs.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les avertissements
//   - info: les informations du package
func checkConstructorsExist(pass *analysis.Pass, info *packageInfo) {
	for _, pi := range info.publicInterfaces {
		if !info.constructors[pi.name] {
			// Avertissement seulement si l'interface a des méthodes
			if pi.methodCount > 0 {
				pass.Reportf(pi.pos,
					"[KTN-INTERFACE-006] Interface '%s' sans constructeur.\n"+
						"Ajoutez un constructeur qui retourne l'interface.\n"+
						"Exemple:\n"+
						"  func New%s(deps ...) %s {\n"+
						"      return &%s{...}\n"+
						"  }",
					pi.name, pi.name, pi.name, strings.ToLower(pi.name[:1])+pi.name[1:]+"Impl")
			}
		}
	}
}

// isAllowedPublicType vérifie si un type public est autorisé sans être une interface.
//
// Params:
//   - typeName: le nom du type
//
// Returns:
//   - bool: true si le type est autorisé
func isAllowedPublicType(typeName string) bool {
	// Types qui peuvent être publics sans être des interfaces:
	// - Types qui ressemblent à des enums (suffixe comme Status, Type, Kind)
	// - Types qui sont clairement des aliases (ID, Count, Name, etc.)

	allowedSuffixes := []string{
		"ID", "Id",
		"Type", "Kind", "Status", "State",
		"Count", "Size", "Index",
		"Name", "Title",
		"Config", "Options", "Settings", // Structs de configuration OK
		"Data", // DTOs (Data Transfer Objects)
	}

	for _, suffix := range allowedSuffixes {
		if strings.HasSuffix(typeName, suffix) {
			return true
		}
	}

	return false
}
