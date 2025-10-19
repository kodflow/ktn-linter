package ktnfunc

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer011 checks that all branches, returns, and significant logic blocks have comments
var Analyzer011 = &analysis.Analyzer{
	Name:     "ktnfunc011",
	Doc:      "KTN-FUNC-011: Tous les blocs de contrôle (if/else/switch/for), returns et logique significative doivent être commentés",
	Run:      runFunc011,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runFunc011(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body
		if funcDecl.Body == nil {
			return
		}

		// Skip test functions (Test*, Benchmark*, Example*, Fuzz*)
		funcName := funcDecl.Name.Name
		if isTestFunction(funcName) {
			return
		}

		// Check all statements in the function body
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			switch stmt := node.(type) {
			case *ast.IfStmt:
				checkIfStmt(pass, stmt)
			case *ast.SwitchStmt:
				checkSwitchStmt(pass, stmt)
			case *ast.TypeSwitchStmt:
				checkTypeSwitchStmt(pass, stmt)
			case *ast.ForStmt, *ast.RangeStmt:
				checkLoopStmt(pass, stmt)
			case *ast.ReturnStmt:
				checkReturnStmt(pass, stmt)
			}
			return true
		})
	})

	return nil, nil
}

// checkIfStmt checks if an if statement has a comment
func checkIfStmt(pass *analysis.Pass, stmt *ast.IfStmt) {
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
func checkSwitchStmt(pass *analysis.Pass, stmt *ast.SwitchStmt) {
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc 'switch' doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
		for _, caseClause := range stmt.Body.List {
			if clause, ok := caseClause.(*ast.CaseClause); ok {
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-FUNC-011: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkTypeSwitchStmt checks if a type switch statement and its cases have comments
func checkTypeSwitchStmt(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) {
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc 'switch' (type) doit avoir un commentaire explicatif")
	}

	// Check each case
	if stmt.Body != nil {
		for _, caseClause := range stmt.Body.List {
			if clause, ok := caseClause.(*ast.CaseClause); ok {
				if !hasCommentBefore(pass, clause.Pos()) {
					pass.Reportf(clause.Pos(), "KTN-FUNC-011: chaque 'case' doit avoir un commentaire explicatif")
				}
			}
		}
	}
}

// checkLoopStmt checks if a loop statement has a comment
func checkLoopStmt(pass *analysis.Pass, stmt ast.Node) {
	if !hasCommentBefore(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le bloc de boucle doit avoir un commentaire explicatif")
	}
}

// checkReturnStmt checks if a return statement has a comment
func checkReturnStmt(pass *analysis.Pass, stmt *ast.ReturnStmt) {
	if !hasCommentBefore(pass, stmt.Pos()) && !hasInlineComment(pass, stmt.Pos()) {
		pass.Reportf(stmt.Pos(), "KTN-FUNC-011: le 'return' doit avoir un commentaire explicatif")
	}
}

// hasCommentBefore checks if there's a comment on the line before the given position
func hasCommentBefore(pass *analysis.Pass, pos token.Pos) bool {
	position := pass.Fset.Position(pos)
	filename := position.Filename

	// Find the file in the pass
	var file *ast.File
	for _, f := range pass.Files {
		if pass.Fset.Position(f.Pos()).Filename == filename {
			file = f
			break
		}
	}

	if file == nil {
		return false
	}

	// Check all comments in the file
	for _, commentGroup := range file.Comments {
		commentPos := pass.Fset.Position(commentGroup.End())
		// Comment should be on the line immediately before the statement
		if commentPos.Filename == filename && commentPos.Line == position.Line-1 {
			// Check if it's a real comment (not a "want" directive)
			for _, comment := range commentGroup.List {
				text := strings.TrimSpace(comment.Text)
				if !strings.Contains(text, "want") {
					return true
				}
			}
		}
	}

	return false
}

// hasInlineComment checks if there's a comment on the same line as the given position
func hasInlineComment(pass *analysis.Pass, pos token.Pos) bool {
	position := pass.Fset.Position(pos)
	filename := position.Filename

	// Find the file in the pass
	var file *ast.File
	for _, f := range pass.Files {
		if pass.Fset.Position(f.Pos()).Filename == filename {
			file = f
			break
		}
	}

	if file == nil {
		return false
	}

	// Check all comments in the file
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			commentPos := pass.Fset.Position(comment.Pos())
			// Comment should be on the same line as the statement
			if commentPos.Filename == filename && commentPos.Line == position.Line {
				// Check if it's a real comment (not a "want" directive)
				text := strings.TrimSpace(comment.Text)
				if !strings.Contains(text, "want") {
					return true
				}
			}
		}
	}

	return false
}

// hasCommentBeforeOrInside checks if there's a comment before or at the start of an else block
func hasCommentBeforeOrInside(pass *analysis.Pass, elseStmt ast.Stmt) bool {
	// Check if there's a comment before the else
	if hasCommentBefore(pass, elseStmt.Pos()) {
		return true
	}

	// If it's a block statement, check for a comment at the beginning
	if blockStmt, ok := elseStmt.(*ast.BlockStmt); ok && len(blockStmt.List) > 0 {
		firstStmt := blockStmt.List[0]
		if hasCommentBefore(pass, firstStmt.Pos()) {
			return true
		}
	}

	// If it's another if statement (else if), check for a comment before it
	if ifStmt, ok := elseStmt.(*ast.IfStmt); ok {
		if hasCommentBefore(pass, ifStmt.Pos()) {
			return true
		}
	}

	return false
}
