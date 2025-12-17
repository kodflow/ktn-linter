// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar018 is the rule code for this analyzer
	ruleCodeVar018 string = "KTN-VAR-018"
)

// Analyzer018 checks that variables don't use snake_case naming
var Analyzer018 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar018",
	Doc:      "KTN-VAR-018: Les variables ne doivent pas utiliser snake_case (underscores)",
	Run:      runVar018,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar018 exécute l'analyse KTN-VAR-018.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar018(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar018) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar018, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Only check var declarations
		if genDecl.Tok != token.VAR {
			// Continue traversing AST nodes
			return
		}

		// Iterate over specifications
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			// Check each variable name
			checkVar018Names(pass, valueSpec)
		}
	})

	// Return analysis result
	return nil, nil
}

// checkVar018Names vérifie les noms de variables pour snake_case.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVar018Names(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	// Iterate over variable names
	for _, name := range valueSpec.Names {
		varName := name.Name

		// Skip blank identifier
		if varName == "_" {
			continue
		}

		// Check for snake_case (lowercase with underscores)
		if isSnakeCase(varName) {
			msg, _ := messages.Get(ruleCodeVar018)
			pass.Reportf(
				name.Pos(),
				"%s: %s",
				ruleCodeVar018,
				msg.Format(config.Get().Verbose, varName),
			)
		}
	}
}

// isSnakeCase vérifie si un nom utilise snake_case (minuscules avec underscores).
// Exclut SCREAMING_SNAKE_CASE qui est déjà géré par VAR-001.
//
// Params:
//   - name: nom à vérifier
//
// Returns:
//   - bool: true si le nom est en snake_case
func isSnakeCase(name string) bool {
	// Must contain underscore
	if !strings.Contains(name, "_") {
		// Retour de la fonction
		return false
	}

	// Check if it's NOT all uppercase (SCREAMING_SNAKE_CASE is handled by VAR-001)
	for _, ch := range name {
		// If we find a lowercase letter, it's snake_case
		if ch >= 'a' && ch <= 'z' {
			// Retour de la fonction
			return true
		}
	}

	// All uppercase = SCREAMING_SNAKE_CASE, not snake_case
	return false
}

