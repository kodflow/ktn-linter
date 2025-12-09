// Analyzer 001 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/token"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// maxCommentLength max chars for inline comments
	maxCommentLength int = 80
)

var (
	// urlPattern matches URLs in comments (http://, https://, file://)
	urlPattern *regexp.Regexp = regexp.MustCompile(`https?://\S+|file://\S+`)

	// Analyzer001 detects inline comments exceeding 80 characters.
	Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktncomment001",
		Doc:      "KTN-COMMENT-001: commentaire inline trop long (>80 chars)",
		Run:      runComment001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

// runComment001 analyzes inline comments for excessive length.
//
// Params:
//   - pass: Analysis pass
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runComment001(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Parcours des groupes de commentaires
		for _, commentGroup := range file.Comments {
			// Parcours des commentaires individuels
			for _, comment := range commentGroup.List {
				// Skip doc comments (start at beginning of line)
				if isDocComment(pass, file, comment) {
					// Continue au commentaire suivant
					continue
				}

				// Check inline comment length
				text := comment.Text

				// Handle multi-line block comments /* ... */
				if strings.HasPrefix(text, "/*") {
					// Check each line separately for multi-line comments
					checkMultiLineComment(pass, comment, text)
					// Continue au commentaire suivant
					continue
				}

				// Remove // prefix for single-line comments
				text = strings.TrimPrefix(text, "//")

				// Skip comments containing URLs
				if containsURL(text) {
					// Continue au commentaire suivant
					continue
				}

				// Check length exceeds limit
				if len(text) > maxCommentLength {
					pass.Reportf(
						comment.Pos(),
						"KTN-COMMENT-001: commentaire inline trop long (>%d chars)",
						maxCommentLength,
					)
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkMultiLineComment checks each line of a block comment separately.
//
// Params:
//   - pass: Analysis pass
//   - comment: The comment node
//   - text: The comment text
func checkMultiLineComment(pass *analysis.Pass, comment *ast.Comment, text string) {
	// Remove /* and */ markers
	content := strings.TrimPrefix(text, "/*")
	content = strings.TrimSuffix(content, "*/")

	// Check each line using SplitSeq
	for line := range strings.SplitSeq(content, "\n") {
		// Trim leading/trailing whitespace for length check
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			// Continue à la ligne suivante
			continue
		}

		// Skip lines containing URLs
		if containsURL(trimmed) {
			// Continue à la ligne suivante
			continue
		}

		// Check length exceeds limit
		if len(trimmed) > maxCommentLength {
			pass.Reportf(
				comment.Pos(),
				"KTN-COMMENT-001: ligne de commentaire trop longue (>%d chars)",
				maxCommentLength,
			)
			// Only report once per block comment
			return
		}
	}
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
			// Commentaire de doc détecté
			return true
		}

		// Check function declarations
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Check if comment is doc comment
			if funcDecl.Doc != nil && slices.Contains(funcDecl.Doc.List, comment) {
				// Commentaire de doc fonction
				return true
			}
		}

		// Check general declarations
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Check if comment is doc comment
			if genDecl.Doc != nil && slices.Contains(genDecl.Doc.List, comment) {
				// Commentaire de doc général
				return true
			}
		}
	}

	// Vérification position début de ligne
	return isCommentAtLineStart(pass, comment)
}

// isCommentAtLineStart checks if comment starts at line beginning.
//
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

// containsURL checks if text contains a URL.
//
// Params:
//   - text: text to check
//
// Returns:
//   - bool: true if text contains a URL
func containsURL(text string) bool {
	// Check if URL pattern matches
	return urlPattern.MatchString(text)
}
