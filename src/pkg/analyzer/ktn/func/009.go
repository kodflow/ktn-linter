package ktn_func

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule009 vérifie la profondeur d'imbrication des blocs.
//
// KTN-FUNC-009: La profondeur d'imbrication ne doit pas dépasser 3 niveaux.
// Trop d'imbrication rend le code difficile à lire.
//
// Incorrect:
//
//	if a {
//	    if b {
//	        if c {
//	            if d { // 4 niveaux !
//	                ...
//	            }
//	        }
//	    }
//	}
//
// Correct: décomposer en sous-fonctions ou utiliser early returns
var Rule009 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_FUNC_009",
	Doc:  "Vérifie que la profondeur d'imbrication ne dépasse pas 3 niveaux",
	Run:  runRule009,
}

// runRule009 exécute la vérification KTN-FUNC-009.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule009(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkNestingDepth(pass, funcDecl)
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

// checkNestingDepth vérifie la profondeur d'imbrication.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkNestingDepth(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	if funcDecl.Body == nil {
		// Early return from function.
		return
	}

	maxDepth := calculateMaxNestingDepth(funcDecl.Body, 0)
	if maxDepth > 3 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-009] Fonction '%s' a une profondeur d'imbrication trop élevée (%d > 3).\nLimitez l'imbrication à 3 niveaux maximum. Extrayez des sous-fonctions pour réduire la complexité.",
			funcName, maxDepth)
	}
}

// calculateMaxNestingDepth calcule la profondeur maximale d'imbrication.
//
// Params:
//   - node: le nœud AST
//   - currentDepth: la profondeur actuelle
//
// Returns:
//   - int: la profondeur maximale trouvée
func calculateMaxNestingDepth(node ast.Node, currentDepth int) int {
	maxDepth := currentDepth
	ast.Inspect(node, func(n ast.Node) bool {
		maxDepth = inspectNestingNode(n, maxDepth, currentDepth)
		// Early return from function.
		return shouldContinueInspection(n)
	})
	// Early return from function.
	return maxDepth
}

// inspectNestingNode met à jour la profondeur max pour un nœud.
//
// Params:
//   - n: le nœud à inspecter
//   - currentMax: profondeur maximale actuelle
//   - depth: profondeur courante
//
// Returns:
//   - int: nouvelle profondeur maximale
func inspectNestingNode(n ast.Node, currentMax, depth int) int {
	switch stmt := n.(type) {
	case *ast.IfStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
		if stmt.Else != nil {
			currentMax = updateMaxDepth(currentMax, stmt.Else, depth)
		}
	case *ast.ForStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.RangeStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.SwitchStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.SelectStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	}
	// Early return from function.
	return currentMax
}

// shouldContinueInspection détermine si l'inspection doit continuer.
//
// Params:
//   - n: le nœud à vérifier
//
// Returns:
//   - bool: false pour les structures imbriquées
func shouldContinueInspection(n ast.Node) bool {
	switch n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.SelectStmt:
		// Condition not met, return false.
		return false
	}
	// Continue traversing AST nodes.
	return true
}

// updateMaxDepth met à jour la profondeur maximale.
//
// Params:
//   - currentMax: la profondeur maximale actuelle
//   - node: le nœud à analyser
//   - depth: la profondeur actuelle
//
// Returns:
//   - int: la nouvelle profondeur maximale
func updateMaxDepth(currentMax int, node ast.Node, depth int) int {
	newDepth := calculateMaxNestingDepth(node, depth+1)
	if newDepth > currentMax {
		// Early return from function.
		return newDepth
	}
	// Early return from function.
	return currentMax
}
