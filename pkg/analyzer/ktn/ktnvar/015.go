// Analyzer 015 for the ktnvar package.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar015 is the rule code for this analyzer
	ruleCodeVar015 string = "KTN-VAR-015"
)

// Analyzer015 detects map allocations without capacity hints
var Analyzer015 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar015",
	Doc:      "KTN-VAR-015: Préallouer maps avec capacité si connue",
	Run:      runVar015,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar015 exécute l'analyse KTN-VAR-015.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar015(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar015) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar015, pass.Fset.Position(n.Pos()).Filename) {
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
		pass.Reportf(
			callExpr.Pos(),
			"KTN-VAR-015: préallouer la map avec une capacité (make(map[K]V, capacity))",
		)
	})

	// Retour de la fonction
	return nil, nil
}
