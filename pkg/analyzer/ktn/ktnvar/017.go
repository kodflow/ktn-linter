// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar017 is the rule code for this analyzer
	ruleCodeVar017 string = "KTN-VAR-017"
)

// Analyzer017 detects map allocations without capacity hints
var Analyzer017 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar017",
	Doc:      "KTN-VAR-017: Préallouer maps avec capacité si connue",
	Run:      runVar017,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar017 exécute l'analyse KTN-VAR-017.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar017(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar017) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr, ok := n.(*ast.CallExpr)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar017, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification que c'est un appel à "make"
		if !utils.IsMakeCall(callExpr) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que le type est une map
		if len(callExpr.Args) == 0 || !utils.IsMapTypeWithPass(pass, callExpr.Args[0]) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que make a exactement 1 argument (type seulement)
		if len(callExpr.Args) != 1 {
			// make() avec capacité fournie, conforme
			return
		}

		// Signaler l'erreur
		msg, ok := messages.Get(ruleCodeVar017)
		// Defensive: avoid panic if message is missing
		if !ok {
			pass.Reportf(
				callExpr.Pos(),
				"%s: préallouer map avec une capacité quand elle est connue",
				ruleCodeVar017,
			)
			return
		}
		pass.Reportf(
			callExpr.Pos(),
			"%s: %s",
			ruleCodeVar017,
			msg.Format(config.Get().Verbose),
		)
	})

	// Retour de la fonction
	return nil, nil
}
