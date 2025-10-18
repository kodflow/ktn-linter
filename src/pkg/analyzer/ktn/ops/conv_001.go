package ktn_ops

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleConv001 = &analysis.Analyzer{
	Name: "KTN_CONV_002",
	Doc:  "Détecte les conversions de type redondantes",
	Run:  runRuleConv001,
}

func runRuleConv001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 1 {
				return true
			}

			// Vérifier si Fun est un type
			typeIdent, ok := call.Fun.(*ast.Ident)
			if !ok {
				return true
			}

			// Vérifier si l'argument est le même identifiant
			if argIdent, ok := call.Args[0].(*ast.Ident); ok {
				if typeIdent.Name == argIdent.Name {
					pass.Reportf(call.Pos(),
						"[KTN-CONV-001] Conversion de type redondante: %s(%s).\n"+
							"Cette conversion est inutile car la variable a déjà le bon type.\n"+
							"Les conversions redondantes nuisent à la lisibilité.\n"+
							"Exemple:\n"+
							"  // ❌ MAUVAIS\n"+
							"  var x int = 5\n"+
							"  y := int(x)  // x est déjà int!\n"+
							"\n"+
							"  // ✅ CORRECT\n"+
							"  var x int = 5\n"+
							"  y := x",
						typeIdent.Name, argIdent.Name)
				}
			}
			return true
		})
	}
	return nil, nil
}
