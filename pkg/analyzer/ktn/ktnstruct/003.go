// Analyzer 003 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// initialStructTypesCap initial capacity for struct types map
	initialStructTypesCap int = 32
	// getPrefixLen length of "Get" prefix
	getPrefixLen int = 3
)

// Analyzer003 checks getters don't have "Get" prefix (Go idiom)
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct003",
	Doc:      "KTN-STRUCT-003: Les getters ne doivent pas avoir le préfixe 'Get' (convention Go idiomatique)",
	Run:      runStruct003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct003 exécute l'analyse KTN-STRUCT-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct003(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les types struct pour pouvoir vérifier les méthodes
	structTypes := collectStructTypes(pass)

	// Filtrer les déclarations de fonctions
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir les déclarations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le noeud suivant
			return
		}

		// Vérifier si c'est une méthode (a un receiver)
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas une méthode, continuer
			return
		}

		// Vérifier si la méthode est exportée
		if !ast.IsExported(funcDecl.Name.Name) {
			// Méthode privée, continuer
			return
		}

		// Vérifier si la struct receiver est exportée
		receiverTypeName := getReceiverTypeName(funcDecl.Recv.List[0].Type)
		// Ignorer les méthodes sur structs privées
		if receiverTypeName != "" && !ast.IsExported(receiverTypeName) {
			// Struct privée, continuer
			return
		}

		// Vérifier si le nom commence par "Get"
		methodName := funcDecl.Name.Name
		// Vérifier préfixe Get
		if !strings.HasPrefix(methodName, "Get") {
			// Ne commence pas par Get, continuer
			return
		}

		// Vérifier qu'il y a au moins un caractère après "Get"
		if len(methodName) <= getPrefixLen {
			// Juste "Get", pas un getter
			return
		}

		// Vérifier si c'est un getter simple (retourne un champ de la struct)
		if !isSimpleGetter(funcDecl, structTypes) {
			// Pas un getter simple, continuer
			return
		}

		// Construire le nom suggéré sans le préfixe "Get"
		suggestedName := methodName[getPrefixLen:]

		// Reporter la violation
		pass.Reportf(
			funcDecl.Name.Pos(),
			"KTN-STRUCT-003: la méthode '%s' devrait être renommée '%s' (convention Go idiomatique, voir https://go.dev/doc/effective_go#Getters)",
			methodName,
			suggestedName,
		)
	})

	// Retour de la fonction
	return nil, nil
}

// collectStructTypes collecte les types struct et leurs champs.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]string: map du nom du type vers la liste des champs
func collectStructTypes(pass *analysis.Pass) map[string][]string {
	result := make(map[string][]string, initialStructTypesCap)

	// Parcourir chaque fichier
	for _, file := range pass.Files {
		// Parcourir l'AST du fichier
		ast.Inspect(file, func(n ast.Node) bool {
			// Vérifier si c'est une déclaration de type
			typeSpec, ok := n.(*ast.TypeSpec)
			// Si ce n'est pas une TypeSpec, continuer
			if !ok {
				// Continuer traversal
				return true
			}

			// Vérifier si c'est une struct
			structType, ok := typeSpec.Type.(*ast.StructType)
			// Si ce n'est pas une struct, continuer
			if !ok {
				// Continuer traversal
				return true
			}

			// Collecter les noms des champs
			var fields []string
			// Vérifier si la struct a des champs
			if structType.Fields != nil {
				// Parcourir les champs
				for _, field := range structType.Fields.List {
					// Parcourir les noms de champs
					for _, name := range field.Names {
						fields = append(fields, strings.ToLower(name.Name))
					}
				}
			}

			// Stocker le type et ses champs
			result[typeSpec.Name.Name] = fields

			// Continuer traversal
			return true
		})
	}

	// Retourner le résultat
	return result
}

// isSimpleGetter vérifie si une méthode est un getter simple.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - structTypes: map des types struct et leurs champs
//
// Returns:
//   - bool: true si c'est un getter simple
func isSimpleGetter(funcDecl *ast.FuncDecl, structTypes map[string][]string) bool {
	// Vérifier que la méthode a un retour
	if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) == 0 {
		// Pas de retour, pas un getter
		return false
	}

	// Vérifier que la méthode n'a pas de paramètres
	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
		// A des paramètres, pas un getter simple
		return false
	}

	// Extraire le nom du type receiver
	receiverType := getReceiverTypeName(funcDecl.Recv.List[0].Type)
	// Vérifier si le type est connu
	if receiverType == "" {
		// Type inconnu, considérer comme getter
		return true
	}

	// Vérifier si le champ correspondant existe dans la struct
	methodName := funcDecl.Name.Name
	expectedFieldName := strings.ToLower(methodName[getPrefixLen:])

	// Obtenir les champs du type receiver
	fields, ok := structTypes[receiverType]
	// Si type non trouvé, considérer comme getter
	if !ok {
		// Type non trouvé
		return true
	}

	// Vérifier si le champ existe (comparaison directe)
	if slices.Contains(fields, expectedFieldName) {
		// Champ trouvé, c'est un getter
		return true
	}

	// Champ non trouvé, mais c'est quand même une méthode Get sans paramètres
	return true
}

// getReceiverTypeName extrait le nom du type du receiver.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - string: nom du type
func getReceiverTypeName(expr ast.Expr) string {
	// Gérer le cas *Type
	starExpr, ok := expr.(*ast.StarExpr)
	// Si c'est un pointeur
	if ok {
		// Récursion pour obtenir le type sous-jacent
		return getReceiverTypeName(starExpr.X)
	}

	// Gérer le cas Type
	ident, ok := expr.(*ast.Ident)
	// Si c'est un identifiant
	if ok {
		// Retourner le nom
		return ident.Name
	}

	// Type non reconnu
	return ""
}
