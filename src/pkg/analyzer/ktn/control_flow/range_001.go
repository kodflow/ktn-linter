package ktn_control_flow

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleRange001 = &analysis.Analyzer{
	Name: "KTN_RANGE_003",
	Doc:  "Détecte la capture de variable de range dans une closure",
	Run:  runRuleRange001,
}

func runRuleRange001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			rangeStmt, ok := n.(*ast.RangeStmt)
			if !ok || rangeStmt.Body == nil {
				return true
			}

			// Récupérer les noms des variables de range
			rangeVars := getRangeVarNames(rangeStmt)
			if len(rangeVars) == 0 {
				return true
			}

			// Chercher les closures dans le body
			ast.Inspect(rangeStmt.Body, func(n ast.Node) bool {
				funcLit, ok := n.(*ast.FuncLit)
				if !ok {
					return true
				}

				// Vérifier si les variables ont été copiées localement
				copiedVars := findCopiedVars(rangeStmt.Body, funcLit, rangeVars)

				// Vérifier si la closure utilise des variables non copiées
				for _, varName := range rangeVars {
					if varName != "_" && !copiedVars[varName] {
						if usesVariable(funcLit.Body, varName) {
							pass.Reportf(funcLit.Pos(),
								"[KTN-RANGE-003] Variable de range '%s' capturée dans une closure.\n"+
									"Les variables de range sont réutilisées à chaque itération.\n"+
									"Capturer directement dans une goroutine/closure cause des bugs.\n"+
									"Créez une copie locale avant la closure.\n"+
									"Exemple:\n"+
									"  // ❌ MAUVAIS\n"+
									"  for _, v := range items {\n"+
									"      go func() { process(v) }()  // BUG!\n"+
									"  }\n"+
									"\n"+
									"  // ✅ CORRECT - copie locale\n"+
									"  for _, v := range items {\n"+
									"      v := v\n"+
									"      go func() { process(v) }()\n"+
									"  }\n"+
									"\n"+
									"  // ✅ CORRECT - passage par paramètre\n"+
									"  for _, v := range items {\n"+
									"      go func(item Item) { process(item) }(v)\n"+
									"  }",
								varName)
						}
					}
				}
				return true
			})
			return true
		})
	}
	return nil, nil
}

func getRangeVarNames(rangeStmt *ast.RangeStmt) []string {
	var names []string
	if rangeStmt.Key != nil {
		if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok {
			names = append(names, keyIdent.Name)
		}
	}
	if rangeStmt.Value != nil {
		if valueIdent, ok := rangeStmt.Value.(*ast.Ident); ok {
			names = append(names, valueIdent.Name)
		}
	}
	return names
}

func findCopiedVars(body *ast.BlockStmt, funcLit *ast.FuncLit, rangeVars []string) map[string]bool {
	copiedVars := make(map[string]bool)
	ast.Inspect(body, func(n ast.Node) bool {
		if n == funcLit {
			return false
		}
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok || assignStmt.Tok.String() != ":=" {
			return true
		}
		for i, lhs := range assignStmt.Lhs {
			if i >= len(assignStmt.Rhs) {
				continue
			}
			lhsIdent, ok1 := lhs.(*ast.Ident)
			rhsIdent, ok2 := assignStmt.Rhs[i].(*ast.Ident)
			if ok1 && ok2 && lhsIdent.Name == rhsIdent.Name {
				copiedVars[lhsIdent.Name] = true
			}
		}
		return true
	})
	return copiedVars
}

func usesVariable(body *ast.BlockStmt, varName string) bool {
	uses := false
	ast.Inspect(body, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if ok && ident.Name == varName {
			uses = true
			return false
		}
		return true
	})
	return uses
}
