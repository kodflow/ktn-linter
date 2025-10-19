package ktnfunc

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer007 checks that all functions have proper documentation
var Analyzer007 = &analysis.Analyzer{
	Name:     "ktnfunc007",
	Doc:      "KTN-FUNC-007: Toutes les fonctions doivent avoir une documentation au format strict (description, Params, Returns)",
	Run:      runFunc007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	// Pattern for function description (first line)
	descPattern = regexp.MustCompile(`^//\s*\w+\s+.+`)

	// Pattern for Params section
	paramsHeaderPattern = regexp.MustCompile(`^//\s*Params:\s*$`)
	paramItemPattern    = regexp.MustCompile(`^//\s*-\s*\w+:\s*.+`)

	// Pattern for Returns section
	returnsHeaderPattern = regexp.MustCompile(`^//\s*Returns:\s*$`)
	returnItemPattern    = regexp.MustCompile(`^//\s*-\s*.+:\s*.+`)

	// Pattern for Example section (optional)
	exampleHeaderPattern = regexp.MustCompile(`^//\s*Example:\s*$`)
)

func runFunc007(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Get function documentation
		doc := funcDecl.Doc

		// Check if documentation exists
		if doc == nil || len(doc.List) == 0 {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-007: la fonction '%s' doit avoir une documentation",
				funcDecl.Name.Name,
			)
			return
		}

		// Extract comment lines
		comments := extractCommentLines(doc)

		// Validate documentation format
		hasParams := funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0
		hasReturns := funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0

		err := validateDocFormat(comments, funcDecl.Name.Name, hasParams, hasReturns)
		if err != "" {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-007: %s",
				err,
			)
		}
	})

	return nil, nil
}

// extractCommentLines extracts individual comment lines from a CommentGroup
func extractCommentLines(cg *ast.CommentGroup) []string {
	var lines []string
	for _, comment := range cg.List {
		text := comment.Text
		// Remove // prefix
		if strings.HasPrefix(text, "//") {
			lines = append(lines, text)
		}
	}
	return lines
}

// validateDocFormat validates the documentation format
func validateDocFormat(comments []string, funcName string, hasParams, hasReturns bool) string {
	if len(comments) == 0 {
		return "documentation vide"
	}

	// Check first line (description)
	if !strings.HasPrefix(comments[0], "// "+funcName) {
		return "la description doit commencer par '// " + funcName + " '"
	}

	idx := 1

	// Skip blank line if present
	if idx < len(comments) && strings.TrimSpace(comments[idx]) == "//" {
		idx++
	}

	// Check for Params section if function has parameters
	if hasParams {
		if idx >= len(comments) || !paramsHeaderPattern.MatchString(comments[idx]) {
			return "section 'Params:' manquante (fonction avec paramètres)"
		}
		idx++

		// Validate param items
		foundParam := false
		for idx < len(comments) {
			line := comments[idx]
			if paramItemPattern.MatchString(line) {
				foundParam = true
				idx++
			} else {
				break
			}
		}

		if !foundParam {
			return "au moins un paramètre doit être documenté dans 'Params:'"
		}

		// Skip blank line
		if idx < len(comments) && strings.TrimSpace(comments[idx]) == "//" {
			idx++
		}
	}

	// Check for Returns section if function has return values
	if hasReturns {
		if idx >= len(comments) || !returnsHeaderPattern.MatchString(comments[idx]) {
			return "section 'Returns:' manquante (fonction avec valeurs de retour)"
		}
		idx++

		// Validate return items
		foundReturn := false
		for idx < len(comments) {
			line := comments[idx]
			if returnItemPattern.MatchString(line) {
				foundReturn = true
				idx++
			} else {
				break
			}
		}

		if !foundReturn {
			return "au moins une valeur de retour doit être documentée dans 'Returns:'"
		}
	}

	return ""
}
