package ktn_control_flow

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleFor001 = &analysis.Analyzer{
	Name: "KTN_FOR_001",
	Doc:  "Détecte l'utilisation inutile de _ dans les range loops",
	Run:  runRuleFor001,
}

func runRuleFor001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			rangeStmt, ok := n.(*ast.RangeStmt)
			if !ok {
				return true
			}

			// Vérifier si la valeur est ignorée avec _
			if rangeStmt.Value != nil {
				valueIdent, ok := rangeStmt.Value.(*ast.Ident)
				if ok && valueIdent.Name == "_" {
					// Cas 1: for _, _ := range
					if rangeStmt.Key != nil {
						if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok && keyIdent.Name == "_" {
							pass.Reportf(rangeStmt.Pos(),
								"[KTN-FOR-001] Index et valeur de range ignorés inutilement avec _.\n"+
									"Si vous n'utilisez ni l'index ni la valeur, omettez-les complètement.\n"+
									"Utilisez `for range items` au lieu de `for _, _ := range items`.\n"+
									"Exemple:\n"+
									"  // ❌ MAUVAIS\n"+
									"  for _, _ := range items { doSomething() }\n"+
									"\n"+
									"  // ✅ CORRECT\n"+
									"  for range items { doSomething() }")
							return true
						}
					}

					// Cas 2: for i, _ := range
					pass.Reportf(rangeStmt.Pos(),
						"[KTN-FOR-001] Valeur de range ignorée inutilement avec _.\n"+
							"Si vous n'utilisez que l'index, omettez la valeur.\n"+
							"Utilisez `for i := range items` au lieu de `for i, _ := range items`.\n"+
							"Exemple:\n"+
							"  // ❌ MAUVAIS\n"+
							"  for i, _ := range items { use(i) }\n"+
							"\n"+
							"  // ✅ CORRECT\n"+
							"  for i := range items { use(i) }\n"+
							"\n"+
							"Note: `for _, v := range items` est CORRECT pour itérer sur les valeurs.")
				}
			}
			return true
		})
	}
	return nil, nil
}
