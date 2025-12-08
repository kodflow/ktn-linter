// Analyzer 007 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// METHODS_MAP_CAP est la capacité initiale pour la map des méthodes
	METHODS_MAP_CAP int = 16
)

// Analyzer007 checks that non-DTO structs have getters/setters for private fields
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct007",
	Doc:      "KTN-STRUCT-007: Les structs non-DTO doivent avoir des getters pour les champs privés",
	Run:      runStruct007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// nonDTOStructInfo contient les informations sur une struct non-DTO.
type nonDTOStructInfo struct {
	name          string
	privateFields []string
	pos           ast.Node
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

	// Collecter les structs non-DTO avec champs privés
	structs := collectNonDTOStructs(insp)

	// Collecter les méthodes existantes
	methods := collectMethods(insp)

	// Vérifier les getters manquants
	checkMissingGetters(pass, structs, methods)

	// Retour de la fonction
	return nil, nil
}

// collectNonDTOStructs collecte les structs non-DTO avec champs privés.
//
// Params:
//   - insp: inspecteur AST
//
// Returns:
//   - []nonDTOStructInfo: liste des structs avec leurs champs privés
func collectNonDTOStructs(insp *inspector.Inspector) []nonDTOStructInfo {
	var structs []nonDTOStructInfo

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

		// Ignorer les DTOs
		if shared.IsSerializableStruct(structType, typeSpec.Name.Name) {
			return
		}

		// Collecter les champs privés
		privateFields := collectPrivateFields(structType)
		// Si pas de champs privés, ignorer
		if len(privateFields) == 0 {
			return
		}

		structs = append(structs, nonDTOStructInfo{
			name:          typeSpec.Name.Name,
			privateFields: privateFields,
			pos:           typeSpec,
		})
	})

	// Retour de la liste des structs collectées
	return structs
}

// collectPrivateFields collecte les noms des champs privés.
//
// Params:
//   - structType: type de la struct
//
// Returns:
//   - []string: liste des noms de champs privés
func collectPrivateFields(structType *ast.StructType) []string {
	var fields []string

	// Vérifier si la struct a des champs
	if structType.Fields == nil {
		// Retour liste vide si pas de champs
		return fields
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Parcourir les noms
		for _, name := range field.Names {
			// Vérifier si le champ est privé
			if len(name.Name) > 0 && unicode.IsLower(rune(name.Name[0])) {
				fields = append(fields, name.Name)
			}
		}
	}

	// Retour de la liste des champs privés
	return fields
}

// collectMethods collecte les méthodes par type receiver.
//
// Params:
//   - insp: inspecteur AST
//
// Returns:
//   - map[string][]string: map receiver -> liste de méthodes
func collectMethods(insp *inspector.Inspector) map[string][]string {
	methods := make(map[string][]string, METHODS_MAP_CAP)

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

		// Ajouter la méthode
		methods[receiverType] = append(methods[receiverType], funcDecl.Name.Name)
	})

	// Retour de la map des méthodes
	return methods
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

// checkMissingGetters vérifie les getters manquants.
//
// Params:
//   - pass: contexte d'analyse
//   - structs: liste des structs à vérifier
//   - methods: map des méthodes existantes
func checkMissingGetters(pass *analysis.Pass, structs []nonDTOStructInfo, methods map[string][]string) {
	// Parcourir les structs
	for _, s := range structs {
		structMethods := methods[s.name]

		// Vérifier chaque champ privé
		for _, field := range s.privateFields {
			// Construire le nom du getter attendu
			expectedGetter := buildGetterName(field)

			// Vérifier si le getter existe
			if !hasMethod(structMethods, expectedGetter) {
				pass.Reportf(
					s.pos.Pos(),
					"KTN-STRUCT-007: la struct '%s' devrait avoir un getter '%s()' pour le champ privé '%s'",
					s.name,
					expectedGetter,
					field,
				)
			}
		}
	}
}

// buildGetterName construit le nom du getter pour un champ.
//
// Params:
//   - fieldName: nom du champ privé
//
// Returns:
//   - string: nom du getter (ex: "name" -> "Name" ou "GetName")
func buildGetterName(fieldName string) string {
	// Capitaliser le premier caractère
	if len(fieldName) == 0 {
		// Retour chaîne vide si nom vide
		return ""
	}
	// En Go, le getter standard est juste le nom capitalisé (sans "Get")
	return strings.ToUpper(fieldName[:1]) + fieldName[1:]
}

// hasMethod vérifie si une méthode existe.
//
// Params:
//   - methods: liste des méthodes
//   - name: nom recherché
//
// Returns:
//   - bool: true si la méthode existe
func hasMethod(methods []string, name string) bool {
	// Parcourir les méthodes
	for _, m := range methods {
		// Vérifier le nom exact ou avec préfixe Get
		if m == name || m == "Get"+name {
			// Retour true si méthode trouvée
			return true
		}
	}
	// Retour false si méthode non trouvée
	return false
}
