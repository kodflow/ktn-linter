// Analyzer 001 for the ktninterface package.
package ktninterface

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	// INITIAL_INTERFACES_CAP initial capacity for interfaces map
	INITIAL_INTERFACES_CAP int = 16
	// INTERFACE_SUFFIX_LEN length of "Interface" suffix
	INTERFACE_SUFFIX_LEN int = 9
)

// Analyzer001 detects unused interface declarations.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktninterface001",
	Doc:      "KTN-INTERFACE-001: interface non utilisée",
	Run:      runInterface001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface001 analyzes interfaces to detect unused ones.
// Params:
//   - pass: Analysis pass
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runInterface001(pass *analysis.Pass) (any, error) {
	// Collect all interface declarations
	interfaces := make(map[string]*ast.TypeSpec, INITIAL_INTERFACES_CAP)
	usedInterfaces := make(map[string]bool, INITIAL_INTERFACES_CAP)
	structNames := make(map[string]bool, INITIAL_INTERFACES_CAP)

	// First pass: collect all interface and struct declarations
	collectDeclarations(pass, interfaces, structNames)

	// Second pass: find interface usages
	findInterfaceUsages(pass, usedInterfaces)

	// Report unused interfaces
	reportUnusedInterfaces(pass, interfaces, usedInterfaces, structNames)

	// Retour de la fonction
	return nil, nil
}

// collectDeclarations collecte les déclarations d'interfaces et de structs.
//
// Params:
//   - pass: contexte d'analyse
//   - interfaces: map pour stocker les interfaces
//   - structNames: map pour stocker les noms de structs
func collectDeclarations(pass *analysis.Pass, interfaces map[string]*ast.TypeSpec, structNames map[string]bool) {
	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			genDecl, isGenDecl := node.(*ast.GenDecl)
			// Continue if not general declaration
			if !isGenDecl {
				return true
			}

			// Parcourir les specs
			collectTypeSpecs(genDecl.Specs, interfaces, structNames)
			// Continuer l'itération
			return true
		})
	}
}

// collectTypeSpecs collecte les TypeSpecs d'une déclaration.
//
// Params:
//   - specs: liste des specs à analyser
//   - interfaces: map pour stocker les interfaces
//   - structNames: map pour stocker les noms de structs
func collectTypeSpecs(specs []ast.Spec, interfaces map[string]*ast.TypeSpec, structNames map[string]bool) {
	// Itération sur les specs
	for _, spec := range specs {
		typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
		// Continue if not type spec
		if !isTypeSpec {
			continue
		}

		_, isInterface := typeSpec.Type.(*ast.InterfaceType)
		// Store if interface type
		if isInterface {
			interfaces[typeSpec.Name.Name] = typeSpec
		}

		_, isStruct := typeSpec.Type.(*ast.StructType)
		// Store struct names detection
		if isStruct {
			structNames[typeSpec.Name.Name] = true
		}
	}
}

// findInterfaceUsages trouve les usages d'interfaces dans le code.
//
// Params:
//   - pass: contexte d'analyse
//   - usedInterfaces: map pour marquer les interfaces utilisées
func findInterfaceUsages(pass *analysis.Pass, usedInterfaces map[string]bool) {
	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			checkNodeForInterfaceUsage(node, usedInterfaces)
			// Continuer l'itération
			return true
		})
	}
}

// checkNodeForInterfaceUsage vérifie un nœud AST pour les usages d'interface.
//
// Params:
//   - node: nœud AST à vérifier
//   - usedInterfaces: map pour marquer les interfaces utilisées
func checkNodeForInterfaceUsage(node ast.Node, usedInterfaces map[string]bool) {
	// Verification de la condition
	switch n := node.(type) {
	// Cas FuncDecl
	case *ast.FuncDecl:
		checkFuncDeclForInterfaces(n, usedInterfaces)
	// Cas Field
	case *ast.Field:
		checkType(n.Type, usedInterfaces)
	// Cas InterfaceType
	case *ast.InterfaceType:
		checkEmbeddedInterfaces(n, usedInterfaces)
	}
}

// checkFuncDeclForInterfaces vérifie les paramètres et retours d'une fonction.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - usedInterfaces: map pour marquer les interfaces utilisées
func checkFuncDeclForInterfaces(funcDecl *ast.FuncDecl, usedInterfaces map[string]bool) {
	// Check parameters
	if funcDecl.Type.Params != nil {
		checkFieldList(funcDecl.Type.Params, usedInterfaces)
	}
	// Check results
	if funcDecl.Type.Results != nil {
		checkFieldList(funcDecl.Type.Results, usedInterfaces)
	}
}

// checkEmbeddedInterfaces vérifie les interfaces embarquées.
//
// Params:
//   - interfaceType: type interface
//   - usedInterfaces: map pour marquer les interfaces utilisées
func checkEmbeddedInterfaces(interfaceType *ast.InterfaceType, usedInterfaces map[string]bool) {
	// Vérifier si les méthodes existent
	if interfaceType.Methods == nil {
		return
	}
	// Itération sur les méthodes
	for _, method := range interfaceType.Methods.List {
		// Embedded interface has no function type
		if method.Type != nil {
			checkType(method.Type, usedInterfaces)
		}
	}
}

// reportUnusedInterfaces reporte les interfaces non utilisées.
//
// Params:
//   - pass: contexte d'analyse
//   - interfaces: map des interfaces trouvées
//   - usedInterfaces: map des interfaces utilisées
//   - structNames: map des noms de structs
func reportUnusedInterfaces(pass *analysis.Pass, interfaces map[string]*ast.TypeSpec, usedInterfaces map[string]bool, structNames map[string]bool) {
	// Itération sur les interfaces
	for name, typeSpec := range interfaces {
		// Skip if interface is used
		if usedInterfaces[name] {
			continue
		}

		// Skip if interface follows XXXInterface pattern for struct XXX
		// Ces interfaces sont légitimes pour le mocking de la struct
		if isStructInterfacePattern(name, structNames) {
			continue
		}

		// Skip if a struct with same name exists (interface for struct)
		// L'interface est légitime car elle permet de mocker la struct
		if hasCorrespondingStruct(name, structNames) {
			continue
		}

		// Report unused interface with helpful message for developers/AI
		pass.Reportf(
			typeSpec.Pos(),
			"KTN-INTERFACE-001: interface '%s' non utilisée. "+
				"Options: (1) créer une struct '%s' qui l'implémente pour permettre le mocking, "+
				"(2) utiliser cette interface en paramètre/retour de fonction, "+
				"(3) supprimer si vraiment inutile. "+
				"Les interfaces permettent de créer des mocks et d'améliorer la couverture de tests",
			name,
			name,
		)
	}
}

// hasCorrespondingStruct vérifie si une struct correspondante existe.
// Par exemple, pour une interface "UserService", vérifie si "UserService" struct existe.
//
// Params:
//   - interfaceName: nom de l'interface
//   - structs: map des noms de structs
//
// Returns:
//   - bool: true si une struct correspondante existe
func hasCorrespondingStruct(interfaceName string, structs map[string]bool) bool {
	// Vérifier si une struct avec le même nom existe
	return structs[interfaceName]
}

// isStructInterfacePattern checks if interface follows XXXInterface pattern.
//
// Params:
//   - interfaceName: Interface name to check
//   - structs: Map of struct names
//
// Returns:
//   - bool: true if interface follows the pattern
func isStructInterfacePattern(interfaceName string, structs map[string]bool) bool {
	// Check if interface name ends with "Interface"
	if len(interfaceName) <= INTERFACE_SUFFIX_LEN || interfaceName[len(interfaceName)-INTERFACE_SUFFIX_LEN:] != "Interface" {
		return false
	}

	// Extract potential struct name (remove "Interface" suffix)
	structName := interfaceName[:len(interfaceName)-INTERFACE_SUFFIX_LEN]

	// Check if corresponding struct exists
	return structs[structName]
}

// checkFieldList checks field list for interface usage.
// Params:
//   - fields: Field list to check
//   - used: Map to track used interfaces
func checkFieldList(fields *ast.FieldList, used map[string]bool) {
	// Verification de la condition
	for _, field := range fields.List {
		checkType(field.Type, used)
	}
}

// checkType checks type for interface usage.
// Params:
//   - expr: Expression to check
//   - used: Map to track used interfaces
func checkType(expr ast.Expr, used map[string]bool) {
	// Verification de la condition
	switch t := expr.(type) {
	// Verification de la condition
	case *ast.Ident:
		// Mark identifier as used
		used[t.Name] = true
	// Verification de la condition
	case *ast.StarExpr:
		// Check pointer type
		checkType(t.X, used)
	// Verification de la condition
	case *ast.ArrayType:
		// Check array element type
		checkType(t.Elt, used)
	// Verification de la condition
	case *ast.MapType:
		// Check map key and value types
		checkType(t.Key, used)
		checkType(t.Value, used)
	// Verification de la condition
	case *ast.ChanType:
		// Check channel element type
		checkType(t.Value, used)
	// Verification de la condition
	case *ast.SelectorExpr:
		// Check selector expression
		checkType(t.X, used)
	}
}
