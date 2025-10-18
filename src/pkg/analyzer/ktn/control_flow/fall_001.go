package ktn_control_flow

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var RuleFall001 = &analysis.Analyzer{
	Name: "KTN_FALL_001",
	Doc:  "Détecte fallthrough hors d'un switch",
	Run:  runRuleFall001,
}

func runRuleFall001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			branchStmt, ok := n.(*ast.BranchStmt)
			if !ok || branchStmt.Tok != token.FALLTHROUGH {
				return true
			}

			if !isInsideSwitchCase(file, branchStmt) {
				pass.Reportf(branchStmt.Pos(),
					"[KTN-FALL-001] fallthrough ne peut être utilisé que dans un switch.\n"+
						"Le mot-clé fallthrough est uniquement valide dans un case de switch.\n"+
						"L'utiliser ailleurs est une erreur de syntaxe.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS - fallthrough hors switch\n"+
						"  if x > 0 {\n"+
						"      fallthrough  // ERREUR\n"+
						"  }\n"+
						"\n"+
						"  // ✅ CORRECT - fallthrough dans switch\n"+
						"  switch x {\n"+
						"  case 1:\n"+
						"      doOne()\n"+
						"      fallthrough\n"+
						"  case 2:\n"+
						"      doTwo()\n"+
						"  }")
			}
			return true
		})
	}
	return nil, nil
}

func isInsideSwitchCase(file *ast.File, target ast.Node) bool {
	inCase := false
	ast.Inspect(file, func(n ast.Node) bool {
		if switchStmt, ok := n.(*ast.SwitchStmt); ok {
			if switchStmt.Body != nil {
				for _, stmt := range switchStmt.Body.List {
					if caseClause, ok := stmt.(*ast.CaseClause); ok {
						for _, s := range caseClause.Body {
							if s == target {
								inCase = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	})
	return inCase
}
