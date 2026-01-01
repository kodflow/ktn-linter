// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar024 is the rule code for this analyzer
	ruleCodeVar024 string = "KTN-VAR-024"
)

// Analyzer024 detects usage of interface{} and suggests using any.
//
// Since Go 1.18, any is an alias for interface{}. Using any is more
// readable and idiomatic in modern Go code.
var Analyzer024 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar024",
	Doc:      "KTN-VAR-024: Preferer any a interface{}",
	Run:      runVar024,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar024 executes the analysis for KTN-VAR-024.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: always nil
//   - error: any error encountered
func runVar024(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar024) {
		// Regle desactivee
		return nil, nil
	}

	// Recuperation de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Types de noeuds a analyser
	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	// Parcours des noeuds
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar024, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Verification du type de noeud
		interfaceType, ok := n.(*ast.InterfaceType)
		// Verification de la condition
		if !ok {
			// Pas un interface type
			return
		}

		// Verification si c'est une interface vide
		checkEmptyInterface(pass, interfaceType)
	})

	// Traitement
	return nil, nil
}

// checkEmptyInterface checks if an interface type is empty (interface{}).
//
// Params:
//   - pass: analysis context
//   - interfaceType: interface type to check
func checkEmptyInterface(pass *analysis.Pass, interfaceType *ast.InterfaceType) {
	// Verifier si c'est une interface vide
	if !isEmptyInterface(interfaceType) {
		// Interface non vide
		return
	}

	// Report de l'erreur
	msg, _ := messages.Get(ruleCodeVar024)
	pass.Reportf(
		interfaceType.Pos(),
		"%s: %s",
		ruleCodeVar024,
		msg.Format(config.Get().Verbose),
	)
}

// isEmptyInterface checks if an interface type is empty.
//
// Params:
//   - interfaceType: interface type to check
//
// Returns:
//   - bool: true if it's an empty interface
func isEmptyInterface(interfaceType *ast.InterfaceType) bool {
	// Interface vide si Methods est nil ou vide
	if interfaceType.Methods == nil {
		// Methods nil
		return true
	}

	// Verifier si la liste est vide
	return len(interfaceType.Methods.List) == 0
}
