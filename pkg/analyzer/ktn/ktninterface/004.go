// Package ktninterface provides analyzers for interface-related lint rules.
package ktninterface

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeInterface004 code de la règle KTN-INTERFACE-004
	ruleCodeInterface004 string = "KTN-INTERFACE-004"
)

// Analyzer004 checks for overuse of empty interface (interface{} or any).
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktninterface004",
	Doc:      "KTN-INTERFACE-004: Utilisation excessive de interface{}/any",
	Run:      runInterface004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface004 exécute l'analyse KTN-INTERFACE-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runInterface004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeInterface004) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeInterface004, filename) {
			// Fichier exclu
			return
		}

		// Vérifier les paramètres de la fonction
		checkFuncParams(pass, funcDecl)

		// Vérifier les retours de la fonction
		checkFuncReturns(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncParams vérifie les paramètres pour interface{}.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
func checkFuncParams(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Vérifier si paramètres présents
	if funcDecl.Type.Params == nil {
		// Pas de paramètres
		return
	}

	// Parcourir les paramètres
	for _, param := range funcDecl.Type.Params.List {
		// Vérifier si interface vide
		if isEmptyInterface(param.Type) {
			// Construire le nom du paramètre
			paramName := buildParamName(param, funcDecl.Name.Name)
			reportEmptyInterface(pass, param, paramName)
		}
	}
}

// checkFuncReturns vérifie les retours pour interface{}.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
func checkFuncReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Vérifier si retours présents
	if funcDecl.Type.Results == nil {
		// Pas de retours
		return
	}

	// Parcourir les retours
	for i, result := range funcDecl.Type.Results.List {
		// Vérifier si interface vide
		if isEmptyInterface(result.Type) {
			// Construire le nom du retour
			returnName := buildReturnName(result, i, funcDecl.Name.Name)
			reportEmptyInterface(pass, result, returnName)
		}
	}
}

// isEmptyInterface vérifie si un type est interface{} ou any.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - bool: true si interface vide
func isEmptyInterface(expr ast.Expr) bool {
	// Vérifier interface{}
	if iface, ok := expr.(*ast.InterfaceType); ok {
		// Interface vide si pas de méthodes
		return iface.Methods == nil || len(iface.Methods.List) == 0
	}

	// Vérifier any (alias de interface{})
	if ident, ok := expr.(*ast.Ident); ok {
		// Comparer avec "any"
		return ident.Name == "any"
	}

	// Pas interface vide
	return false
}

// buildParamName construit le nom d'un paramètre.
//
// Params:
//   - param: champ du paramètre
//   - funcName: nom de la fonction
//
// Returns:
//   - string: nom descriptif du paramètre
func buildParamName(param *ast.Field, funcName string) string {
	// Si le paramètre a un nom
	if len(param.Names) > 0 {
		// Retourner le nom
		return param.Names[0].Name
	}

	// Sinon utiliser le nom de fonction
	return "paramètre de " + funcName
}

// buildReturnName construit le nom d'un retour.
//
// Params:
//   - result: champ du retour
//   - index: index du retour
//   - funcName: nom de la fonction
//
// Returns:
//   - string: nom descriptif du retour
func buildReturnName(result *ast.Field, index int, funcName string) string {
	// Si le retour a un nom
	if len(result.Names) > 0 {
		// Retourner le nom
		return result.Names[0].Name
	}

	// Sinon utiliser une description
	if index == 0 {
		// Premier retour
		return "retour de " + funcName
	}

	// Retour indexé
	return "retour de " + funcName
}

// reportEmptyInterface signale l'utilisation de interface{}.
//
// Params:
//   - pass: contexte d'analyse
//   - field: champ concerné
//   - name: nom descriptif
func reportEmptyInterface(pass *analysis.Pass, field *ast.Field, name string) {
	msg, _ := messages.Get(ruleCodeInterface004)
	pass.Reportf(
		field.Pos(),
		"%s: %s",
		ruleCodeInterface004,
		msg.Format(config.Get().Verbose, name),
	)
}
