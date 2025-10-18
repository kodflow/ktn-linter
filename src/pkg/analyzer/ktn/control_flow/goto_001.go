package ktn_control_flow

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var RuleGoto001 = &analysis.Analyzer{
	Name: "KTN_GOTO_001",
	Doc:  "Détecte l'utilisation non idiomatique de goto",
	Run:  runRuleGoto001,
}

func runRuleGoto001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			branchStmt, ok := n.(*ast.BranchStmt)
			if !ok || branchStmt.Tok != token.GOTO {
				return true
			}

			pass.Reportf(branchStmt.Pos(),
				"[KTN-GOTO-001] goto est considéré non idiomatique en Go.\n"+
					"L'utilisation de goto rend le code difficile à comprendre et maintenir.\n"+
					"Utilisez des structures de contrôle standards (if, for, return, break, continue).\n"+
					"Exception: cleanup dans du code bas niveau (rare).\n"+
					"Exemple:\n"+
					"  // ❌ MAUVAIS - goto\n"+
					"  if err != nil { goto cleanup }\n"+
					"  doWork()\n"+
					"  cleanup:\n"+
					"      close()\n"+
					"\n"+
					"  // ✅ CORRECT - defer\n"+
					"  defer close()\n"+
					"  if err != nil { return err }\n"+
					"  doWork()")
			return true
		})
	}
	return nil, nil
}
