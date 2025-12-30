// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar009 is the rule code for this analyzer
	ruleCodeVar009 string = "KTN-VAR-009"
	// defaultMaxStructFields max fields for struct without pointer
	defaultMaxStructFields int = 3
)

// Analyzer009 checks for large struct passed by value in function parameters
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar009",
	Doc:      "KTN-VAR-009: Utilise des pointeurs pour les structs avec >3 champs",
	Run:      runVar009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar009 exécute l'analyse KTN-VAR-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar009(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar009) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	maxFields := cfg.GetThreshold(ruleCodeVar009, defaultMaxStructFields)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Vérification de l'exclusion
		if cfg.IsFileExcluded(ruleCodeVar009, filename) {
			// Fichier exclu
			return
		}

		// Vérifier les receivers (méthodes)
		if funcDecl.Recv != nil {
			// Analyse des receivers
			checkFuncParams009(pass, funcDecl.Recv, maxFields)
		}

		// Vérifier les paramètres de la fonction
		if funcDecl.Type.Params != nil {
			// Analyse des paramètres
			checkFuncParams009(pass, funcDecl.Type.Params, maxFields)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncParams009 vérifie les paramètres d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - params: liste des paramètres
//   - maxFields: nombre max de champs
func checkFuncParams009(pass *analysis.Pass, params *ast.FieldList, maxFields int) {
	// Parcours des paramètres
	for _, param := range params.List {
		// Utiliser la position du nom si disponible, sinon la position du type
		pos := param.Pos()
		// Vérifier si des noms sont disponibles et non nil
		if len(param.Names) > 0 && param.Names[0] != nil {
			pos = param.Names[0].NamePos
		}
		// Vérification du type de paramètre
		checkParamType009(pass, param.Type, pos, maxFields)
	}
}

// checkParamType009 vérifie le type d'un paramètre.
//
// Params:
//   - pass: contexte d'analyse
//   - typ: type du paramètre
//   - pos: position du paramètre
//   - maxFields: nombre max de champs
func checkParamType009(pass *analysis.Pass, typ ast.Expr, pos token.Pos, maxFields int) {
	// Ignorer les pointeurs (déjà passés par référence)
	if _, isPointer := typ.(*ast.StarExpr); isPointer {
		// C'est un pointeur, OK
		return
	}

	// Récupération du type réel
	typeInfo := pass.TypesInfo.TypeOf(typ)
	// Vérification du type
	if typeInfo == nil {
		// Type inconnu
		return
	}

	// Ignorer les types externes (frameworks comme Terraform)
	if isExternalType009(typeInfo, pass) {
		// Retour de la fonction
		return
	}

	// Vérification que c'est une struct
	structType, ok := typeInfo.Underlying().(*types.Struct)
	// Vérification du type struct
	if !ok {
		// Pas une struct
		return
	}

	// Comptage des champs
	numFields := structType.NumFields()
	// Vérification du nombre de champs
	if numFields > maxFields {
		// Grande struct détectée
		msg, _ := messages.Get(ruleCodeVar009)
		pass.Reportf(
			pos,
			"%s: %s",
			ruleCodeVar009,
			msg.Format(config.Get().Verbose, numFields),
		)
	}
}

// isExternalType009 checks if type is from external package.
//
// Params:
//   - typeInfo: Type to check
//   - pass: Analysis pass
//
// Returns:
//   - bool: true if type is from external package
func isExternalType009(typeInfo types.Type, pass *analysis.Pass) bool {
	// Check if it's a named type
	named, ok := typeInfo.(*types.Named)
	// Verification de la condition
	if !ok {
		// Retour de la fonction
		return false
	}

	// Get package of the type
	obj := named.Obj()
	// Verification de la condition
	if obj == nil || obj.Pkg() == nil {
		// Retour de la fonction
		return false
	}

	// Check if type is from current package
	return obj.Pkg().Path() != pass.Pkg.Path()
}
