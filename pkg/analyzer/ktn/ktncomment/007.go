// Analyzer 007 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeComment007 is the rule code for this analyzer
	ruleCodeComment007 string = "KTN-COMMENT-007"
)

// Analyzer007 checks that all branches, returns, and significant logic blocks have comments
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment007",
	Doc:      "KTN-COMMENT-007: Tous les blocs de contrôle (if/else/switch/for), returns et logique significative doivent être commentés",
	Run:      runComment007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment007 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runComment007(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment007) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Check if function should be skipped
		if shouldSkipFunction(pass, cfg, funcDecl) {
			// Function should be skipped
			return
		}

		// Check all statements in the function body
		checkFunctionBody(pass, funcDecl.Body)
	})

	// Retour de la fonction
	return nil, nil
}

// shouldSkipFunction checks if a function should be skipped from analysis.
//
// Params:
//   - pass: analysis pass
//   - cfg: configuration
//   - funcDecl: function declaration to check
//
// Returns:
//   - bool: true if function should be skipped
func shouldSkipFunction(
	pass *analysis.Pass,
	cfg *config.Config,
	funcDecl *ast.FuncDecl,
) bool {
	// Skip test files entirely
	filename := pass.Fset.Position(funcDecl.Pos()).Filename

	// Skip excluded files
	if cfg.IsFileExcluded(ruleCodeComment007, filename) {
		// File excluded by configuration
		return true
	}

	// Skip test files
	if shared.IsTestFile(filename) {
		// Test file should be skipped
		return true
	}

	// Skip if no body
	if funcDecl.Body == nil {
		// External function has no body
		return true
	}

	// Skip test functions
	if shared.IsTestFunction(funcDecl) {
		// Test function should be skipped
		return true
	}

	// Function should be analyzed
	return false
}

// checkFunctionBody checks all statements in a function body.
//
// Params:
//   - pass: analysis pass
//   - body: function body to check
func checkFunctionBody(pass *analysis.Pass, body *ast.BlockStmt) {
	// Check all statements
	ast.Inspect(body, func(node ast.Node) bool {
		// Check different statement types
		switch stmt := node.(type) {
		// Check if statement
		case *ast.IfStmt:
			checkIfStmt(pass, stmt)
		// Check switch statement
		case *ast.SwitchStmt:
			checkSwitchStmt(pass, stmt)
		// Check type switch statement
		case *ast.TypeSwitchStmt:
			checkTypeSwitchStmt(pass, stmt)
		// Check loop statements
		case *ast.ForStmt, *ast.RangeStmt:
			checkLoopStmt(pass, stmt)
		// Check return statement
		case *ast.ReturnStmt:
			checkReturnStmt(pass, stmt)
		}
		// Continue traversal
		return true
	})
}

// checkIfStmt checks if an if statement has a comment
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: if statement to check
func checkIfStmt(pass *analysis.Pass, stmt *ast.IfStmt) {
	// Vérification que le if a un commentaire (règle stricte, pas d'exception)
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-COMMENT-007: le bloc 'if' doit avoir un commentaire explicatif")
	}

	// Check else clause
	if stmt.Else != nil {
		// For else, check if there's a comment before or at the start of the else block
		if !hasCommentBeforeOrInside(pass, stmt.Else) {
			pass.Reportf(stmt.Else.Pos(), "KTN-COMMENT-007: le bloc 'else' doit avoir un commentaire explicatif")
		}
	}
}

// checkSwitchStmt checks if a switch statement and its cases have comments
// Params:
//   - pass: contexte d'analyse
func checkSwitchStmt(pass *analysis.Pass, stmt *ast.SwitchStmt) {
	// Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-COMMENT-007: le bloc 'switch' doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
		// Itération sur les éléments
		for _, caseClause := range stmt.Body.List {
			// Vérification de la condition
			if clause, ok := caseClause.(*ast.CaseClause); ok {
				// Vérification de la condition
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-COMMENT-007: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkTypeSwitchStmt checks if a type switch statement and its cases have comments
// Params:
//   - pass: contexte d'analyse
func checkTypeSwitchStmt(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) {
	// Vérification de la condition
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-COMMENT-007: le bloc 'switch' (type) doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
		// Itération sur les éléments
		for _, caseClause := range stmt.Body.List {
			// Vérification de la condition
			if clause, ok := caseClause.(*ast.CaseClause); ok {
				// Vérification de la condition
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-COMMENT-007: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkLoopStmt checks if a loop statement has a comment
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: loop statement to check
func checkLoopStmt(pass *analysis.Pass, stmt ast.Node) {
	// Vérification que la boucle a un commentaire
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-COMMENT-007: le bloc de boucle doit avoir un commentaire explicatif")
	}
}

// checkReturnStmt checks if a return statement has a comment
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: return statement to check
func checkReturnStmt(pass *analysis.Pass, stmt *ast.ReturnStmt) {
	// Vérification que le return a un commentaire (règle stricte, pas d'exception)
	if !hasCommentBefore(pass, stmt.Pos()) && !hasInlineComment(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-COMMENT-007: le 'return' doit avoir un commentaire explicatif")
	}
}

// hasCommentBefore checks if there's a comment on the line before the given position
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si commentaire présent
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
