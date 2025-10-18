package ktn_data_structures

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleSlice001 = &analysis.Analyzer{
	Name: "KTN_SLICE_001",
	Doc:  "Détecte l'indexation de slice sans vérification de bounds",
	Run:  runRuleSlice001,
}

func runRuleSlice001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			index, ok := n.(*ast.IndexExpr)
			if !ok {
				return true
			}

			// Vérifier si l'index est un littéral (ignoré - peut être vérifié statiquement)
			if _, ok := index.Index.(*ast.BasicLit); ok {
				return true
			}

			// Vérifier si l'index est une variable
			indexIdent, ok := index.Index.(*ast.Ident)
			if !ok {
				return true
			}

			// Vérifier si l'index vient d'un range (sûr)
			if isIndexFromRange(file, indexIdent, index) {
				return true
			}

			// Vérifier si l'index a été vérifié contre len()
			if !isIndexChecked(file, indexIdent, index) {
				if sliceIdent, ok := index.X.(*ast.Ident); ok {
					pass.Reportf(index.Pos(),
						"[KTN-SLICE-001] Indexation du slice '%s' avec '%s' sans vérification de bounds.\n"+
							"Accéder à un index hors limites cause un panic.\n"+
							"Vérifiez toujours que l'index est valide avant d'accéder.\n"+
							"Exemple:\n"+
							"  // ❌ MAUVAIS - panic si i >= len(items)\n"+
							"  v := items[i]\n"+
							"\n"+
							"  // ✅ CORRECT - vérifier les bounds\n"+
							"  if i < len(items) {\n"+
							"      v := items[i]\n"+
							"  }",
						sliceIdent.Name, indexIdent.Name)
				}
			}
			return true
		})
	}
	return nil, nil
}

func isIndexFromRange(file *ast.File, indexIdent *ast.Ident, usage ast.Node) bool {
	fromRange := false
	ast.Inspect(file, func(n ast.Node) bool {
		if n == usage {
			return false
		}
		rangeStmt, ok := n.(*ast.RangeStmt)
		if ok && rangeStmt.Key != nil {
			if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok {
				if keyIdent.Name == indexIdent.Name {
					fromRange = true
					return false
				}
			}
		}
		return true
	})
	return fromRange
}

func isIndexChecked(file *ast.File, indexIdent *ast.Ident, usage ast.Node) bool {
	checked := false
	ast.Inspect(file, func(n ast.Node) bool {
		if n == usage {
			return false
		}
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}
		// Vérifier si la condition est index < len(...)
		if binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			if binaryExpr.Op.String() == "<" {
				if xIdent, ok := binaryExpr.X.(*ast.Ident); ok && xIdent.Name == indexIdent.Name {
					if yCall, ok := binaryExpr.Y.(*ast.CallExpr); ok {
						if yFunc, ok := yCall.Fun.(*ast.Ident); ok && yFunc.Name == "len" {
							checked = true
							return false
						}
					}
				}
			}
		}
		return true
	})
	return checked
}
