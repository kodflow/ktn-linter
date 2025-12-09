// Analyzer 007 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// METHODS_MAP_CAP est la capacité initiale pour la map des méthodes
	METHODS_MAP_CAP int = 16
)

// Analyzer007 checks getter/setter naming conventions.
// Getters/setters are OPTIONAL, but if present must follow conventions.
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct007",
	Doc:      "KTN-STRUCT-007: Convention de nommage des getters/setters (Field() et SetField())",
	Run:      runStruct007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// methodInfo contient les informations sur une méthode.
type methodInfo struct {
	name       string
	funcDecl   *ast.FuncDecl
	receiverTy string
	returnType string
}

// runStruct007 exécute l'analyse KTN-STRUCT-007.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct007(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les structs avec leurs champs privés
	structFields := collectStructPrivateFields(pass, insp)

	// Collecter les méthodes avec infos détaillées
	methods := collectMethodsDetailed(pass, insp)

	// Vérifier la convention de nommage
	checkNamingConventions(pass, structFields, methods)

	// Retour de la fonction
	return nil, nil
}

// structFieldsInfo contient les infos sur les champs d'une struct.
type structFieldsInfo struct {
	name          string
	privateFields map[string]bool
	pos           ast.Node
}

// collectStructPrivateFields collecte les champs privés des structs exportées.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]structFieldsInfo: map nom struct -> champs privés
func collectStructPrivateFields(pass *analysis.Pass, insp *inspector.Inspector) map[string]structFieldsInfo {
	result := make(map[string]structFieldsInfo)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec := n.(*ast.TypeSpec)

		// Vérifier si c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Si pas une struct, ignorer
		if !ok {
			return
		}

		// Ignorer les structs privées
		if !ast.IsExported(typeSpec.Name.Name) {
			return
		}

		// Collecter les champs privés
		privateFields := make(map[string]bool)
		// Vérifier si la struct a des champs
		if structType.Fields != nil {
			// Parcourir les champs
			for _, field := range structType.Fields.List {
				// Parcourir les noms
				for _, name := range field.Names {
					// Vérifier si le champ est privé
					if len(name.Name) > 0 && unicode.IsLower(rune(name.Name[0])) {
						privateFields[name.Name] = true
					}
				}
			}
		}

		// Stocker les infos
		result[typeSpec.Name.Name] = structFieldsInfo{
			name:          typeSpec.Name.Name,
			privateFields: privateFields,
			pos:           typeSpec,
		}
	})

	// Retour de la map
	return result
}

// collectMethodsDetailed collecte les méthodes avec informations détaillées.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string][]methodInfo: map receiver -> liste de méthodes
func collectMethodsDetailed(pass *analysis.Pass, insp *inspector.Inspector) map[string][]methodInfo {
	methods := make(map[string][]methodInfo, METHODS_MAP_CAP)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Vérifier si c'est une méthode
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			return
		}

		// Extraire le nom du receiver
		receiverType := extractReceiverType(funcDecl.Recv.List[0].Type)
		// Si pas de receiver valide, ignorer
		if receiverType == "" {
			return
		}

		// Extraire le type de retour simple
		returnType := extractSimpleReturnType(pass, funcDecl)

		// Ajouter la méthode
		methods[receiverType] = append(methods[receiverType], methodInfo{
			name:       funcDecl.Name.Name,
			funcDecl:   funcDecl,
			receiverTy: receiverType,
			returnType: returnType,
		})
	})

	// Retour de la map des méthodes
	return methods
}

// extractSimpleReturnType extrait le type de retour simple d'une méthode.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - string: type de retour ou vide
func extractSimpleReturnType(pass *analysis.Pass, funcDecl *ast.FuncDecl) string {
	// Pas de résultats
	if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) == 0 {
		return ""
	}

	// Si plus d'un résultat, ignorer
	if len(funcDecl.Type.Results.List) > 1 {
		return ""
	}

	// Extraire le type
	result := funcDecl.Type.Results.List[0]
	// Utiliser le type info pour obtenir le nom du type
	if tv, ok := pass.TypesInfo.Types[result.Type]; ok {
		return tv.Type.String()
	}

	// Retour vide si pas de type info
	return ""
}

// extractReceiverType extrait le nom du type receiver.
//
// Params:
//   - expr: expression du type receiver
//
// Returns:
//   - string: nom du type ou vide
func extractReceiverType(expr ast.Expr) string {
	// Gérer le cas *Type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}

	// Extraire l'identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom du type
		return ident.Name
	}

	// Retour chaîne vide si pas d'identifiant
	return ""
}

// checkNamingConventions vérifie les conventions de nommage des getters/setters.
// Note: La détection du préfixe Get est gérée par STRUCT-003, on ne la duplique pas ici.
//
// Params:
//   - pass: contexte d'analyse
//   - structFields: map des structs avec champs privés
//   - methods: map des méthodes par receiver
func checkNamingConventions(pass *analysis.Pass, structFields map[string]structFieldsInfo, methods map[string][]methodInfo) {
	// Parcourir les méthodes
	for receiverType, methodList := range methods {
		// Récupérer les champs de la struct
		structInfo, ok := structFields[receiverType]
		// Si pas de struct connue, ignorer
		if !ok {
			continue
		}

		// Vérifier chaque méthode
		for _, method := range methodList {
			// Vérifier la cohérence getter/champ (mismatch nom getter vs champ retourné)
			// Note: préfixe Get est vérifié par STRUCT-003, pas ici
			checkGetterFieldMismatch(pass, method, structInfo)
		}
	}
}

// checkGetPrefixGetter vérifie qu'un getter n'a pas le préfixe Get.
//
// Params:
//   - pass: contexte d'analyse
//   - method: informations sur la méthode
//   - structInfo: informations sur la struct
func checkGetPrefixGetter(pass *analysis.Pass, method methodInfo, structInfo structFieldsInfo) {
	// Vérifier si c'est un getter avec préfixe Get
	if !strings.HasPrefix(method.name, "Get") {
		return
	}

	// Extraire le nom du champ attendu
	fieldName := strings.TrimPrefix(method.name, "Get")
	// Convertir en lowercase pour le champ privé
	privateField := strings.ToLower(fieldName[:1]) + fieldName[1:]

	// Vérifier si le champ existe
	if structInfo.privateFields[privateField] {
		pass.Reportf(
			method.funcDecl.Pos(),
			"KTN-STRUCT-007: getter '%s()' devrait être nommé '%s()' (sans préfixe Get)",
			method.name,
			fieldName,
		)
	}
}

// checkGetterFieldMismatch vérifie la cohérence entre nom du getter et champ retourné.
//
// Params:
//   - pass: contexte d'analyse
//   - method: informations sur la méthode
//   - structInfo: informations sur la struct
func checkGetterFieldMismatch(pass *analysis.Pass, method methodInfo, structInfo structFieldsInfo) {
	// Ignorer les méthodes sans corps
	if method.funcDecl.Body == nil {
		return
	}

	// Ignorer les méthodes avec trop de statements
	if len(method.funcDecl.Body.List) != 1 {
		return
	}

	// Vérifier si c'est un simple return
	retStmt, ok := method.funcDecl.Body.List[0].(*ast.ReturnStmt)
	// Pas un return simple
	if !ok || len(retStmt.Results) != 1 {
		return
	}

	// Vérifier si le return est un sélecteur sur receiver
	fieldName := extractReturnedField(retStmt.Results[0])
	// Pas un champ retourné
	if fieldName == "" {
		return
	}

	// Vérifier si le champ est privé
	if !structInfo.privateFields[fieldName] {
		return
	}

	// Calculer le nom de getter attendu
	expectedGetter := strings.ToUpper(fieldName[:1]) + fieldName[1:]

	// Si le nom ne correspond pas, signaler
	if method.name != expectedGetter && !strings.HasPrefix(method.name, "Get") {
		pass.Reportf(
			method.funcDecl.Pos(),
			"KTN-STRUCT-007: getter '%s()' retourne le champ '%s', devrait être nommé '%s()'",
			method.name,
			fieldName,
			expectedGetter,
		)
	}
}

// extractReturnedField extrait le nom du champ retourné.
//
// Params:
//   - expr: expression retournée
//
// Returns:
//   - string: nom du champ ou vide
func extractReturnedField(expr ast.Expr) string {
	// Vérifier si c'est un sélecteur
	sel, ok := expr.(*ast.SelectorExpr)
	// Pas un sélecteur
	if !ok {
		return ""
	}

	// Vérifier si X est un identifiant (receiver)
	if _, ok := sel.X.(*ast.Ident); !ok {
		return ""
	}

	// Retour du nom du champ
	return sel.Sel.Name
}

// setterInfo contient les informations sur un setter.
type setterInfo struct {
	name     string
	funcDecl *ast.FuncDecl
}

// checkSetterNaming vérifie le nommage des setters.
//
// Params:
//   - pass: contexte d'analyse
//   - method: informations sur la méthode
//   - structInfo: informations sur la struct
func checkSetterNaming(pass *analysis.Pass, method methodInfo, structInfo structFieldsInfo) {
	// Ignorer les méthodes sans corps
	if method.funcDecl.Body == nil {
		return
	}

	// Vérifier si la méthode modifie un champ
	modifiedField := findModifiedField(method.funcDecl.Body)
	// Pas de champ modifié
	if modifiedField == "" {
		return
	}

	// Vérifier si le champ est privé
	if !structInfo.privateFields[modifiedField] {
		return
	}

	// Calculer le nom de setter attendu
	expectedSetter := "Set" + strings.ToUpper(modifiedField[:1]) + modifiedField[1:]

	// Si le nom ne correspond pas, signaler
	if method.name != expectedSetter {
		pass.Reportf(
			method.funcDecl.Pos(),
			"KTN-STRUCT-007: setter pour '%s' devrait être nommé '%s()', pas '%s()'",
			modifiedField,
			expectedSetter,
			method.name,
		)
	}
}

// findModifiedField cherche un champ modifié dans le corps de la fonction.
//
// Params:
//   - body: corps de la fonction
//
// Returns:
//   - string: nom du champ modifié ou vide
func findModifiedField(body *ast.BlockStmt) string {
	// Vérifier si le corps est nil
	if body == nil {
		return ""
	}

	// Parcourir les statements
	for _, stmt := range body.List {
		// Vérifier si c'est une assignation
		assign, ok := stmt.(*ast.AssignStmt)
		// Pas une assignation
		if !ok {
			continue
		}

		// Vérifier si le LHS est un sélecteur
		if len(assign.Lhs) > 0 {
			// Extraire le champ assigné
			if sel, ok := assign.Lhs[0].(*ast.SelectorExpr); ok {
				// Vérifier si X est un identifiant (receiver)
				if _, ok := sel.X.(*ast.Ident); ok {
					// Retour du nom du champ
					return sel.Sel.Name
				}
			}
		}
	}

	// Retour vide si pas de champ modifié
	return ""
}

// isSetterMethod vérifie si une méthode est un setter.
//
// Params:
//   - method: informations sur la méthode
//
// Returns:
//   - bool: true si setter
//   - string: nom du champ setté
func isSetterMethod(method methodInfo) (bool, string) {
	// Vérifier le préfixe Set
	if !strings.HasPrefix(method.name, "Set") {
		return false, ""
	}

	// Extraire le nom du champ
	fieldName := strings.TrimPrefix(method.name, "Set")
	// Convertir en lowercase
	if len(fieldName) > 0 {
		fieldName = strings.ToLower(fieldName[:1]) + fieldName[1:]
	}

	// Vérifier si la méthode a un paramètre et pas de retour
	if method.funcDecl.Type.Params == nil || len(method.funcDecl.Type.Params.List) == 0 {
		return false, ""
	}

	// Retour true avec le nom du champ
	return true, fieldName
}

// Utility function to capitalize first letter
func capitalizeFirst(s string) string {
	// Vérifier si la chaîne est vide
	if len(s) == 0 {
		return ""
	}
	// Retour de la chaîne capitalisée
	return strings.ToUpper(s[:1]) + s[1:]
}

// hasGetter vérifie si un getter existe pour un champ.
//
// Params:
//   - methods: liste des méthodes
//   - fieldName: nom du champ
//
// Returns:
//   - bool: true si getter trouvé
func hasGetter(methods []methodInfo, fieldName string) bool {
	expectedName := capitalizeFirst(fieldName)
	// Parcourir les méthodes
	for _, m := range methods {
		// Vérifier le nom exact ou avec préfixe Get
		if m.name == expectedName {
			return true
		}
	}
	// Retour false si pas trouvé
	return false
}

// Placeholder to use types package
var _ types.Type = nil
