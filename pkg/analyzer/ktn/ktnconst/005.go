// Package ktnconst implements KTN linter rules.
package ktnconst

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeConst005 est le code de la regle KTN-CONST-005.
	ruleCodeConst005 string = "KTN-CONST-005"
	// maxConstNameLen is the maximum length for constant names.
	maxConstNameLen int = 30
)

// Analyzer005 checks that constants have names with at most 30 characters.
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnconst005",
	Doc:      "KTN-CONST-005: Verifie que les noms de constantes ont au maximum 30 caracteres",
	Run:      runConst005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runConst005 executes KTN-CONST-005 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runConst005(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeConst005) {
		// Regle desactivee
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeConst005, filename) {
			// File is excluded
			return
		}
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

				// Check if name is too long
				if isConstNameTooLong(constName) {
					msg, _ := messages.Get(ruleCodeConst005)
					pass.Reportf(
						name.Pos(),
						"%s: %s",
						ruleCodeConst005,
						msg.Format(cfg.Verbose, constName, len(constName), maxConstNameLen),
					)
				}
			}
		}
	})

	// Return result
	return nil, nil
}

// isConstNameTooLong checks if a constant name is too long.
// Blank identifier (_) is always allowed.
//
// Params:
//   - name: constant name to check
//
// Returns:
//   - bool: true if the name is too long
func isConstNameTooLong(name string) bool {
	// Blank identifier is always allowed
	if name == "_" {
		// Skip blank identifier
		return false
	}

	// Check maximum length
	return len(name) > maxConstNameLen
}
