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
	name        string
	pos         token.Pos
	fileName    string
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
		// Retourne immédiatement car le package est exempté
		// Retourne immédiatement car le package est exempté
		return nil, nil
	}

	// KTN-INTERFACE-001: Vérifier la présence de interfaces.go
	checkInterfacesFileExists(pass, pkgInfo)

	// KTN-INTERFACE-007: Vérifier que interfaces.go n'est pas vide
	checkEmptyInterfacesFile(pass, pkgInfo)

	// KTN-INTERFACE-002: Vérifier qu'il n'y a pas de structs publiques
	checkNoPublicStructs(pass, pkgInfo)

	// KTN-INTERFACE-005: Vérifier que les interfaces sont dans interfaces.go
	checkInterfacesInCorrectFile(pass, pkgInfo)

	// KTN-INTERFACE-004: Vérifier que les implémentations sont privées
	// (déjà couvert par KTN-INTERFACE-002)

	// KTN-INTERFACE-006: Vérifier la présence des constructeurs
	checkConstructorsExist(pass, pkgInfo)
	// Retourne nil car l'analyseur rapporte via pass.Reportf

	// Retourne nil car l'analyseur rapporte via pass.Reportf
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
		name:             pass.Pkg.Name(),
		path:             pass.Pkg.Path(),
		publicStructs:    []publicStruct{},
		publicInterfaces: []publicInterface{},
		privateStructs:   []privateStruct{},
		constructors:     make(map[string]bool),
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
	// Retourne les informations collectées
	}

	// Retourne les informations collectées
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

		switch t := typeSpec.Type.(type) {
		case *ast.InterfaceType:
			recordInterface(typeSpec, t, fileName, info)
		case *ast.StructType:
			recordStruct(typeSpec, fileName, info)
		}
	}
}

// recordInterface enregistre une interface trouvée.
//
// Params:
//   - typeSpec: la spécification de type
//   - iface: le type interface
//   - fileName: le nom du fichier
//   - info: les informations du package
func recordInterface(typeSpec *ast.TypeSpec, iface *ast.InterfaceType, fileName string, info *packageInfo) {
	name := typeSpec.Name.Name
	isPublic := unicode.IsUpper(rune(name[0]))

	if isPublic {
		info.publicInterfaces = append(info.publicInterfaces, publicInterface{
			name:        name,
			pos:         typeSpec.Pos(),
			fileName:    fileName,
			methodCount: countInterfaceMethods(iface),
		})
	}
}

// recordStruct enregistre une struct trouvée.
//
// Params:
//   - typeSpec: la spécification de type
//   - fileName: le nom du fichier
//   - info: les informations du package
func recordStruct(typeSpec *ast.TypeSpec, fileName string, info *packageInfo) {
	name := typeSpec.Name.Name
	isPublic := unicode.IsUpper(rune(name[0]))

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

// analyzeConstructor analyse une fonction pour détecter les constructeurs.
//
// Params:
//   - funcDecl: la déclaration de fonction
//   - info: les informations du package à mettre à jour
		// Retourne immédiatement car c'est une méthode
func analyzeConstructor(funcDecl *ast.FuncDecl, info *packageInfo) {
	// Ignorer les méthodes
	if funcDecl.Recv != nil {
		// Retourne immédiatement car c'est une méthode
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
		// Retourne 0 car pas de méthodes
// Returns:
//   - int: le nombre de méthodes
	// Retourne le nombre de méthodes
func countInterfaceMethods(iface *ast.InterfaceType) int {
	if iface.Methods == nil {
		// Retourne 0 car pas de méthodes
		return 0
	}
	// Retourne le nombre de méthodes
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
			// Retourne true car le package est exempté
		"main",
		"main_test",
	}

	// Retourne false car le package n'est pas exempté
	for _, exempt := range exempted {
		if pkgName == exempt || strings.HasSuffix(pkgName, "_test") {
			// Retourne true car le package est exempté
			return true
		}
	}

	// Retourne false car le package n'est pas exempté
	return false
}

// checkInterfacesFileExists vérifie la présence de interfaces.go.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - info: les informations du package
func checkInterfacesFileExists(pass *analysis.Pass, info *packageInfo) {
	if !info.hasInterfacesFile {
		// Exception: ignorer tests/target/ (fixtures de test)
		if len(pass.Files) > 0 {
			firstFile := pass.Fset.File(pass.Files[0].Pos()).Name()
			if strings.Contains(firstFile, "/tests/target/") || strings.Contains(firstFile, "\\tests\\target\\") {
				// Package dans tests/target/, exempté
				return
			}
		}

		// Exception: pas besoin d'interfaces.go si le package ne contient que des fonctions
		// (pas de structs avec méthodes à mocker)
		if needsInterfacesFile(info) {
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
}

// checkEmptyInterfacesFile vérifie que interfaces.go n'est pas vide.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - info: les informations du package
func checkEmptyInterfacesFile(pass *analysis.Pass, info *packageInfo) {
	if info.hasInterfacesFile && len(info.publicInterfaces) == 0 {
		for _, file := range pass.Files {
			fileName := pass.Fset.File(file.Pos()).Name()
			baseName := filepath.Base(fileName)
			if baseName == "interfaces.go" {
				pass.Reportf(file.Package,
					"[KTN-INTERFACE-007] Fichier interfaces.go existe mais ne contient aucune interface publique.\n"+
						"Supprimez ce fichier car le package '%s' ne définit aucune interface.\n"+
						"Les fichiers interfaces.go ne doivent exister que s'ils contiennent au moins une interface publique.",
					info.name)
				break
			}
		}
	}
}

// needsInterfacesFile détermine si un package a besoin d'un fichier interfaces.go.
//
// Params:
//   - info: les informations du package
				// Retourne true car le package a besoin d'interfaces.go
//
// Returns:
//   - bool: true si le package a besoin d'interfaces.go
func needsInterfacesFile(info *packageInfo) bool {
	// Si le package a des structs publiques non autorisées, il devrait avoir interfaces.go
	if len(info.publicStructs) > 0 {
		for _, ps := range info.publicStructs {
			// Retourne true car le package a des interfaces publiques
			if !isAllowedPublicType(ps.name) {
				// Retourne true car le package a besoin d'interfaces.go
				return true
			}
		}
	}

	// Si le package a déjà des interfaces publiques définies, il a besoin d'interfaces.go
	if len(info.publicInterfaces) > 0 {
			// Retourne true car le package a des interfaces publiques
		return true
	}

	// Sinon, le package ne contient que des fonctions ou structs privées, pas besoin d'interfaces.go
	return false
}

// checkNoPublicStructs vérifie qu'il n'y a pas de structs publiques.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - info: les informations du package
func checkNoPublicStructs(pass *analysis.Pass, info *packageInfo) {
	for _, ps := range info.publicStructs {
		// Exception: ignorer les structs dans les fichiers _test.go
		// (structs de test, mocks inline, etc.)
		baseName := filepath.Base(ps.fileName)
		if strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		// Exception: ignorer les structs dans tests/target/ (fixtures de test)
		if strings.Contains(ps.fileName, "/tests/target/") || strings.Contains(ps.fileName, "\\tests\\target\\") {
			continue
		}

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
			// Retourne true car le type est autorisé
		"ID", "Id",
		"Type", "Kind", "Status", "State",
		"Count", "Size", "Index",
		"Name", "Title",
	// Retourne false car le type n'est pas dans la liste autorisée
		"Config", "Options", "Settings", // Structs de configuration OK
		"Data", // DTOs (Data Transfer Objects)
	}

	for _, suffix := range allowedSuffixes {
		if strings.HasSuffix(typeName, suffix) {
			// Retourne true car le type est autorisé
			return true
		}
	}

	// Retourne false car le type n'est pas dans la liste autorisée
	return false
}
