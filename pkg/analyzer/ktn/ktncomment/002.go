// Analyzer 002 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MAX_COMMENT_LENGTH max chars for inline comments
	MAX_COMMENT_LENGTH int = 80
)

// Analyzer002 detects inline comments exceeding 80 characters.
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment002",
	Doc:      "KTN-COMMENT-002: commentaire inline trop long (>80 chars)",
	Run:      runComment002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment002 analyzes inline comments for excessive length.
// Params:
//   - pass: Analysis pass
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runComment002(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Verification de la condition
		for _, commentGroup := range file.Comments {
			// Verification de la condition
			for _, comment := range commentGroup.List {
				// Skip doc comments (start at beginning of line)
				if isDocComment(pass, file, comment) {
					continue
				}

				// Check inline comment length
				text := comment.Text
				// Remove comment markers
				text = strings.TrimPrefix(text, "//")
				text = strings.TrimPrefix(text, "/*")
				text = strings.TrimSuffix(text, "*/")

				// Check length exceeds limit
				if len(text) > MAX_COMMENT_LENGTH {
					pass.Reportf(
						comment.Pos(),
						"KTN-COMMENT-002: commentaire inline trop long (>%d chars)",
						MAX_COMMENT_LENGTH,
					)
				}
			}
		}
	})

	return nil, nil
}

// isDocComment checks if comment is documentation comment.
//
// Params:
//   - pass: Analysis pass
//   - file: File containing comment
//   - comment: Comment to check
//
// Returns:
//   - bool: true if comment is documentation comment
func isDocComment(
	pass *analysis.Pass,
	file *ast.File,
	comment *ast.Comment,
) bool {
	// Get comment line using helper function
	line := getCommentLine(pass.Fset, comment)

	// Check if comment is at start of line
	for _, decl := range file.Decls {
		declLine := pass.Fset.Position(decl.Pos()).Line
		// Check if comment precedes declaration
		if declLine == line+1 || declLine == line {
			return true
		}

		// Check function declarations
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Check if comment is doc comment
			if funcDecl.Doc != nil && slices.Contains(funcDecl.Doc.List, comment) {
				return true
			}
		}

		// Check general declarations
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Check if comment is doc comment
			if genDecl.Doc != nil && slices.Contains(genDecl.Doc.List, comment) {
				return true
			}
		}
	}

	// Verification de la condition
	return isCommentAtLineStart(pass, comment)
}

// isCommentAtLineStart checks if comment starts at line beginning.
// Params:
//   - pass: Analysis pass
//   - comment: Comment to check
//
// Returns:
//   - bool: true if at line start
func isCommentAtLineStart(pass *analysis.Pass, comment *ast.Comment) bool {
	pos := pass.Fset.Position(comment.Pos())
	// Check if comment is at column 1
	return pos.Column == 1
}

// getCommentLine returns line number of a comment.
//
// Params:
//   - fset: File set
//   - comment: Comment to get line for
//
// Returns:
//   - int: line number of the comment
func getCommentLine(fset *token.FileSet, comment *ast.Comment) int {
	// Return the line number
	return fset.Position(comment.Pos()).Line
}

