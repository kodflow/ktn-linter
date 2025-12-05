// Analyzer 001 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 checks that error is always the last return value
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktnfunc001",
	Doc:      "KTN-FUNC-001: L'erreur doit toujours être en dernière position dans les valeurs de retour",
	Run:      runFunc001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// validateErrorInReturns vérifie que l'erreur est en dernière position.
//
// Params:
//   - pass: contexte d'analyse
//   - funcType: type de la fonction
func validateErrorInReturns(pass *analysis.Pass, funcType *ast.FuncType) {
	// Vérification présence de résultats
	if funcType == nil || funcType.Results == nil {
		// Pas de résultats à vérifier
		return
	}

	results := funcType.Results.List

	// Recherche des positions d'erreur
	var errorPositions []int
	// Itération sur les résultats
	for i, result := range results {
		// Vérification si type error
		if isErrorType(pass, result.Type) {
			errorPositions = append(errorPositions, i)
		}
	}

	// Vérification erreurs mal placées
	if len(errorPositions) > 0 {
		lastPos := len(results) - 1
		// Itération sur les positions d'erreur
		for _, pos := range errorPositions {
			// Vérification position incorrecte
			if pos != lastPos {
				pass.Reportf(
					funcType.Results.Pos(),
					"KTN-FUNC-001: l'erreur doit être en dernière position dans les valeurs de retour (trouvée en position %d sur %d)",
					pos+1,
					len(results),
				)
				// Retour après premier rapport
				return
			}
		}
	}
}

// runFunc001 exécute l'analyse KTN-FUNC-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc001(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		var funcType *ast.FuncType
		// Sélection selon la valeur
		switch node := n.(type) {
		// Traitement FuncDecl
		case *ast.FuncDecl:
			funcType = node.Type
		// Traitement FuncLit
		case *ast.FuncLit:
			funcType = node.Type
		}
		validateErrorInReturns(pass, funcType)
	})

	// Retour de la fonction
	return nil, nil
}

// isErrorType checks if a type expression represents the error interface
// or a type alias for error.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de type à vérifier
//
// Returns:
//   - bool: true si type error ou alias
func isErrorType(pass *analysis.Pass, expr ast.Expr) bool {
	tv := pass.TypesInfo.Types[expr]
	// Vérification si type valide
	if tv.Type == nil {
		return false
	}

	// Vérifier si c'est directement l'interface error
	if isBuiltinError(tv.Type) {
		return true
	}

	// Vérifier si c'est un type nommé qui est l'interface error
	named, ok := tv.Type.(*types.Named)
	if ok {
		obj := named.Obj()
		// Vérification builtin error
		if obj != nil && obj.Name() == "error" && obj.Pkg() == nil {
			return true
		}
		// Vérifier le type sous-jacent (pour les alias)
		if isBuiltinError(named.Underlying()) {
			return true
		}
	}

	// Vérifier si le type implémente l'interface error
	return implementsError(tv.Type)
}

// isBuiltinError vérifie si le type est l'interface error builtin.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si c'est l'interface error
func isBuiltinError(t types.Type) bool {
	// Vérifier si c'est une interface avec la signature Error() string
	iface, ok := t.Underlying().(*types.Interface)
	if !ok {
		return false
	}
	// Vérifier si l'interface a exactement une méthode Error() string
	if iface.NumMethods() != 1 {
		return false
	}
	method := iface.Method(0)
	// Vérifier le nom et la signature
	if method.Name() != "Error" {
		return false
	}
	sig, ok := method.Type().(*types.Signature)
	if !ok {
		return false
	}
	// Vérifier paramètres (aucun) et retour (string)
	if sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	// Vérifier que le retour est string
	basic, ok := sig.Results().At(0).Type().(*types.Basic)
	return ok && basic.Kind() == types.String
}

// implementsError vérifie si un type implémente l'interface error.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si implémente error
func implementsError(t types.Type) bool {
	// Créer l'interface error pour comparaison
	errorMethod := types.NewFunc(0, nil, "Error",
		types.NewSignatureType(nil, nil, nil, nil,
			types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
			false))
	errorIface := types.NewInterfaceType([]*types.Func{errorMethod}, nil)
	errorIface.Complete()

	// Vérifier si le type implémente error
	return types.Implements(t, errorIface)
}
