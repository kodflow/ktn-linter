// Analyzer 005 for the ktnvar package.
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
	// ruleCodeVar005 is the rule code for this analyzer
	ruleCodeVar005 string = "KTN-VAR-005"
	// minMakeArgsVar008 is the minimum number of arguments for make call
	minMakeArgsVar008 int = 2
)

// Analyzer005 checks that make with length > 0 is avoided when append is used
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar005",
	Doc:      "KTN-VAR-005: Vérifie d'éviter make([]T, length) si utilisation avec append",
	Run:      runVar005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// checkMakeCallVar008 vérifie un appel à make avec length > 0.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction à vérifier
func checkMakeCallVar008(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérification que c'est un appel à make
	if !utils.IsMakeCall(call) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification du nombre d'arguments (2 ou 3: type, length, [capacity])
	if len(call.Args) < minMakeArgsVar008 {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(pass, call.Args[0]) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que la longueur est > 0
	if !utils.HasPositiveLength(pass, call.Args[1]) {
		// Continue traversing AST nodes.
		return
	}

	// Skip si VAR-016 s'applique (constante <= 1024 sans capacité)
	if len(call.Args) == minMakeArgsVar008 && utils.IsSmallConstantSize(pass, call.Args[1]) {
		// VAR-016 gère ce cas
		return
	}

	// Signalement de l'erreur
	pass.Reportf(
		call.Pos(),
		"KTN-VAR-005: utiliser make([]T, 0, capacity) au lieu de make([]T, length) pour éviter les zéro-values inutiles avant append",
	)
}

// runVar005 exécute l'analyse KTN-VAR-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar005(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar005) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		if cfg.IsFileExcluded(ruleCodeVar005, filename) {
			// Fichier exclu
			return
		}

		checkMakeCallVar008(pass, call)
	})

	// Retour de la fonction
	return nil, nil
}
