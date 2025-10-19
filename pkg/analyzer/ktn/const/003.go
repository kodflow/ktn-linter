package ktnconst

import (
	"go/ast"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks that constants use CAPITAL_UNDERSCORE naming convention
var Analyzer003 = &analysis.Analyzer{
	Name:     "ktnconst003",
	Doc:      "KTN-CONST-003: Vérifie que les constantes utilisent la convention CAPITAL_UNDERSCORE",
	Run:      runConst003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// validConstNamePattern matches valid CAPITAL_UNDERSCORE constant names
// Must start with uppercase letter, followed by uppercase letters, digits, or underscores
// Must contain at least one underscore for multi-word constants
var validConstNamePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

func runConst003(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			return
		}

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			for _, name := range valueSpec.Names {
				constName := name.Name

				// Skip blank identifiers
				if constName == "_" {
					continue
				}

				// Check if the constant name follows CAPITAL_UNDERSCORE convention
				if !isValidConstantName(constName) {
					pass.Reportf(
						name.Pos(),
						"KTN-CONST-003: la constante '%s' doit utiliser la convention CAPITAL_UNDERSCORE (ex: MAX_SIZE, API_KEY, HTTP_TIMEOUT)",
						constName,
					)
				}
			}
		}
	})

	return nil, nil
}

// isValidConstantName checks if a constant name follows CAPITAL_UNDERSCORE convention
func isValidConstantName(name string) bool {
	// Must match the pattern: starts with uppercase, contains only uppercase, digits, underscores
	if !validConstNamePattern.MatchString(name) {
		return false
	}

	// Single letter constants are valid (e.g., A, B, C)
	if len(name) == 1 {
		return true
	}

	// Pattern already validates SCREAMING_SNAKE_CASE:
	// - Single letters: A, B, C
	// - Acronyms: API, HTTP, URL, HTTPS, EOF
	// - Underscored names: MAX_SIZE, API_KEY, HTTP_TIMEOUT
	// - With numbers: HTTP2, TLS1_2_VERSION
	return true
}
