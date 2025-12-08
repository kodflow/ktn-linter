// Analyzer 003 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks for unnecessary else blocks after return/continue/break/panic
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc003",
	Doc:      "KTN-FUNC-003: Éviter else après return/continue/break/panic (early return préféré)",
	Run:      runFunc003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc003 exécute l'analyse KTN-FUNC-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc003(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		ifStmt := n.(*ast.IfStmt)

		// Vérifier si le bloc if est vide
		if ifStmt.Body == nil || len(ifStmt.Body.List) == 0 {
			// Bloc vide, ignorer
			return
		}

		// Vérifier s'il y a une clause else
		if ifStmt.Else == nil {
			// Pas de else, ignorer
			return
		}

		// Récupérer la dernière instruction du bloc if
		lastStmt := ifStmt.Body.List[len(ifStmt.Body.List)-1]

		// Déterminer le type de sortie anticipée
		hasEarlyExit, exitType := checkEarlyExit(lastStmt)

		// Si sortie anticipée détectée, reporter l'erreur
		if hasEarlyExit {
			pass.Reportf(
				ifStmt.Else.Pos(),
				"KTN-FUNC-003: else inutile après %s, utiliser early return",
				exitType,
			)
		}
	})

	// Retour succès
	return nil, nil
}

// checkEarlyExit vérifie si une instruction est une sortie anticipée.
//
// Params:
//   - stmt: instruction à vérifier
//
// Returns:
//   - bool: true si sortie anticipée
//   - string: type de sortie (return, continue, break, panic)
func checkEarlyExit(stmt ast.Stmt) (bool, string) {
	// Switch sur le type d'instruction
	switch s := stmt.(type) {
	// Cas return
	case *ast.ReturnStmt:
		// Retour true car c'est une sortie anticipée de type return
		return true, "return"
	// Cas branch (continue/break)
	case *ast.BranchStmt:
		// Si c'est continue
		if s.Tok.String() == "continue" {
			// Retour true car c'est une sortie anticipée de type continue
			return true, "continue"
		}
		// Si c'est break
		if s.Tok.String() == "break" {
			// Retour true car c'est une sortie anticipée de type break
			return true, "break"
		}
	// Cas expression statement (peut contenir panic)
	case *ast.ExprStmt:
		// Vérifier si c'est un appel à panic
		if isPanicCall(s.X) {
			// Retour true car c'est une sortie anticipée de type panic
			return true, "panic"
		}
	}

	// Pas de sortie anticipée
	return false, ""
}

// isPanicCall vérifie si une expression est un appel à panic().
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un appel à panic
func isPanicCall(expr ast.Expr) bool {
	// Vérifier si c'est un appel de fonction
	call, ok := expr.(*ast.CallExpr)
	// Si pas un appel de fonction
	if !ok {
		// Retour false
		return false
	}

	// Vérifier si c'est un identifiant
	ident, ok := call.Fun.(*ast.Ident)
	// Si pas un identifiant
	if !ok {
		// Retour false
		return false
	}

	// Retour true si c'est panic
	return ident.Name == "panic"
}
