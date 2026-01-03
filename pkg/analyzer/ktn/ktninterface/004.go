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
	// ruleCodeInterface004 définit le code de la règle KTN-INTERFACE-004
	ruleCodeInterface004 string = "KTN-INTERFACE-004"
	// expectedAnalyzerRunReturnCount définit le nombre de retours attendus pour une fonction analyzer.Run
	expectedAnalyzerRunReturnCount int = 2
)

// Analyzer004 checks for overuse of empty interface (interface{} or any).
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktninterface004",
	Doc:      "KTN-INTERFACE-004: Utilisation excessive de interface{}/any",
	Run:      runInterface004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface004 exécute l'analyse KTN-INTERFACE-004.
// Tested via integration tests in 004_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runInterface004(pass *analysis.Pass) (any, error) {
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeInterface004) {
		// Règle désactivée - retour immédiat
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir les déclarations de fonction
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeInterface004, filename) {
			// Fichier exclu - skip
			return
		}

		// Vérifier les paramètres et retours
		checkFuncParams(pass, funcDecl)
		checkFuncReturns(pass, funcDecl)
	})

	// Retour succès
	return nil, nil
}

// checkFuncParams vérifie les paramètres pour interface{}.
// Tested via integration tests in 004_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
func checkFuncParams(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Vérifier si la fonction a des paramètres
	if funcDecl.Type.Params == nil {
		// Pas de paramètres - retour
		return
	}

	// Parcourir les paramètres
	for _, param := range funcDecl.Type.Params.List {
		// Vérifier si c'est une interface vide
		if isEmptyInterface(param.Type) {
			// Interface vide détectée - construire nom et reporter
			paramName := buildParamName(param, funcDecl.Name.Name)
			reportEmptyInterface(pass, param, paramName)
		}
	}
}

// checkFuncReturns vérifie les retours pour interface{}.
// Tested via integration tests in 004_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
func checkFuncReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Vérifier si la fonction a des retours
	if funcDecl.Type.Results == nil {
		// Pas de retours - skip
		return
	}

	// Exclure les fonctions analysis.Analyzer.Run (signature: func(*analysis.Pass) (any, error))
	if isAnalyzerRunFunction(funcDecl) {
		// Fonction analysis.Analyzer.Run - skip (signature imposée par le framework)
		return
	}

	// Parcourir les valeurs de retour
	for i, result := range funcDecl.Type.Results.List {
		// Vérifier si c'est une interface vide
		if isEmptyInterface(result.Type) {
			// Interface vide détectée - construire nom et reporter
			returnName := buildReturnName(result, i, funcDecl.Name.Name)
			reportEmptyInterface(pass, result, returnName)
		}
	}
}

// isAnalyzerRunFunction vérifie si la fonction est une implémentation de analysis.Analyzer.Run.
// La signature attendue est: func(*analysis.Pass) (any, error)
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si c'est une fonction analysis.Analyzer.Run
func isAnalyzerRunFunction(funcDecl *ast.FuncDecl) bool {
	// Vérifier le paramètre unique *analysis.Pass
	if !hasAnalysisPassParam(funcDecl) {
		// Pas le bon paramètre
		return false
	}

	// Vérifier les retours (any, error)
	return hasAnyErrorReturns(funcDecl)
}

// hasAnalysisPassParam vérifie si la fonction a un unique paramètre *analysis.Pass.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si le paramètre est *analysis.Pass
func hasAnalysisPassParam(funcDecl *ast.FuncDecl) bool {
	// Vérifier le nombre de paramètres (1 seul)
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) != 1 {
		// Pas le bon nombre de paramètres
		return false
	}

	// Vérifier le type du paramètre
	return isAnalysisPassType(funcDecl.Type.Params.List[0].Type)
}

// isAnalysisPassType vérifie si le type est *analysis.Pass.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - bool: true si c'est *analysis.Pass
func isAnalysisPassType(expr ast.Expr) bool {
	// Vérifier si c'est un pointeur
	starExpr, ok := expr.(*ast.StarExpr)
	// Type non pointeur
	if !ok {
		// Pas un type pointeur
		return false
	}

	// Vérifier si c'est un selector (package.Type)
	selExpr, ok := starExpr.X.(*ast.SelectorExpr)
	// Type non qualifié
	if !ok {
		// Pas un sélecteur de type
		return false
	}

	// Vérifier package et type
	pkgIdent, ok := selExpr.X.(*ast.Ident)
	// Retourne true si c'est analysis.Pass
	return ok && pkgIdent.Name == "analysis" && selExpr.Sel.Name == "Pass"
}

// hasAnyErrorReturns vérifie si la fonction retourne (any, error).
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si les retours sont (any, error)
func hasAnyErrorReturns(funcDecl *ast.FuncDecl) bool {
	// Vérifier le nombre de retours (2: any, error)
	if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) != expectedAnalyzerRunReturnCount {
		// Pas le bon nombre de retours
		return false
	}

	// Vérifier le premier retour (any ou interface{})
	if !isEmptyInterface(funcDecl.Type.Results.List[0].Type) {
		// Premier retour n'est pas any/interface{}
		return false
	}

	// Vérifier le second retour (error)
	return isErrorType(funcDecl.Type.Results.List[1].Type)
}

// isErrorType vérifie si le type est error.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - bool: true si c'est error
func isErrorType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	// Retourne true si c'est le type error
	return ok && ident.Name == "error"
}

// isEmptyInterface vérifie si un type est interface{} ou any.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - bool: true si interface vide
func isEmptyInterface(expr ast.Expr) bool {
	// Vérifier si c'est interface{}
	if iface, ok := expr.(*ast.InterfaceType); ok {
		// Interface vide si pas de méthodes
		return iface.Methods == nil || len(iface.Methods.List) == 0
	}

	// Vérifier si c'est le mot-clé 'any'
	if ident, ok := expr.(*ast.Ident); ok {
		// 'any' est équivalent à interface{}
		return ident.Name == "any"
	}

	// Pas une interface vide
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
	// Vérifier si le paramètre a un nom
	if len(param.Names) > 0 {
		// Utiliser le nom du paramètre
		return param.Names[0].Name
	}

	// Paramètre anonyme - utiliser nom générique
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
	// Vérifier si le retour a un nom
	if len(result.Names) > 0 {
		// Utiliser le nom du retour
		return result.Names[0].Name
	}

	// Premier retour anonyme
	if index == 0 {
		// Retour générique
		return "retour de " + funcName
	}

	// Autre retour anonyme
	return "retour de " + funcName
}

// reportEmptyInterface signale l'utilisation de interface{}.
// Tested via integration tests in 004_external_test.go
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
