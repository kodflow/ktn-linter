package ktn_ops

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RuleReturn001 analyzer for return statements.
var RuleReturn001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_RETURN_004",
	Doc:  "Détecte les naked returns dans des fonctions longues",
	Run:  runRuleReturn001,
}

func runRuleReturn001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Type == nil || funcDecl.Type.Results == nil {
				// Continue traversing AST nodes.
				return true
			}

			// Vérifier si la fonction a des named returns
			hasNamedReturns := false
			for _, field := range funcDecl.Type.Results.List {
				if len(field.Names) > 0 {
					hasNamedReturns = true
					break
				}
			}

			if !hasNamedReturns || funcDecl.Body == nil {
				// Continue traversing AST nodes.
				return true
			}

			// Chercher les naked returns dans le body
			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				returnStmt, ok := n.(*ast.ReturnStmt)
				if !ok {
					// Continue traversing AST nodes.
					return true
				}

				// Vérifier si c'est un naked return
				if len(returnStmt.Results) == 0 {
					// Vérifier si la fonction est longue
					if isFunctionLong(funcDecl) {
						pass.Reportf(returnStmt.Pos(),
							"[KTN-RETURN-001] Naked return dans la fonction '%s'.\n"+
								"Les naked returns (return sans valeurs avec named returns) réduisent la lisibilité.\n"+
								"Dans les fonctions longues, il est difficile de savoir quelles valeurs sont retournées.\n"+
								"Retournez explicitement les valeurs.\n"+
								"Exemple:\n"+
								"  // ❌ MAUVAIS - naked return dans fonction longue\n"+
								"  func process() (result int, err error) {\n"+
								"      // ... 50 lignes ...\n"+
								"      result = 42\n"+
								"      return  // Pas clair!\n"+
								"  }\n"+
								"\n"+
								"  // ✅ CORRECT - return explicite\n"+
								"  func process() (result int, err error) {\n"+
								"      // ... code ...\n"+
								"      return result, err  // Clair!\n"+
								"  }",
							funcDecl.Name.Name)
					}
				}
				// Continue traversing AST nodes.
				return true
			})
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

func isFunctionLong(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Body == nil {
		// Condition not met, return false.
		return false
	}
	statementCount := 0
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ExprStmt, *ast.AssignStmt, *ast.ReturnStmt,
			*ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.SelectStmt, *ast.DeferStmt:
			statementCount++
		}
		// Continue traversing AST nodes.
		return true
	})
	// Early return from function.
	return statementCount > 10
}
