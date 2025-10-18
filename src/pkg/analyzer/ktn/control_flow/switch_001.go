package ktn_control_flow

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RuleSwitch001 analyzer for switch statements.
var RuleSwitch001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_SWITCH_005",
	Doc:  "Détecte les switch avec un seul case",
	Run:  runRuleSwitch001,
}

func runRuleSwitch001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switchStmt, ok := n.(*ast.SwitchStmt)
			if !ok || switchStmt.Body == nil {
				// Continue traversing AST nodes.
				return true
			}

			caseCount := 0
			for _, stmt := range switchStmt.Body.List {
				if caseClause, ok := stmt.(*ast.CaseClause); ok {
					// default case a List == nil
					if caseClause.List != nil {
						caseCount++
					}
				}
			}

			if caseCount == 1 {
				pass.Reportf(switchStmt.Pos(),
					"[KTN-SWITCH-005] switch avec un seul case devrait être un if.\n"+
						"Un switch avec un seul case est moins lisible qu'un simple if.\n"+
						"Utilisez if/else au lieu de switch pour 1-2 cas.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS\n"+
						"  switch x {\n"+
						"  case 1:\n"+
						"      doSomething()\n"+
						"  default:\n"+
						"      doDefault()\n"+
						"  }\n"+
						"\n"+
						"  // ✅ CORRECT\n"+
						"  if x == 1 {\n"+
						"      doSomething()\n"+
						"  } else {\n"+
						"      doDefault()\n"+
						"  }")
			}
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}
