package ktncomment

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 detects redundant comments that state the obvious.
// Params:
//   - N/A
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktncomment001",
	Doc:      "KTN-COMMENT-001: commentaire redondant",
	Run:      runComment001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// redundantPatterns defines common redundant comment patterns.
var redundantPatterns = []struct {
	comment string
	code    string
}{
	{"return nil", "return nil"},
	{"return", "return"},
	{"set", ":="},
	{"increment", "++"},
	{"decrement", "--"},
	{"loop", "for"},
	{"iterate", "for"},
	{"create", ":="},
	{"initialize", ":="},
}

// runComment001 analyzes comments for redundancy.
// Params:
//   - pass: Analysis pass
func runComment001(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Build map of comments to their following statements
		commentStmtMap := buildCommentStmtMap(pass, file)

		for comment, stmt := range commentStmtMap {
			// Skip doc comments
			if isDocComment(pass, file, comment) {
				continue
			}

			commentText := normalizeComment(comment.Text)
			stmtText := getStmtText(pass, stmt)

			// Check for redundancy
			if isRedundant(commentText, stmtText) {
				pass.Reportf(
					comment.Pos(),
					"KTN-COMMENT-001: commentaire redondant",
				)
			}
		}
	})

	return nil, nil
}

// buildCommentStmtMap maps comments to following statements.
// Params:
//   - pass: Analysis pass
//   - file: File to analyze
func buildCommentStmtMap(
	pass *analysis.Pass,
	file *ast.File,
) map[*ast.Comment]ast.Stmt {
	result := make(map[*ast.Comment]ast.Stmt)

	ast.Inspect(file, func(node ast.Node) bool {
		// Check for block statements
		blockStmt, ok := node.(*ast.BlockStmt)
		// Continue if not block statement
		if !ok {
			return true
		}

		// Match comments to statements
		for _, stmt := range blockStmt.List {
			stmtPos := pass.Fset.Position(stmt.Pos())

			// Find comments on previous line
			for _, cg := range file.Comments {
				for _, c := range cg.List {
					cPos := pass.Fset.Position(c.Pos())
					// Check if comment is on previous line
					if cPos.Line == stmtPos.Line-1 {
						result[c] = stmt
						break
					}
				}
			}
		}

		return true
	})

	return result
}

// normalizeComment normalizes comment text for comparison.
// Params:
//   - text: Comment text
func normalizeComment(text string) string {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return text
}

// getStmtText extracts statement text for comparison.
// Params:
//   - pass: Analysis pass
//   - stmt: Statement to analyze
func getStmtText(pass *analysis.Pass, stmt ast.Stmt) string {
	// Handle different statement types
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		// Check for nil return
		if len(s.Results) == 1 {
			if ident, ok := s.Results[0].(*ast.Ident); ok {
				if ident.Name == "nil" {
					return "return nil"
				}
			}
		}
		return "return"
	case *ast.AssignStmt:
		// Check for assignment operator
		if s.Tok == token.DEFINE {
			return ":="
		}
		return "="
	case *ast.IncDecStmt:
		// Check for increment/decrement
		if s.Tok == token.INC {
			return "++"
		}
		return "--"
	case *ast.ForStmt, *ast.RangeStmt:
		return "for"
	}
	return ""
}

// isRedundant checks if comment is redundant with code.
// Params:
//   - comment: Normalized comment text
//   - code: Extracted code pattern
func isRedundant(comment, code string) bool {
	// Check for exact matches
	for _, pattern := range redundantPatterns {
		if strings.Contains(comment, pattern.comment) &&
			strings.Contains(code, pattern.code) {
			return true
		}
	}

	// Check for obvious redundancy
	if comment == code {
		return true
	}

	return false
}
