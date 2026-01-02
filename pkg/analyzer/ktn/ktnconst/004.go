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
	// ruleCodeConst004 est le code de la regle KTN-CONST-004.
	ruleCodeConst004 string = "KTN-CONST-004"
	// minConstNameLen is the minimum length for constant names.
	minConstNameLen int = 2
)

// Analyzer004 checks that constants have names with at least 2 characters.
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnconst004",
	Doc:      "KTN-CONST-004: Verifie que les noms de constantes ont au moins 2 caracteres",
	Run:      runConst004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runConst004 executes KTN-CONST-004 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runConst004(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeConst004) {
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
		if cfg.IsFileExcluded(ruleCodeConst004, filename) {
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

				// Check if name is too short
				if isConstNameTooShort(constName) {
					msg, _ := messages.Get(ruleCodeConst004)
					pass.Reportf(
						name.Pos(),
						"%s: %s",
						ruleCodeConst004,
						msg.Format(cfg.Verbose, constName, minConstNameLen),
					)
				}
			}
		}
	})

	// Return result
	return nil, nil
}

// isConstNameTooShort checks if a constant name is too short.
// Blank identifier (_) is always allowed.
//
// Params:
//   - name: constant name to check
//
// Returns:
//   - bool: true if the name is too short
func isConstNameTooShort(name string) bool {
	// Blank identifier is always allowed
	if name == "_" {
		// Skip blank identifier
		return false
	}

	// Check minimum length
	return len(name) < minConstNameLen
}
