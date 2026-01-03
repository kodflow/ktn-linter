// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct001 code de la règle KTN-STRUCT-001
	ruleCodeStruct001 string = "KTN-STRUCT-001"
	// methodsMapCap est la capacité initiale pour la map des méthodes
	methodsMapCap int = 16
)

// Analyzer001 checks getter/setter naming conventions.
// Getters/setters are OPTIONAL, but if present must follow conventions.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct001",
	Doc:      "KTN-STRUCT-001: Convention de nommage des getters/setters (Field() et SetField())",
	Run:      runStruct001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
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

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les structs avec leurs champs privés
	structFields := collectStructPrivateFields(pass, insp, cfg)

	// Collecter les méthodes avec infos détaillées
	methods := collectMethodsDetailed(pass, insp, cfg)

	// Vérifier la convention de nommage
	checkNamingConventions(pass, structFields, methods)

	// Retour de la fonction
	return nil, nil
}

// collectStructPrivateFields collecte les champs privés des structs exportées.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
//
// Returns:
//   - map[string]structFieldsInfo: map nom struct -> champs privés
func collectStructPrivateFields(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string]structFieldsInfo {
	result := make(map[string]structFieldsInfo, methodsMapCap)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec := n.(*ast.TypeSpec)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct001, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Si pas une struct, ignorer
		if !ok {
			// Retour anticipé
			return
		}

		// Ignorer les structs privées
		if !ast.IsExported(typeSpec.Name.Name) {
			// Retour anticipé
			return
		}

		// Collecter les champs privés
		privateFields := make(map[string]bool, methodsMapCap)
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
//   - cfg: configuration
//
// Returns:
//   - map[string][]methodInfo: map receiver -> liste de méthodes
func collectMethodsDetailed(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string][]methodInfo {
	methods := make(map[string][]methodInfo, methodsMapCap)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct001, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une méthode
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Retour anticipé
			return
		}

		// Extraire le nom du receiver
		receiverType := extractReceiverType(funcDecl.Recv.List[0].Type)
		// Si pas de receiver valide, ignorer
		if receiverType == "" {
			// Retour anticipé
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
		// Retour vide
		return ""
	}

	// Si plus d'un résultat, ignorer
	if len(funcDecl.Type.Results.List) > 1 {
		// Retour vide
		return ""
	}

	// Extraire le type
	result := funcDecl.Type.Results.List[0]
	// Utiliser le type info pour obtenir le nom du type
	if tv, ok := pass.TypesInfo.Types[result.Type]; ok {
		// Retour du type trouvé
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
	// Variable pour stocker les infos de struct
	var structInfo structFieldsInfo
	// Parcourir les méthodes
	for receiverType, methodList := range methods {
		// Récupérer les champs de la struct
		structInfo = structFields[receiverType]
		// Si pas de struct connue, ignorer
		if structInfo.privateFields == nil {
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

// checkGetterFieldMismatch vérifie la cohérence entre nom du getter et champ retourné.
//
// Params:
//   - pass: contexte d'analyse
//   - method: informations sur la méthode
//   - structInfo: informations sur la struct
func checkGetterFieldMismatch(pass *analysis.Pass, method methodInfo, structInfo structFieldsInfo) {
	// Ignorer les méthodes sans corps
	if method.funcDecl.Body == nil {
		// Retour anticipé
		return
	}

	// Ignorer les méthodes avec trop de statements
	if len(method.funcDecl.Body.List) != 1 {
		// Retour anticipé
		return
	}

	// Vérifier si c'est un simple return
	retStmt, ok := method.funcDecl.Body.List[0].(*ast.ReturnStmt)
	// Pas un return simple
	if !ok || len(retStmt.Results) != 1 {
		// Retour anticipé
		return
	}

	// Vérifier si le return est un sélecteur sur receiver
	fieldName := extractReturnedField(retStmt.Results[0])
	// Pas un champ retourné
	if fieldName == "" {
		// Retour anticipé
		return
	}

	// Vérifier si le champ est privé
	if !structInfo.privateFields[fieldName] {
		// Retour anticipé
		return
	}

	// Calculer le nom de getter attendu
	expectedGetter := strings.ToUpper(fieldName[:1]) + fieldName[1:]

	// Si le nom ne correspond pas, signaler
	if method.name != expectedGetter && !strings.HasPrefix(method.name, "Get") {
		msg, _ := messages.Get(ruleCodeStruct001)
		pass.Reportf(
			method.funcDecl.Pos(),
			"%s: %s",
			ruleCodeStruct001,
			msg.Format(config.Get().Verbose, method.name, expectedGetter, fieldName),
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
		// Retour vide
		return ""
	}

	// Vérifier si X est un identifiant (receiver)
	if _, ok := sel.X.(*ast.Ident); !ok {
		// Retour vide
		return ""
	}

	// Retour du nom du champ
	return sel.Sel.Name
}
