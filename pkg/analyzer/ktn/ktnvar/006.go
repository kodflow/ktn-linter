// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

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
	// ruleCodeVar006 is the rule code for this analyzer
	ruleCodeVar006 string = "KTN-VAR-006"
)

// builtinIdentifiers006 contains all Go built-in identifiers (44 total).
// Types (22): bool, byte, complex64, complex128, error, float32, float64,
//
//	int, int8, int16, int32, int64, rune, string, uint, uint8,
//	uint16, uint32, uint64, uintptr, any, comparable
//
// Constants (3): true, false, iota
// Zero-value (1): nil
// Functions (18): append, cap, clear, close, complex, copy, delete, imag,
//
//	len, make, max, min, new, panic, print, println, real, recover
var builtinIdentifiers006 map[string]bool = map[string]bool{
	// Types (22)
	"bool":       true,
	"byte":       true,
	"complex64":  true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"any":        true,
	"comparable": true,
	// Constants (3)
	"true":  true,
	"false": true,
	"iota":  true,
	// Zero-value (1)
	"nil": true,
	// Functions (18)
	"append":  true,
	"cap":     true,
	"clear":   true,
	"close":   true,
	"complex": true,
	"copy":    true,
	"delete":  true,
	"imag":    true,
	"len":     true,
	"make":    true,
	"max":     true,
	"min":     true,
	"new":     true,
	"panic":   true,
	"print":   true,
	"println": true,
	"real":    true,
	"recover": true,
}

// Analyzer006 checks that variables do not shadow built-in identifiers.
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar006",
	Doc:      "KTN-VAR-006: Verifie que les variables ne masquent pas les identifiants built-in",
	Run:      runVar006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar006 executes KTN-VAR-006 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runVar006(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar006) {
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
		if cfg.IsFileExcluded(ruleCodeVar006, filename) {
			// File is excluded
			return
		}
		genDecl := n.(*ast.GenDecl)

		// Only check var declarations
		if genDecl.Tok != token.VAR {
			// Return early
			return
		}

		// Iterate over specs
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Iterate over names
			for _, name := range valueSpec.Names {
				varName := name.Name

				// Check if name shadows built-in
				if isBuiltinIdentifier006(varName) {
					msg, _ := messages.Get(ruleCodeVar006)
					pass.Reportf(
						name.Pos(),
						"%s: %s",
						ruleCodeVar006,
						msg.Format(cfg.Verbose, varName),
					)
				}
			}
		}
	})

	// Return result
	return nil, nil
}

// isBuiltinIdentifier006 checks if a name is a Go built-in identifier.
// Blank identifier (_) is always allowed.
//
// Params:
//   - name: variable name to check
//
// Returns:
//   - bool: true if the name shadows a built-in
func isBuiltinIdentifier006(name string) bool {
	// Blank identifier is always allowed
	if name == "_" {
		// Skip blank identifier
		return false
	}

	// Check if name is in built-in map
	return builtinIdentifiers006[name]
}
