// Analyzer 001 for the ktninterface package.
package ktninterface

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	// initialInterfacesCap initial capacity for interfaces map
	initialInterfacesCap int = 16
	// interfaceSuffixLen length of "Interface" suffix
	interfaceSuffixLen int = 9
	// ruleCodeInterface001 rule code for INTERFACE-001
	ruleCodeInterface001 string = "KTN-INTERFACE-001"
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeInterface001) {
		// Règle désactivée
		return nil, nil
	}

	// Collect all interface declarations
	interfaces := make(map[string]*ast.TypeSpec, initialInterfacesCap)
	usedInterfaces := make(map[string]bool, initialInterfacesCap)
	structNames := make(map[string]bool, initialInterfacesCap)

	// First pass: collect all interface and struct declarations
	collectDeclarations(pass, cfg, interfaces, structNames)

	// Second pass: find interface usages
	findInterfaceUsages(pass, usedInterfaces)

	// Third pass: collect compile-time checks (var _ I = (*S)(nil))
	compileTimeChecks := collectCompileTimeChecks(pass)

	// Report unused interfaces
	reportUnusedInterfaces(pass, interfaces, usedInterfaces, structNames, compileTimeChecks)

	// Retour de la fonction
	return nil, nil
}

// collectDeclarations collecte les déclarations d'interfaces et de structs.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration du linter
//   - interfaces: map pour stocker les interfaces
//   - structNames: map pour stocker les noms de structs
func collectDeclarations(pass *analysis.Pass, cfg *config.Config, interfaces map[string]*ast.TypeSpec, structNames map[string]bool) {
	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			genDecl, isGenDecl := node.(*ast.GenDecl)
			// Continue if not general declaration
			if !isGenDecl {
				// Continuer l'itération
				return true
			}

			// Vérifier si le fichier est exclu
			filename := pass.Fset.Position(genDecl.Pos()).Filename
			// Vérification si le fichier courant est exclu par la configuration de la règle
			if cfg.IsFileExcluded(ruleCodeInterface001, filename) {
				// Fichier exclu
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

// collectCompileTimeChecks collecte les compile-time interface checks.
// Pattern: var _ InterfaceName = (*StructName)(nil)
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string]bool: map des interfaces qui ont un compile-time check
func collectCompileTimeChecks(pass *analysis.Pass) map[string]bool {
	checks := make(map[string]bool, initialInterfacesCap)

	// Parcourir tous les fichiers du package
	for _, file := range pass.Files {
		// Parcourir les déclarations
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Vérifier si c'est une déclaration var
			if !ok || genDecl.Tok != token.VAR {
				// Continuer l'itération
				continue
			}

			// Parcourir les specs
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				// Vérifier si c'est une ValueSpec
				if !ok {
					// Continuer l'itération
					continue
				}

				// Extraire le nom de l'interface
				ifaceName := extractInterfaceFromCheck(valueSpec)
				// Ajouter si trouvé
				if ifaceName != "" {
					checks[ifaceName] = true
				}
			}
		}
	}

	// Retour de la map
	return checks
}

// extractInterfaceFromCheck extrait le nom de l'interface d'un compile-time check.
// Pattern: var _ InterfaceName = (*StructName)(nil)
//
// Params:
//   - spec: ValueSpec à analyser
//
// Returns:
//   - string: nom de l'interface ou vide
func extractInterfaceFromCheck(spec *ast.ValueSpec) string {
	// Vérifier que le nom est "_"
	if len(spec.Names) != 1 || spec.Names[0].Name != "_" {
		// Pas le pattern attendu
		return ""
	}

	// Vérifier qu'il y a un type explicite (l'interface)
	if spec.Type == nil {
		// Pas de type explicite
		return ""
	}

	// Extraire le nom de l'interface
	return extractInterfaceNameFromExpr(spec.Type)
}

// extractInterfaceNameFromExpr extrait le nom d'interface d'une expression.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - string: nom de l'interface ou vide
func extractInterfaceNameFromExpr(expr ast.Expr) string {
	// Vérifier le type de l'expression
	switch t := expr.(type) {
	// Identifiant simple (InterfaceName)
	case *ast.Ident:
		// Retour du nom
		return t.Name
	// Sélecteur (pkg.InterfaceName)
	case *ast.SelectorExpr:
		// Retour du sélecteur
		return t.Sel.Name
	}

	// Type non reconnu
	return ""
}

// checkNodeForInterfaceUsage vérifie un nœud AST pour les usages d'interface.
//
// Params:
//   - node: nœud AST à vérifier
//   - usedInterfaces: map pour marquer les interfaces utilisées
func checkNodeForInterfaceUsage(node ast.Node, usedInterfaces map[string]bool) {
	// Verification de la condition
	switch n := node.(type) {
	// Cas FuncDecl - paramètres et retours de fonction
	case *ast.FuncDecl:
		checkFuncDeclForInterfaces(n, usedInterfaces)
	// Cas Field - types de champs (struct, params, etc.)
	case *ast.Field:
		checkType(n.Type, usedInterfaces)
	// Cas InterfaceType - interfaces embarquées
	case *ast.InterfaceType:
		checkEmbeddedInterfaces(n, usedInterfaces)
	// Cas ValueSpec - déclarations de variables (var x MyInterface)
	case *ast.ValueSpec:
		checkValueSpec(n, usedInterfaces)
	// Cas TypeAssertExpr - type assertions (x.(MyInterface))
	case *ast.TypeAssertExpr:
		checkTypeAssert(n, usedInterfaces)
	// Cas TypeSwitchStmt - type switch (switch v := x.(type))
	case *ast.TypeSwitchStmt:
		checkTypeSwitch(n, usedInterfaces)
	// Cas CompositeLit - littéraux composites
	case *ast.CompositeLit:
		checkType(n.Type, usedInterfaces)
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
		// Retour de la fonction
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
//   - compileTimeChecks: map des interfaces vérifiées via var _ I = (*S)(nil)
func reportUnusedInterfaces(pass *analysis.Pass, interfaces map[string]*ast.TypeSpec, usedInterfaces map[string]bool, structNames map[string]bool, compileTimeChecks map[string]bool) {
	// Itération sur les interfaces
	for name, typeSpec := range interfaces {
		// Skip if interface is used directly in code
		if usedInterfaces[name] {
			// Interface utilisée
			continue
		}

		// Skip if interface has a compile-time check (var _ Interface = (*Struct)(nil))
		if compileTimeChecks[name] {
			// Interface vérifiée via compile-time check
			continue
		}

		// Skip if interface follows XXXInterface pattern for struct XXX
		// Ces interfaces sont légitimes pour le mocking de la struct
		if isStructInterfacePattern(name, structNames) {
			// Pattern XXXInterface détecté
			continue
		}

		// Skip if a struct with same name exists (interface for struct)
		// L'interface est légitime car elle permet de mocker la struct
		if hasCorrespondingStruct(name, structNames) {
			// Struct correspondante trouvée
			continue
		}

		// Report based on export status
		reportUnusedInterface(pass, typeSpec, name)
	}
}

// reportUnusedInterface génère le message d'erreur approprié.
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: spécification du type interface
//   - name: nom de l'interface
func reportUnusedInterface(pass *analysis.Pass, typeSpec *ast.TypeSpec, name string) {
	// Message différent selon si l'interface est exportée ou non
	if ast.IsExported(name) {
		// Interface exportée - peut être utilisée par d'autres packages
		pass.Reportf(
			typeSpec.Pos(),
			"KTN-INTERFACE-001: interface '%s' non utilisée dans ce package. "+
				"Si elle est destinée à être utilisée par d'autres packages, ajoutez 'var _ %s = (*StructName)(nil)' "+
				"pour vérifier l'implémentation à la compilation",
			name,
			name,
		)
	} else {
		// Interface privée - message standard
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
	if len(interfaceName) <= interfaceSuffixLen || interfaceName[len(interfaceName)-interfaceSuffixLen:] != "Interface" {
		// Retour de la fonction
		return false
	}

	// Extract potential struct name (remove "Interface" suffix)
	structName := interfaceName[:len(interfaceName)-interfaceSuffixLen]

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
	// Vérifier si l'expression est nil
	if expr == nil {
		// Retour si nil
		return
	}

	// Verification de la condition
	switch t := expr.(type) {
	// Cas Ident - identifiant simple
	case *ast.Ident:
		// Mark identifier as used
		used[t.Name] = true
	// Cas StarExpr - type pointeur
	case *ast.StarExpr:
		// Check pointer type
		checkType(t.X, used)
	// Cas ArrayType - type tableau/slice
	case *ast.ArrayType:
		// Check array element type
		checkType(t.Elt, used)
	// Cas MapType - type map
	case *ast.MapType:
		// Check map key and value types
		checkType(t.Key, used)
		checkType(t.Value, used)
	// Cas ChanType - type channel
	case *ast.ChanType:
		// Check channel element type
		checkType(t.Value, used)
	// Cas SelectorExpr - accès qualifié (pkg.Type)
	case *ast.SelectorExpr:
		// Marquer le sélecteur comme utilisé (ex: pkg.MyInterface)
		used[t.Sel.Name] = true
	}
}

// checkValueSpec vérifie les déclarations de variables pour les usages d'interface.
//
// Params:
//   - spec: spécification de variable
//   - used: map pour marquer les interfaces utilisées
func checkValueSpec(spec *ast.ValueSpec, used map[string]bool) {
	// Vérifier le type déclaré
	if spec.Type != nil {
		checkType(spec.Type, used)
	}
}

// checkTypeAssert vérifie les type assertions pour les usages d'interface.
//
// Params:
//   - expr: expression de type assertion
//   - used: map pour marquer les interfaces utilisées
func checkTypeAssert(expr *ast.TypeAssertExpr, used map[string]bool) {
	// Vérifier le type dans l'assertion (x.(Type))
	if expr.Type != nil {
		checkType(expr.Type, used)
	}
}

// checkTypeSwitch vérifie les type switches pour les usages d'interface.
//
// Params:
//   - stmt: statement type switch
//   - used: map pour marquer les interfaces utilisées
func checkTypeSwitch(stmt *ast.TypeSwitchStmt, used map[string]bool) {
	// Vérifier le corps du switch
	if stmt.Body == nil {
		// Retour si pas de corps
		return
	}

	// Parcourir les cases
	for _, s := range stmt.Body.List {
		caseClause, ok := s.(*ast.CaseClause)
		// Vérifier si c'est un case clause
		if !ok {
			continue
		}

		// Vérifier chaque type dans le case
		for _, typeExpr := range caseClause.List {
			checkType(typeExpr, used)
		}
	}
}
