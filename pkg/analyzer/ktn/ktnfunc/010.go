// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc010 is the rule code for this analyzer
	ruleCodeFunc010 string = "KTN-FUNC-010"
	// defaultMaxLinesForNakedReturn max lines for naked return
	defaultMaxLinesForNakedReturn int = 5
)

// Analyzer010 checks that naked returns are only used in very short functions
var Analyzer010 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc010",
	Doc:      "KTN-FUNC-010: Les naked returns sont interdits sauf pour les fonctions très courtes (<5 lignes)",
	Run:      runFunc010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc010 exécute l'analyse KTN-FUNC-010.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc010(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc010) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	maxLinesNaked := cfg.GetThreshold(ruleCodeFunc010, defaultMaxLinesForNakedReturn)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeFunc010, filename) {
			// Fichier exclu
			return
		}

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		// Analyze naked returns
		analyzeNakedReturns(pass, funcDecl, maxLinesNaked)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeNakedReturns analyzes naked returns in a function.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: fonction à analyser
//   - maxLinesNaked: max lines for naked return
func analyzeNakedReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl, maxLinesNaked int) {
	funcName := funcDecl.Name.Name

	// Skip if function doesn't have named return values
	if !hasNamedReturns(funcDecl.Type.Results) {
		// Retour de la fonction
		return
	}

	// Count the lines of the function
	pureLines := countPureCodeLines(pass, funcDecl.Body)

	// Check for naked returns
	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		ret, ok := node.(*ast.ReturnStmt)
		// Vérification de la condition
		if !ok {
			// Retour de la fonction
			return true
		}

		// Naked return has no results specified
		if len(ret.Results) == 0 {
			// Allow naked returns in very short functions
			if pureLines >= maxLinesNaked {
				// Rapport d'erreur pour naked return interdit
				pass.Reportf(
					ret.Pos(),
					"KTN-FUNC-010: naked return interdit dans la fonction '%s' (%d lignes, max: %d pour naked return)",
					funcName,
					pureLines,
					maxLinesNaked-1,
				)
			}
		}

		// Retour de la fonction
		return true
	})
}

// hasNamedReturns checks if the function has named return values
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si retours nommés
func hasNamedReturns(results *ast.FieldList) bool {
	// Vérification de la condition
	if results == nil || len(results.List) == 0 {
		// Retour de la fonction
		return false
	}

	// Itération sur les éléments
	for _, field := range results.List {
		// Vérification de la condition
		if len(field.Names) > 0 {
			// Retour de la fonction
			return true
		}
	}

	// Retour de la fonction
	return false
}
