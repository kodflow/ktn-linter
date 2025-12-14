// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc006 is the rule code for this analyzer
	ruleCodeFunc006 string = "KTN-FUNC-006"
	// defaultMaxParams max params allowed in a function (context.Context excluded)
	defaultMaxParams int = 5
)

// Analyzer006 checks that functions don't have more than MAX_PARAMS parameters
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc006",
	Doc:      "KTN-FUNC-006: Les fonctions ne doivent pas dépasser 5 paramètres",
	Run:      runFunc006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc006 exécute l'analyse KTN-FUNC-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc006(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc006) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	maxParams := cfg.GetThreshold(ruleCodeFunc006, defaultMaxParams)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeFunc006, filename) {
			// Fichier exclu
			return
		}

		// Analyze function parameters
		analyzeFuncParams(pass, n, maxParams)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeFuncParams analyzes function parameters for FUNC-006.
//
// Params:
//   - pass: contexte d'analyse
//   - n: node to analyze
//   - maxParams: max parameters allowed
func analyzeFuncParams(pass *analysis.Pass, n ast.Node, maxParams int) {
	var funcType *ast.FuncType
	var pos ast.Node
	var name string

	// Sélection selon la valeur
	switch fn := n.(type) {
	// Traitement
	case *ast.FuncDecl:
		funcType = fn.Type
		pos = fn.Name
		name = fn.Name.Name

		// Skip test functions
		if shared.IsTestFunction(fn) {
			// Retour pour ignorer les fonctions de test
			return
		}
	// Traitement
	case *ast.FuncLit:
		funcType = fn.Type
		pos = fn
		name = "function literal"
	}

	// Count total parameters (excluding context.Context)
	paramCount := countEffectiveParams(pass, funcType.Params)

	// Vérification de la condition
	if paramCount > maxParams {
		msg, _ := messages.Get(ruleCodeFunc006)
		pass.Reportf(
			pos.Pos(),
			"%s: %s",
			ruleCodeFunc006,
			msg.Format(config.Get().Verbose, name, paramCount, maxParams),
		)
	}
}

// countEffectiveParams counts parameters excluding context.Context.
// context.Context is excluded because KTN-FUNC-002 already mandates it as first param.
//
// Params:
//   - pass: contexte d'analyse
//   - params: field list of function parameters
//
// Returns:
//   - int: effective parameter count
func countEffectiveParams(pass *analysis.Pass, params *ast.FieldList) int {
	// Check for nil params
	if params == nil {
		// Retour 0 si pas de paramètres
		return 0
	}

	count := 0
	// Iterate over parameter fields
	for _, field := range params.List {
		// Skip context.Context parameters (including aliases)
		if isContextTypeWithPass(pass, field.Type) {
			// Passage au champ suivant si context.Context
			continue
		}

		// Each field can declare multiple params: func(a, b, c int)
		if len(field.Names) > 0 {
			// Ajout du nombre de noms
			count += len(field.Names)
		} else {
			// Unnamed parameter (e.g., in interface or func literal)
			// Incrément pour paramètre sans nom
			count++
		}
	}

	// Retour du compte effectif
	return count
}
