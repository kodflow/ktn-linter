// Package ktngeneric implements KTN linter rules for generic functions.
package ktngeneric

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeGeneric003 est le code de la regle KTN-GENERIC-003.
	ruleCodeGeneric003 string = "KTN-GENERIC-003"
	// deprecatedConstraintsPath est le chemin du package obsolete.
	deprecatedConstraintsPath string = "golang.org/x/exp/constraints"
)

// Analyzer003 checks that deprecated x/exp/constraints package is not used.
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktngeneric003",
	Doc:      "KTN-GENERIC-003: Deprecated golang.org/x/exp/constraints package should be replaced by cmp",
	Run:      runGeneric003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runGeneric003 execute l'analyse KTN-GENERIC-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runGeneric003(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeGeneric003) {
		// Regle desactivee
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeGeneric003, filename) {
			// File is excluded
			return
		}
		importSpec := n.(*ast.ImportSpec)

		// Analyser l'import
		checkDeprecatedConstraintsImport(pass, importSpec)
	})

	// Retour de la fonction
	return nil, nil
}

// checkDeprecatedConstraintsImport verifie si l'import est le package obsolete.
//
// Params:
//   - pass: contexte d'analyse
//   - importSpec: specification d'import a analyser
func checkDeprecatedConstraintsImport(pass *analysis.Pass, importSpec *ast.ImportSpec) {
	// Verifier si l'import est nil
	if importSpec.Path == nil {
		// Import invalide
		return
	}

	// Extraire le chemin de l'import (enlever les guillemets)
	importPath := importSpec.Path.Value
	// Verifier si c'est le package obsolete
	if importPath != `"`+deprecatedConstraintsPath+`"` {
		// Pas le package obsolete
		return
	}

	// Guard contre nil (pour tests unitaires)
	if pass == nil {
		// Pas de contexte pour reporter
		return
	}

	// Reporter l'erreur
	cfg := config.Get()
	msg, _ := messages.Get(ruleCodeGeneric003)
	pass.Reportf(
		importSpec.Pos(),
		"%s: %s",
		ruleCodeGeneric003,
		msg.Format(cfg.Verbose),
	)
}
