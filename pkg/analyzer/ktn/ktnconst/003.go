// Package ktnconst implements KTN linter rules.
package ktnconst

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks that constants use standard Go naming conventions (CamelCase).
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnconst003",
	Doc:      "KTN-CONST-003: Vérifie que les constantes utilisent les conventions Go standard (CamelCase)",
	Run:      runConst003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runConst003 executes KTN-CONST-003 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runConst003(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			// Return early
			return
		}

		// Iterate over specs
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Iterate over names
			for _, name := range valueSpec.Names {
				constName := name.Name

				// Skip blank identifiers
				if constName == "_" {
					continue
				}

				// Check if the constant name follows Go conventions
				if !isValidGoConstantName(constName) {
					pass.Reportf(
						name.Pos(),
						"KTN-CONST-003: la constante '%s' doit utiliser la convention CamelCase (ex: MaxSize, httpTimeout)",
						constName,
					)
				}
			}
		}
	})

	// Return result
	return nil, nil
}

// isValidGoConstantName checks if a constant name follows Go CamelCase convention.
//
// Params:
//   - name: constant name to check
//
// Returns:
//   - bool: true if the name follows Go conventions
//
// Valid names: PascalCase for exported (MaxSize, HttpTimeout, APIKey) or camelCase for unexported (maxSize, httpTimeout, apiKey).
// Invalid names: SCREAMING_SNAKE_CASE (MAX_SIZE), snake_case (max_size), Contains underscores (Max_Size).
func isValidGoConstantName(name string) bool {
	// Empty name is invalid
	if len(name) == 0 {
		// Return early
		return false
	}

	// Single character is valid if letter
	if len(name) == 1 {
		// Return early
		return unicode.IsLetter(rune(name[0]))
	}

	// Check for underscores (not allowed in Go CamelCase)
	if strings.Contains(name, "_") {
		// Return early
		return false
	}

	// Check first character is a letter
	firstRune := rune(name[0])
	// Vérification que le premier caractère est une lettre
	if !unicode.IsLetter(firstRune) {
		// Return early
		return false
	}

	// Check all characters are valid (letters or digits)
	for _, r := range name {
		// Vérification que chaque caractère est une lettre ou un chiffre
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			// Return early
			return false
		}
	}

	// Name is valid CamelCase
	return true
}
