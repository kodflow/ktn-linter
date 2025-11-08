package ktncomment

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const maxCommentLength = 80

// Analyzer002 detects inline comments exceeding 80 characters.
var Analyzer002 = &analysis.Analyzer{
	Name:     "ktncomment002",
	Doc:      "KTN-COMMENT-002: commentaire inline trop long (>80 chars)",
	Run:      runComment002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment002 analyzes inline comments for excessive length.
// Params:
//   - pass: Analysis pass
func runComment002(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		for _, commentGroup := range file.Comments {
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
				if len(text) > maxCommentLength {
					pass.Reportf(
						comment.Pos(),
						"KTN-COMMENT-002: commentaire inline trop long (>%d chars)",
						maxCommentLength,
					)
				}
			}
		}
	})

	return nil, nil
}

// isDocComment checks if comment is documentation comment.
// Params:
//   - pass: Analysis pass
//   - file: File containing comment
//   - comment: Comment to check
func isDocComment(
	pass *analysis.Pass,
	file *ast.File,
	comment *ast.Comment,
) bool {
	pos := pass.Fset.Position(comment.Pos())
	line := pos.Line

	// Check if comment is at start of line
	for _, decl := range file.Decls {
		declPos := pass.Fset.Position(decl.Pos())
		// Check if comment precedes declaration
		if declPos.Line == line+1 || declPos.Line == line {
			return true
		}

		// Check function declarations
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Check if comment is doc comment
			if funcDecl.Doc != nil {
				for _, c := range funcDecl.Doc.List {
					if c == comment {
						return true
					}
				}
			}
		}

		// Check general declarations
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Check if comment is doc comment
			if genDecl.Doc != nil {
				for _, c := range genDecl.Doc.List {
					if c == comment {
						return true
					}
				}
			}
		}
	}

	return isCommentAtLineStart(pass, comment)
}

// isCommentAtLineStart checks if comment starts at line beginning.
// Params:
//   - pass: Analysis pass
//   - comment: Comment to check
func isCommentAtLineStart(pass *analysis.Pass, comment *ast.Comment) bool {
	pos := pass.Fset.Position(comment.Pos())
	// Check if comment is at column 1
	return pos.Column == 1
}

// getCommentLine returns line number of comment.
// Params:
//   - fset: File set
//   - comment: Comment to check
func getCommentLine(fset *token.FileSet, comment *ast.Comment) int {
	return fset.Position(comment.Pos()).Line
}
