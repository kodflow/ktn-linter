package ktn_func

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule007 vérifie la complexité cyclomatique des fonctions.
//
// KTN-FUNC-007: La complexité cyclomatique ne doit pas dépasser 10 (50 pour tests).
// La complexité est calculée: 1 + nombre de points de décision (if, for, case, &&, ||).
//
// Incorrect: fonction avec 15 points de décision
// Correct: fonction avec moins de 10 points de décision, ou décomposée
var Rule007 = &analysis.Analyzer{
	Name: "KTN_FUNC_007",
	Doc:  "Vérifie que la complexité cyclomatique reste sous 10 (50 pour tests)",
	Run:  runRule007,
}

// runRule007 exécute la vérification KTN-FUNC-007.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule007(pass *analysis.Pass) (any, error) {
	isTestFile := isTestFile(pass)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkComplexity(pass, funcDecl, isTestFile)
		}
	}

	return nil, nil
}

// checkComplexity vérifie la complexité cyclomatique.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
//   - isTestFile: true si c'est un fichier de test
func checkComplexity(pass *analysis.Pass, funcDecl *ast.FuncDecl, isTestFile bool) {
	funcName := funcDecl.Name.Name
	complexity := calculateCyclomaticComplexity(funcDecl)
	maxComplexity := 10
	if isTestFile {
		maxComplexity = 50
	}

	if complexity > maxComplexity {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-007] Fonction '%s' a une complexité cyclomatique trop élevée (%d > %d).\nRéduisez la complexité en extrayant des sous-fonctions ou en simplifiant la logique.",
			funcName, complexity, maxComplexity)
	}
}

// calculateCyclomaticComplexity calcule la complexité cyclomatique.
//
// Params:
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - int: la complexité cyclomatique
func calculateCyclomaticComplexity(funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		return 1
	}

	complexity := 1
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		complexity += getNodeComplexity(n)
		return true
	})

	return complexity
}

// getNodeComplexity retourne la complexité ajoutée par un nœud.
//
// Params:
//   - n: le nœud AST
//
// Returns:
//   - int: la complexité ajoutée
func getNodeComplexity(n ast.Node) int {
	switch stmt := n.(type) {
	case *ast.IfStmt:
		return 1
	case *ast.ForStmt, *ast.RangeStmt:
		return 1
	case *ast.CaseClause:
		if stmt.List != nil {
			return 1
		}
	case *ast.CommClause:
		if stmt.Comm != nil {
			return 1
		}
	case *ast.BinaryExpr:
		if stmt.Op == token.LAND || stmt.Op == token.LOR {
			return 1
		}
	}
	return 0
}
