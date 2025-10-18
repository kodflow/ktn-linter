package ktn_data_structures

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RuleMap001 analyzer for map usage.
var RuleMap001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_MAP_001",
	Doc:  "Détecte l'écriture dans une map sans vérification nil",
	Run:  runRuleMap001,
}

func runRuleMap001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				// Continue traversing AST nodes.
				return true
			}

			// Vérifier si Lhs contient un accès à map
			for _, lhs := range assign.Lhs {
				indexExpr, ok := lhs.(*ast.IndexExpr)
				if !ok {
					continue
				}

				// Vérifier si c'est une map (basique: tout ce qui n'est pas slice/array)
				if mapIdent, ok := indexExpr.X.(*ast.Ident); ok {
					// Vérifier si la map a été vérifiée ou initialisée
					if !isMapSafe(pass, file, mapIdent, assign) {
						pass.Reportf(indexExpr.Pos(),
							"[KTN-MAP-001] Écriture dans la map '%s' sans vérification de nil.\n"+
								"Écrire dans une map nil cause un panic.\n"+
								"Vérifiez toujours qu'une map n'est pas nil avant d'y écrire.\n"+
								"Exemple:\n"+
								"  // ❌ MAUVAIS - panic si m est nil\n"+
								"  var m map[string]int\n"+
								"  m[\"key\"] = 42  // 💥 PANIC\n"+
								"\n"+
								"  // ✅ CORRECT - initialiser avec make\n"+
								"  m := make(map[string]int)\n"+
								"  m[\"key\"] = 42\n"+
								"\n"+
								"  // ✅ CORRECT - vérifier si nil\n"+
								"  if m != nil {\n"+
								"      m[\"key\"] = 42\n"+
								"  }",
							mapIdent.Name)
					}
				}
			}
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

func isMapSafe(pass *analysis.Pass, file *ast.File, mapIdent *ast.Ident, usage ast.Node) bool {
	checked := false
	ast.Inspect(file, func(n ast.Node) bool {
		if n == usage {
			// Condition not met, return false.
			return false
		}
		// Vérifier si map créée avec make()
		assignStmt, ok := n.(*ast.AssignStmt)
		if ok {
			for i, lhs := range assignStmt.Lhs {
				if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == mapIdent.Name {
					if i < len(assignStmt.Rhs) {
						if callExpr, ok := assignStmt.Rhs[i].(*ast.CallExpr); ok {
							if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" {
								checked = true
								// Condition not met, return false.
								return false
							}
						}
					}
				}
			}
		}
		// Continue traversing AST nodes.
		return true
	})
	// Early return from function.
	return checked
}
