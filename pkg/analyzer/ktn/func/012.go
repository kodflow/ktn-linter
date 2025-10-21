package ktnfunc

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer012 checks for unnecessary else blocks after return/continue/break
var Analyzer012 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc012",
	Doc:      "KTN-FUNC-012: Éviter else après return/continue/break (early return préféré)",
	Run:      runFunc012,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc012 exécute l'analyse KTN-FUNC-012.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc012(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
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
				"KTN-FUNC-012: else inutile après %s, utiliser early return",
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
//   - string: type de sortie (return, continue, break)
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
	}

	// Pas de sortie anticipée
	return false, ""
}
