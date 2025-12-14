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
	// ruleCodeConst001 est le code de la règle KTN-CONST-001.
	ruleCodeConst001 string = "KTN-CONST-001"
)

// Analyzer001 checks that constants have explicit types
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnconst001",
	Doc:      "KTN-CONST-001: Vérifie que les constantes ont un type explicite",
	Run:      runConst001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runConst001 exécute l'analyse KTN-CONST-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runConst001(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeConst001) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeConst001, filename) {
			// File is excluded
			return
		}
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			// Retour de la fonction
			return
		}

		// Itération sur les éléments
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Check if the constant has an explicit type
			if valueSpec.Type == nil {
				// Error if value present (not iota inheritance)
				// OK if no value (inherits from previous)
				if len(valueSpec.Values) > 0 {
					// Itération sur les éléments
					for _, name := range valueSpec.Names {
						msg, _ := messages.Get(ruleCodeConst001)
						pass.Reportf(
							name.Pos(),
							"%s: %s",
							ruleCodeConst001,
							msg.Format(cfg.Verbose, name.Name),
						)
					}
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
