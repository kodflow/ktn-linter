package ktnfunc

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer011 checks that all branches, returns, and significant logic blocks have comments
var Analyzer011 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc011",
	Doc:      "KTN-FUNC-011: Tous les blocs de contrôle (if/else/switch/for), returns et logique significative doivent être commentés",
	Run:      runFunc011,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc011 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc011(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Skip test functions (Test*, Benchmark*, Example*, Fuzz*)
		funcName := funcDecl.Name.Name
  // Vérification de la condition
		if isTestFunction(funcName) {
   // Retour de la fonction
			return
		}

		// Check all statements in the function body
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
   // Sélection selon la valeur
			switch stmt := node.(type) {
   // Traitement
			case *ast.IfStmt:
				checkIfStmt(pass, stmt)
   // Traitement
			case *ast.SwitchStmt:
				checkSwitchStmt(pass, stmt)
   // Traitement
			case *ast.TypeSwitchStmt:
				checkTypeSwitchStmt(pass, stmt)
   // Traitement
			case *ast.ForStmt, *ast.RangeStmt:
				checkLoopStmt(pass, stmt)
   // Traitement
			case *ast.ReturnStmt:
				checkReturnStmt(pass, stmt)
			}
   // Retour de la fonction
			return true
		})
	})

 // Retour de la fonction
	return nil, nil
}

// checkIfStmt checks if an if statement has a comment
// Params:
//   - pass: contexte d'analyse
//
func checkIfStmt(pass *analysis.Pass, stmt *ast.IfStmt) {
 // Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc 'if' doit avoir un commentaire explicatif")
	}

	// Check else clause
	if stmt.Else != nil {
		// For else, check if there's a comment before or at the start of the else block
		if !hasCommentBeforeOrInside(pass, stmt.Else) {
			pass.Reportf(stmt.Else.Pos(), "KTN-FUNC-011: le bloc 'else' doit avoir un commentaire explicatif")
		}
	}
}

// checkSwitchStmt checks if a switch statement and its cases have comments
// Params:
//   - pass: contexte d'analyse
//
func checkSwitchStmt(pass *analysis.Pass, stmt *ast.SwitchStmt) {
 // Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc 'switch' doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
  // Itération sur les éléments
		for _, caseClause := range stmt.Body.List {
   // Vérification de la condition
			if clause, ok := caseClause.(*ast.CaseClause); ok {
    // Vérification de la condition
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-FUNC-011: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkTypeSwitchStmt checks if a type switch statement and its cases have comments
// Params:
//   - pass: contexte d'analyse
//
func checkTypeSwitchStmt(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) {
 // Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc 'switch' (type) doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
  // Itération sur les éléments
		for _, caseClause := range stmt.Body.List {
   // Vérification de la condition
			if clause, ok := caseClause.(*ast.CaseClause); ok {
    // Vérification de la condition
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-FUNC-011: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkLoopStmt checks if a loop statement has a comment
// Params:
//   - pass: contexte d'analyse
//
func checkLoopStmt(pass *analysis.Pass, stmt ast.Node) {
 // Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc de boucle doit avoir un commentaire explicatif")
	}
}

// checkReturnStmt checks if a return statement has a comment
// Params:
//   - pass: contexte d'analyse
//
func checkReturnStmt(pass *analysis.Pass, stmt *ast.ReturnStmt) {
 // Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) && !hasInlineComment(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le 'return' doit avoir un commentaire explicatif")
	}
}

// hasCommentBefore checks if there's a comment on the line before the given position
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si commentaire présent
//
func hasCommentBefore(pass *analysis.Pass, pos token.Pos) bool {
	position := pass.Fset.Position(pos)
	filename := position.Filename

	// Find the file in the pass
	var file *ast.File
 // Itération sur les éléments
	for _, f := range pass.Files {
  // Vérification de la condition
		if pass.Fset.Position(f.Pos()).Filename == filename {
			file = f
			break
		}
	}

	// Check all comments in the file
	for _, commentGroup := range file.Comments {
		commentPos := pass.Fset.Position(commentGroup.End())
		// Comment should be on the line immediately before the statement
		if commentPos.Filename == filename && commentPos.Line == position.Line-1 {
     // Retour de la fonction
			return true
		}
	}

 // Retour de la fonction
	return false
}

// hasInlineComment checks if there's a comment on the same line as the given position
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si commentaire inline
//
func hasInlineComment(pass *analysis.Pass, pos token.Pos) bool {
	position := pass.Fset.Position(pos)
	filename := position.Filename

	// Find the file in the pass
	var file *ast.File
 // Itération sur les éléments
	for _, f := range pass.Files {
  // Vérification de la condition
		if pass.Fset.Position(f.Pos()).Filename == filename {
			file = f
			break
		}
	}

	// Check all comments in the file
	for _, commentGroup := range file.Comments {
  // Itération sur les éléments
		for _, comment := range commentGroup.List {
			commentPos := pass.Fset.Position(comment.Pos())
			// Comment should be on the same line as the statement
			if commentPos.Filename == filename && commentPos.Line == position.Line {
     // Retour de la fonction
				return true
			}
		}
	}

 // Retour de la fonction
	return false
}

// hasCommentBeforeOrInside checks if there's a comment before or at the start of an else block
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si commentaire avant
//
func hasCommentBeforeOrInside(pass *analysis.Pass, elseStmt ast.Stmt) bool {
	// Check if there's a comment before the else
	if hasCommentBefore(pass, elseStmt.Pos()) {
  // Retour de la fonction
		return true
	}

	// If it's a block statement, check for a comment at the beginning
	if blockStmt, ok := elseStmt.(*ast.BlockStmt); ok && len(blockStmt.List) > 0 {
		firstStmt := blockStmt.List[0]
  // Vérification de la condition
		if hasCommentBefore(pass, firstStmt.Pos()) {
   // Retour de la fonction
			return true
		}
	}

 // Retour de la fonction
	return false
}
