package ktnfunc

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 checks that functions don't exceed 35 lines of pure code
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktnfunc001",
	Doc:      "KTN-FUNC-001: Les fonctions ne doivent pas dépasser 35 lignes de code pur (hors commentaires et lignes vides)",
	Run:      runFunc001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const maxPureCodeLines = 35

func runFunc001(pass *analysis.Pass) (any, error) {
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

		// Skip main function
		if funcName == "main" {
			return
		}

		// Count pure code lines
		pureLines := countPureCodeLines(pass, funcDecl.Body)

		if pureLines > maxPureCodeLines {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-001: la fonction '%s' contient %d lignes de code pur (max: %d)",
				funcName,
				pureLines,
				maxPureCodeLines,
			)
		}
	})

	return nil, nil
}

// isTestFunction checks if a function name indicates a test function
func isTestFunction(name string) bool {
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Benchmark") ||
		strings.HasPrefix(name, "Example") ||
		strings.HasPrefix(name, "Fuzz")
}

// countPureCodeLines counts only pure code lines by reading the source file
func countPureCodeLines(pass *analysis.Pass, body *ast.BlockStmt) int {
	if body == nil {
		return 0
	}

	startPos := pass.Fset.Position(body.Lbrace)
	endPos := pass.Fset.Position(body.Rbrace)

	// Read the source file
	filename := startPos.Filename
	if pass.ReadFile == nil {
		return 0
	}
	content, err := pass.ReadFile(filename)
	if err != nil {
		return 0
	}

	lines := strings.Split(string(content), "\n")
	pureCodeLines := 0
	inBlockComment := false

	// Count lines between start and end (excluding braces)
	// Start after the opening brace line
	for i := startPos.Line + 1; i < endPos.Line; i++ {
		if i <= 0 || i > len(lines) {
			continue
		}

		line := lines[i-1] // lines are 0-indexed
		trimmed := strings.TrimSpace(line)

		// Check for block comment start/end
		if strings.Contains(trimmed, "/*") {
			inBlockComment = true
		}
		if strings.Contains(trimmed, "*/") {
			inBlockComment = false
			continue
		}

		// Skip if in block comment
		if inBlockComment {
			continue
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip line comments
		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Skip standalone braces
		if trimmed == "{" || trimmed == "}" {
			continue
		}

		// This is a pure code line
		pureCodeLines++
	}

	return pureCodeLines
}
