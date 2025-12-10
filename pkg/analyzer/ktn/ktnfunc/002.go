// Analyzer 002 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc002 is the rule code for this analyzer
	ruleCodeFunc002 string = "KTN-FUNC-002"
)

// Analyzer002 checks that context.Context is always the first parameter
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc002",
	Doc:      "KTN-FUNC-002: context.Context doit toujours être le premier paramètre (après le receiver pour les méthodes)",
	Run:      runFunc002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc002 exécute l'analyse KTN-FUNC-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc002) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeFunc002, filename) {
			// Fichier exclu
			return
		}

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		// Analyze function parameters
		analyzeContextParams(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeContextParams analyzes context.Context parameters in a function.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: fonction à analyser
func analyzeContextParams(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	// Vérification de la condition
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		// Retour de la fonction
		return
	}

	// Check each parameter and count contexts
	contextParamIndex := -1
	contextCount := 0
	// Itération sur les éléments
	for i, field := range funcDecl.Type.Params.List {
		// Vérification de la condition
		if isContextTypeWithPass(pass, field.Type) {
			// Compte le nombre de noms dans ce champ (pour gérer a, b context.Context)
			nameCount := len(field.Names)
			// Si pas de noms explicites, compter 1
			if nameCount == 0 {
				// Assignation de 1 si pas de noms
				nameCount = 1
			}
			// Ajout au compteur de contextes
			contextCount += nameCount
			// Enregistrer la première position trouvée
			if contextParamIndex == -1 {
				// Enregistrement de la position
				contextParamIndex = i
			}
		}
	}

	// Report multiple contexts
	reportMultipleContexts(pass, funcDecl, funcName, contextCount)

	// Report misplaced context
	reportMisplacedContext(pass, funcDecl, funcName, contextParamIndex)
}

// reportMultipleContexts reports if function has multiple context.Context params.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: fonction à analyser
//   - funcName: nom de la fonction
//   - contextCount: nombre de contextes
func reportMultipleContexts(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, contextCount int) {
	// Signaler si plus d'un context.Context
	if contextCount > 1 {
		// Rapport d'erreur pour contextes multiples
		pass.Reportf(
			funcDecl.Type.Params.Pos(),
			"KTN-FUNC-002: la fonction '%s' a %d paramètres context.Context, ce qui est inhabituel",
			funcName,
			contextCount,
		)
	}
}

// reportMisplacedContext reports if context.Context is not the first parameter.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: fonction à analyser
//   - funcName: nom de la fonction
//   - contextParamIndex: position du contexte
func reportMisplacedContext(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, contextParamIndex int) {
	// If there's a context parameter and it's not first, report error
	if contextParamIndex > 0 {
		// Rapport d'erreur pour position incorrecte
		pass.Reportf(
			funcDecl.Type.Params.List[contextParamIndex].Pos(),
			"KTN-FUNC-002: context.Context doit être le premier paramètre de la fonction '%s'",
			funcName,
		)
	}
}

// isContextTypeWithPass checks if a type is context.Context using type info.
// Handles both direct context.Context and type aliases.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de type à vérifier
//
// Returns:
//   - bool: true si type context.Context ou alias
func isContextTypeWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	// Try using type info first
	tv := pass.TypesInfo.Types[expr]
	// Si le type est disponible, vérifier via types.Type
	if tv.Type != nil {
		// Retour du résultat de la vérification par type
		return isContextTypeByType(tv.Type)
	}
	// Fallback to AST-based check
	return isContextType(expr)
}

// isContextTypeByType checks if a types.Type is context.Context.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si context.Context
func isContextTypeByType(t types.Type) bool {
	// Get the named type
	named, ok := t.(*types.Named)
	// Si ce n'est pas un type nommé, ce n'est pas context.Context
	if !ok {
		// Retour false si pas un type nommé
		return false
	}
	obj := named.Obj()
	// Si pas d'objet associé, invalide
	if obj == nil {
		// Retour false si pas d'objet
		return false
	}
	// Check if it's from context package
	if isContextObj(obj) {
		// Retour true si objet context.Context
		return true
	}
	// Check underlying type for aliases
	return isContextUnderlying(t, obj)
}

// isContextObj vérifie si l'objet est context.Context.
//
// Params:
//   - obj: objet de type à vérifier
//
// Returns:
//   - bool: true si context.Context
func isContextObj(obj *types.TypeName) bool {
	// Vérifier le package et le nom
	return obj.Pkg() != nil && obj.Pkg().Path() == "context" && obj.Name() == "Context"
}

// isContextUnderlying vérifie le type sous-jacent pour les alias.
//
// Params:
//   - t: type à vérifier
//   - obj: objet associé au type
//
// Returns:
//   - bool: true si type sous-jacent est context.Context
func isContextUnderlying(t types.Type, obj *types.TypeName) bool {
	// Si le package n'est pas défini, pas d'alias possible
	if obj.Pkg() == nil {
		// Retour false si package non défini
		return false
	}
	underlying := t.Underlying()
	// Vérifier si le type sous-jacent est nommé
	underNamed, ok := underlying.(*types.Named)
	// Si pas un type nommé, retour false
	if !ok {
		// Retour false si type sous-jacent pas nommé
		return false
	}
	underObj := underNamed.Obj()
	// Vérifier context.Context sous-jacent
	if underObj != nil && underObj.Pkg() != nil {
		// Retour si match avec context.Context
		return underObj.Pkg().Path() == "context" && underObj.Name() == "Context"
	}
	// Retour false par défaut
	return false
}

// isContextType checks if a type is context.Context (AST-based fallback).
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si type context.Context
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	// Vérification de la condition
	if !ok {
		// Retour de la fonction
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	// Retour de la fonction
	return ok && ident.Name == "context" && sel.Sel.Name == "Context"
}
