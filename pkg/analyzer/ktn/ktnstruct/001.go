// Package ktnstruct provides analyzers for struct-related KTN rules.
//
// DEPRECATED: KTN-STRUCT-001 (Analyzer001)
// Cette règle est dépréciée et remplacée par KTN-API-001.
// KTN-STRUCT-001 imposait des "mirror interfaces" (interfaces reprenant 100% des méthodes d'une struct),
// ce qui est un anti-pattern car cela crée un couplage fort au lieu du découplage souhaité.
//
// Le bon pattern est l'Interface Segregation Principle (ISP):
// - Définir des interfaces MINIMALES côté CONSUMER
// - Chaque consumer définit l'interface dont IL a besoin
// - Une struct peut satisfaire plusieurs interfaces différentes
//
// Voir KTN-API-001 pour la règle de remplacement.
package ktnstruct

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	// ruleCodeStruct001 code de la règle KTN-STRUCT-001
	ruleCodeStruct001 string = "KTN-STRUCT-001"
)

// Analyzer001 vérifie qu'une interface existe pour chaque struct avec méthodes publiques.
// DEPRECATED: Cette règle est dépréciée et remplacée par KTN-API-001.
// Les "mirror interfaces" sont un anti-pattern - utilisez KTN-API-001 pour les interfaces minimales côté consumer.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct001",
	Doc:      "KTN-STRUCT-001: [DEPRECATED] Remplacé par KTN-API-001 - mirror interfaces sont un anti-pattern",
	Run:      runStruct001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// structWithMethods stocke une struct et ses méthodes publiques
type structWithMethods struct {
	name       string
	node       *ast.TypeSpec
	structType *ast.StructType
	methods    []shared.MethodSignature
	namedType  *types.Named
}

// runStruct001 exécute l'analyse KTN-STRUCT-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct001(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct001) {
		// Règle désactivée
		return nil, nil
	}

	// Collecter tous les interface checks du package (pas juste fichier)
	interfaceChecks := collectInterfaceChecksWithTypes(pass)

	// Collecter toutes les méthodes du PACKAGE entier (pas juste par fichier)
	// Car en Go, les méthodes peuvent être définies dans différents fichiers
	allMethodsByStruct := collectAllMethodsByStruct(pass)

	// Parcourir chaque fichier du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct001, filename) {
			// Fichier exclu
			continue
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le fichier suivant
			continue
		}

		// Collecter les interfaces locales et leurs méthodes
		localInterfaces := collectInterfaces(file, pass)

		// Collecter les structs du fichier avec méthodes du package entier
		structs := collectStructsFromFile(file, allMethodsByStruct)

		// Vérifier chaque struct
		checkStructs(pass, structs, localInterfaces, interfaceChecks)
	}

	// Retour de la fonction
	return nil, nil
}

// checkStructs vérifie chaque struct et rapporte les violations.
//
// Params:
//   - pass: contexte d'analyse
//   - structs: liste des structs avec méthodes
//   - localInterfaces: interfaces locales du fichier
//   - interfaceChecks: compile-time checks du package
func checkStructs(pass *analysis.Pass, structs []structWithMethods, localInterfaces map[string][]shared.MethodSignature, interfaceChecks []interfaceCheck) {
	// Vérifier chaque struct
	for _, s := range structs {
		// Si la struct n'a pas de méthodes publiques, skip
		if len(s.methods) == 0 {
			// Continuer avec la struct suivante
			continue
		}

		// Exception: les DTOs n'ont pas besoin d'interface
		if shared.IsSerializableStruct(s.structType, s.name) {
			// DTO - pas besoin d'interface
			continue
		}

		// Exception: les consommateurs (structs avec injection de dépendances)
		// Un consommateur utilise des interfaces injectées
		if isConsumerStruct(s.structType, pass) {
			// Consumer pattern - pas besoin d'interface
			continue
		}

		// Vérifier si une interface locale couvre toutes les méthodes
		if hasMatchingInterface(s, localInterfaces) {
			// Interface locale trouvée
			continue
		}

		// Vérifier via compile-time checks avec go/types
		if hasMatchingInterfaceCheck(s, interfaceChecks) {
			// Interface externe valide trouvée
			continue
		}

		// Aucune interface ne couvre toutes les méthodes
		msg, _ := messages.Get(ruleCodeStruct001)
		pass.Reportf(
			s.node.Pos(),
			"%s: %s",
			ruleCodeStruct001,
			msg.Format(config.Get().Verbose, s.name, len(s.methods), s.name),
		)
	}
}

// isConsumerStruct détecte si une struct est un "consommateur" (Service, Controller, Handler).
// Un consommateur utilise l'injection de dépendances via des champs de type interface.
// Ces structs orchestrent les dépendances et n'ont pas besoin de leur propre interface car:
// - Elles ne sont jamais mockées (c'est elles qu'on teste)
// - Elles n'ont qu'une seule implémentation
// - Elles sont des points d'entrée de la logique métier
//
// Params:
//   - structType: type de la struct à analyser
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si la struct est un consommateur
func isConsumerStruct(structType *ast.StructType, pass *analysis.Pass) bool {
	// Vérifier si la struct a des champs de type interface
	if structType.Fields == nil {
		// Pas de champs
		return false
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Obtenir le type du champ via TypesInfo
		if pass.TypesInfo != nil {
			fieldType := pass.TypesInfo.TypeOf(field.Type)
			// Vérifier si c'est une interface
			if fieldType != nil {
				// Vérifier le type sous-jacent (pour les types nommés)
				if _, isInterface := fieldType.Underlying().(*types.Interface); isInterface {
					// Champ de type interface trouvé → c'est un consommateur
					return true
				}
			}
		}
	}

	// Pas de champ interface → pas un consommateur
	return false
}

// collectInterfaceChecksWithTypes collecte les compile-time interface checks avec types.
// Patterns supportés: var _ I = (*S)(nil), var _ I = S{}, var _ I = new(S), var _ I = &S{}
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - []interfaceCheck: liste des checks trouvés
func collectInterfaceChecksWithTypes(pass *analysis.Pass) []interfaceCheck {
	var checks []interfaceCheck

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

				// Extraire le check avec types
				check := extractInterfaceCheckWithTypes(valueSpec, pass)
				// Ajouter si valide
				if check != nil {
					checks = append(checks, *check)
				}
			}
		}
	}

	// Retour de la liste
	return checks
}

// extractInterfaceCheckWithTypes extrait un interface check avec informations de type.
//
// Params:
//   - spec: ValueSpec à analyser
//   - pass: contexte d'analyse
//
// Returns:
//   - *interfaceCheck: check trouvé ou nil
func extractInterfaceCheckWithTypes(spec *ast.ValueSpec, pass *analysis.Pass) *interfaceCheck {
	// Vérifier que le nom est "_"
	if len(spec.Names) != 1 || spec.Names[0].Name != "_" {
		// Pas le pattern attendu
		return nil
	}

	// Vérifier qu'il y a un type explicite (l'interface)
	if spec.Type == nil {
		// Pas de type explicite
		return nil
	}

	// Obtenir le type de l'interface via TypesInfo
	ifaceTypeInfo := pass.TypesInfo.TypeOf(spec.Type)
	// Vérifier si on a pu résoudre le type
	if ifaceTypeInfo == nil {
		// Type non résolu
		return nil
	}

	// Extraire l'interface sous-jacente
	ifaceType, ok := ifaceTypeInfo.Underlying().(*types.Interface)
	// Vérifier si c'est une interface
	if !ok {
		// Pas une interface
		return nil
	}

	// Vérifier qu'il y a une valeur
	if len(spec.Values) != 1 {
		// Pas de valeur
		return nil
	}

	// Extraire le nom de la struct depuis la valeur
	structName := extractStructNameFromValue(spec.Values[0], pass)
	// Vérifier si on a trouvé un nom
	if structName == "" {
		// Nom non trouvé
		return nil
	}

	// Retour du check
	return &interfaceCheck{
		structName:    structName,
		interfaceType: ifaceType,
	}
}

// extractStructNameFromValue extrait le nom de la struct depuis une expression.
// Supporte: (*S)(nil), S{}, new(S), &S{}
//
// Params:
//   - expr: expression à analyser
//   - pass: contexte d'analyse
//
// Returns:
//   - string: nom de la struct ou vide
func extractStructNameFromValue(expr ast.Expr, pass *analysis.Pass) string {
	// Obtenir le type de l'expression via TypesInfo
	exprType := pass.TypesInfo.TypeOf(expr)
	// Vérifier si on a un type
	if exprType == nil {
		// Fallback sur extraction AST
		return extractStructNameFromAST(expr)
	}

	// Extraire le type sous-jacent (déréférencer les pointeurs)
	return extractStructNameFromType(exprType)
}

// extractStructNameFromType extrait le nom de struct depuis un types.Type.
//
// Params:
//   - t: type à analyser
//
// Returns:
//   - string: nom de la struct ou vide
func extractStructNameFromType(t types.Type) string {
	// Déréférencer les pointeurs
	for {
		ptr, ok := t.(*types.Pointer)
		// Sortir si ce n'est pas un pointeur
		if !ok {
			// Sortir de la boucle
			break
		}
		t = ptr.Elem()
	}

	// Extraire le Named type
	named, ok := t.(*types.Named)
	// Vérifier si c'est un Named type
	if !ok {
		// Pas un Named type
		return ""
	}

	// Retour du nom
	return named.Obj().Name()
}

// extractStructNameFromAST extrait le nom de struct depuis l'AST (fallback).
// Supporte: (*S)(nil), S{}, new(S), &S{}
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - string: nom de la struct ou vide
func extractStructNameFromAST(expr ast.Expr) string {
	// Vérification selon le type d'expression
	switch e := expr.(type) {
	// Pattern: (*S)(nil) ou new(S)
	case *ast.CallExpr:
		// Extraction depuis CallExpr
		return extractFromCallExpr(e)
	// Pattern: &S{}
	case *ast.UnaryExpr:
		// Extraction depuis UnaryExpr
		return extractFromUnaryExpr(e)
	// Pattern: S{}
	case *ast.CompositeLit:
		// Extraction depuis CompositeLit
		return extractFromCompositeLit(e)
	}

	// Type non supporté
	return ""
}

// extractFromCallExpr extrait le nom depuis (*S)(nil) ou new(S).
//
// Params:
//   - call: CallExpr à analyser
//
// Returns:
//   - string: nom de la struct ou vide
func extractFromCallExpr(call *ast.CallExpr) string {
	// Vérifier si c'est new(S)
	if funIdent, ok := call.Fun.(*ast.Ident); ok && funIdent.Name == "new" {
		// Extraire l'argument
		if len(call.Args) == 1 {
			// Vérifier si l'argument est un identifiant
			if argIdent, ok := call.Args[0].(*ast.Ident); ok {
				// Retour du nom de la struct
				return argIdent.Name
			}
		}
		// Pas d'argument valide
		return ""
	}

	// Pattern: (*S)(nil) - le Fun est (*S)
	parenExpr, ok := call.Fun.(*ast.ParenExpr)
	// Vérifier si c'est une ParenExpr
	if !ok {
		// Pas une ParenExpr
		return ""
	}

	// Extraire *S depuis les parenthèses
	starExpr, ok := parenExpr.X.(*ast.StarExpr)
	// Vérifier si c'est une StarExpr
	if !ok {
		// Pas une StarExpr
		return ""
	}

	// Extraire le nom de la struct
	structIdent, ok := starExpr.X.(*ast.Ident)
	// Vérifier si c'est un identifiant
	if ok {
		// Retour du nom de la struct
		return structIdent.Name
	}

	// Pattern non reconnu
	return ""
}

// extractFromUnaryExpr extrait le nom depuis &S{}.
//
// Params:
//   - unary: UnaryExpr à analyser
//
// Returns:
//   - string: nom de la struct ou vide
func extractFromUnaryExpr(unary *ast.UnaryExpr) string {
	// Vérifier si c'est &
	if unary.Op.String() != "&" {
		// Pas &
		return ""
	}

	// L'opérande doit être un CompositeLit
	compLit, ok := unary.X.(*ast.CompositeLit)
	// Vérifier si c'est un CompositeLit
	if !ok {
		// Pas un CompositeLit
		return ""
	}

	// Extraire depuis le CompositeLit
	return extractFromCompositeLit(compLit)
}

// extractFromCompositeLit extrait le nom depuis S{}.
//
// Params:
//   - comp: CompositeLit à analyser
//
// Returns:
//   - string: nom de la struct ou vide
func extractFromCompositeLit(comp *ast.CompositeLit) string {
	// Le type doit être un identifiant
	ident, ok := comp.Type.(*ast.Ident)
	// Vérifier si c'est un identifiant
	if ok {
		// Retour du nom
		return ident.Name
	}

	// Pas trouvé
	return ""
}

// hasMatchingInterfaceCheck vérifie si un compile-time check couvre toutes les méthodes.
//
// Params:
//   - s: struct avec méthodes
//   - checks: liste des interface checks
//
// Returns:
//   - bool: true si une interface check couvre toutes les méthodes
func hasMatchingInterfaceCheck(s structWithMethods, checks []interfaceCheck) bool {
	// Parcourir les checks
	for _, check := range checks {
		// Vérifier si c'est pour cette struct
		if check.structName != s.name {
			// Pas pour cette struct
			continue
		}

		// Vérifier que l'interface couvre toutes les méthodes publiques
		if interfaceCoversAllPublicMethods(check.interfaceType, s.methods) {
			// Interface valide
			return true
		}
	}

	// Aucun check valide trouvé
	return false
}

// interfaceCoversAllPublicMethods vérifie si l'interface couvre toutes les méthodes.
//
// Params:
//   - iface: type interface
//   - methods: méthodes publiques de la struct
//
// Returns:
//   - bool: true si toutes les méthodes sont couvertes
func interfaceCoversAllPublicMethods(iface *types.Interface, methods []shared.MethodSignature) bool {
	// Chaque méthode publique de la struct doit être dans l'interface
	for _, m := range methods {
		found := false
		// Parcourir les méthodes de l'interface
		for ifaceMethod := range iface.Methods() {
			// Comparer les noms
			if ifaceMethod.Name() == m.Name {
				// Comparer les signatures
				if signaturesMatch(ifaceMethod, m) {
					found = true
					// Sortir de la boucle
					break
				}
			}
		}

		// Si une méthode n'est pas trouvée, l'interface ne couvre pas tout
		if !found {
			// Retour false
			return false
		}
	}

	// Toutes les méthodes sont couvertes
	return true
}

// signaturesMatch compare une méthode d'interface avec une signature de struct.
//
// Params:
//   - ifaceMethod: méthode de l'interface
//   - structMethod: signature de la méthode de struct
//
// Returns:
//   - bool: true si les signatures correspondent
func signaturesMatch(ifaceMethod *types.Func, structMethod shared.MethodSignature) bool {
	// Obtenir la signature
	sig, ok := ifaceMethod.Type().(*types.Signature)
	// Vérifier si c'est une Signature
	if !ok {
		// Pas une signature
		return false
	}

	// Comparer les paramètres
	paramsStr := formatTypeTuple(sig.Params())
	// Comparer les résultats
	resultsStr := formatTypeTuple(sig.Results())

	// Comparer avec la struct method
	return paramsStr == structMethod.ParamsStr && resultsStr == structMethod.ResultsStr
}

// formatTypeTuple formate un tuple de types en string.
//
// Params:
//   - tuple: tuple de types
//
// Returns:
//   - string: représentation string
func formatTypeTuple(tuple *types.Tuple) string {
	// Si nil ou vide
	if tuple == nil || tuple.Len() == 0 {
		// Retour vide
		return ""
	}

	var parts []string
	// Parcourir les éléments
	for v := range tuple.Variables() {
		parts = append(parts, types.TypeString(v.Type(), nil))
	}

	// Retour de la string jointe
	return strings.Join(parts, ",")
}

// collectInterfaces collecte toutes les interfaces et leurs méthodes.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]shared.MethodSignature: map nom interface -> signatures méthodes
func collectInterfaces(file *ast.File, pass *analysis.Pass) map[string][]shared.MethodSignature {
	interfaces := make(map[string][]shared.MethodSignature, 0)

	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est une TypeSpec
		typeSpec, ok := n.(*ast.TypeSpec)
		// Si ce n'est pas une TypeSpec, continuer
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier si c'est une interface
		ifaceType, isInterface := typeSpec.Type.(*ast.InterfaceType)
		// Si ce n'est pas une interface, continuer
		if !isInterface {
			// Continue traversal
			return true
		}

		// Collecter les méthodes de l'interface
		var methods []shared.MethodSignature
		// Parcourir les méthodes de l'interface
		for _, method := range ifaceType.Methods.List {
			// Vérifier si c'est une méthode (pas un embedded interface)
			funcType, isFunc := method.Type.(*ast.FuncType)
			// Si ce n'est pas une fonction, continuer
			if !isFunc {
				// Continue with next method
				continue
			}

			// Extraire le nom de la méthode
			for _, name := range method.Names {
				methods = append(methods, shared.MethodSignature{
					Name:       name.Name,
					ParamsStr:  formatFieldList(funcType.Params, pass),
					ResultsStr: formatFieldList(funcType.Results, pass),
				})
			}
		}

		interfaces[typeSpec.Name.Name] = methods
		// Continue traversal
		return true
	})

	// Retour de la map
	return interfaces
}

// extractStructNameFromReceiver extrait le nom de la struct depuis le receiver.
//
// Params:
//   - recvType: type du receiver
//
// Returns:
//   - string: nom de la struct
func extractStructNameFromReceiver(recvType ast.Expr) string {
	var structName string
	// Gérer les receivers de type *T ou T
	switch t := recvType.(type) {
	// Traitement
	case *ast.StarExpr:
		// Receiver de type *T
		if ident, ok := t.X.(*ast.Ident); ok {
			structName = ident.Name
		}
	// Traitement
	case *ast.Ident:
		// Receiver de type T
		structName = t.Name
	}
	// Retour du nom
	return structName
}

// collectAllMethodsByStruct collecte les méthodes publiques de TOUT le package.
// En Go, les méthodes peuvent être définies dans différents fichiers du même package.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]shared.MethodSignature: map des méthodes par struct
func collectAllMethodsByStruct(pass *analysis.Pass) map[string][]shared.MethodSignature {
	methodsByStruct := make(map[string][]shared.MethodSignature, 0)

	// Parcourir TOUS les fichiers du package
	for _, file := range pass.Files {
		collectMethodsFromFile(file, pass, methodsByStruct)
	}

	// Retour de la map
	return methodsByStruct
}

// collectMethodsFromFile collecte les méthodes d'un fichier et les ajoute à la map.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//   - methodsByStruct: map à enrichir
func collectMethodsFromFile(file *ast.File, pass *analysis.Pass, methodsByStruct map[string][]shared.MethodSignature) {
	// Collecter les méthodes du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier FuncDecl
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification si FuncDecl
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier receiver
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas de receiver
			return true
		}

		// Vérifier méthode publique
		if !ast.IsExported(funcDecl.Name.Name) {
			// Méthode privée
			return true
		}

		// Extraire nom de struct
		recvType := funcDecl.Recv.List[0].Type
		structName := extractStructNameFromReceiver(recvType)

		// Ajouter méthode
		if structName != "" {
			methodsByStruct[structName] = append(methodsByStruct[structName], shared.MethodSignature{
				Name:       funcDecl.Name.Name,
				ParamsStr:  formatFieldList(funcDecl.Type.Params, pass),
				ResultsStr: formatFieldList(funcDecl.Type.Results, pass),
			})
		}

		// Continue traversal
		return true
	})
}

// hasMatchingInterface vérifie si une interface couvre toutes les méthodes.
//
// Params:
//   - s: struct avec méthodes
//   - interfaces: map des interfaces
//
// Returns:
//   - bool: true si une interface matching existe
func hasMatchingInterface(s structWithMethods, interfaces map[string][]shared.MethodSignature) bool {
	// Parcourir toutes les interfaces
	for _, ifaceMethods := range interfaces {
		// Vérifier si cette interface couvre toutes les méthodes de la struct
		if interfaceCoversAllMethods(s.methods, ifaceMethods) {
			// Interface trouvée
			return true
		}
	}

	// Aucune interface ne couvre toutes les méthodes
	return false
}

// interfaceCoversAllMethods vérifie si l'interface couvre toutes les méthodes.
//
// Params:
//   - structMethods: méthodes de la struct
//   - ifaceMethods: méthodes de l'interface
//
// Returns:
//   - bool: true si toutes les méthodes sont couvertes
func interfaceCoversAllMethods(structMethods []shared.MethodSignature, ifaceMethods []shared.MethodSignature) bool {
	// Chaque méthode de la struct doit être dans l'interface
	for _, sm := range structMethods {
		found := false
		// Chercher la méthode dans l'interface
		for _, im := range ifaceMethods {
			// Comparer nom et signatures
			if sm.Name == im.Name && sm.ParamsStr == im.ParamsStr && sm.ResultsStr == im.ResultsStr {
				found = true
				// Sortir de la boucle
				break
			}
		}

		// Si une méthode n'est pas trouvée, l'interface ne couvre pas tout
		if !found {
			// Retour false
			return false
		}
	}

	// Toutes les méthodes sont couvertes
	return true
}

// collectStructsFromFile collecte les structs d'un fichier avec les méthodes pré-collectées.
//
// Params:
//   - file: fichier AST
//   - allMethodsByStruct: map des méthodes par struct (collectées de tout le package)
//
// Returns:
//   - []structWithMethods: liste des structs avec méthodes
func collectStructsFromFile(file *ast.File, allMethodsByStruct map[string][]shared.MethodSignature) []structWithMethods {
	var structs []structWithMethods

	// Parcourir le fichier pour trouver les structs
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est une TypeSpec
		typeSpec, ok := n.(*ast.TypeSpec)
		// Si ce n'est pas une TypeSpec, continuer
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier si c'est une struct
		structType, isStruct := typeSpec.Type.(*ast.StructType)
		// Si c'est une struct
		if isStruct {
			structs = append(structs, structWithMethods{
				name:       typeSpec.Name.Name,
				node:       typeSpec,
				structType: structType,
				methods:    allMethodsByStruct[typeSpec.Name.Name],
			})
		}

		// Continue traversal
		return true
	})

	// Retour de la liste
	return structs
}

// collectMethodsByStruct collecte les méthodes d'un seul fichier.
// Note: utilisé par collectStructsWithMethods pour la compatibilité.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]shared.MethodSignature: map des méthodes par struct
func collectMethodsByStruct(file *ast.File, pass *analysis.Pass) map[string][]shared.MethodSignature {
	methodsByStruct := make(map[string][]shared.MethodSignature, 0)
	collectMethodsFromFile(file, pass, methodsByStruct)
	// Retour de la map
	return methodsByStruct
}

// formatFieldList formate une liste de champs en string.
// Gère correctement les champs avec plusieurs noms (ex: a, b int).
// Utilise pass.TypesInfo pour obtenir le type complet (avec package path).
//
// Params:
//   - fields: liste de champs
//   - pass: contexte d'analyse
//
// Returns:
//   - string: représentation string
func formatFieldList(fields *ast.FieldList, pass *analysis.Pass) string {
	// Si pas de champs
	if fields == nil {
		// Retour vide
		return ""
	}

	var parts []string
	// Parcourir les champs
	for _, field := range fields.List {
		// Utiliser TypesInfo pour obtenir le type résolu
		var typeStr string
		// Vérifier que pass et TypesInfo sont disponibles
		if pass != nil && pass.TypesInfo != nil {
			// Obtenir le type résolu via TypesInfo
			if t := pass.TypesInfo.TypeOf(field.Type); t != nil {
				// Type résolu avec chemin complet
				typeStr = types.TypeString(t, nil)
			} else {
				// Fallback sur ExprString si type non résolu
				typeStr = types.ExprString(field.Type)
			}
		} else {
			// Fallback sans TypesInfo
			typeStr = types.ExprString(field.Type)
		}
		// Compter combien de noms partagent ce type
		count := len(field.Names)
		// Si pas de noms explicites (type anonyme), count = 1
		if count == 0 {
			count = 1
		}
		// Ajouter le type pour chaque nom
		for range count {
			parts = append(parts, typeStr)
		}
	}

	// Retour de la string jointe
	return strings.Join(parts, ",")
}
