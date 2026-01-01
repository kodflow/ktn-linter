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
	// ruleCodeVar003 is the rule code for this analyzer
	ruleCodeVar003 string = "KTN-VAR-003"
)

// Analyzer003 checks that variables use camelCase (no underscores).
// This rule detects both SCREAMING_SNAKE_CASE and snake_case naming.
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar003",
	Doc:      "KTN-VAR-003: Les variables doivent utiliser camelCase (pas de underscores)",
	Run:      runVar003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar003 exécute l'analyse KTN-VAR-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar003(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar003) {
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
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar003, pass.Fset.Position(n.Pos()).Filename) {
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
			checkVar003Names(pass, valueSpec)
		}
	})

	// Return analysis result
	return nil, nil
}

// checkVar003Names vérifie les noms de variables pour les underscores.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVar003Names(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	// Iterate over variable names
	for _, name := range valueSpec.Names {
		varName := name.Name

		// Skip blank identifier
		if varName == "_" {
			continue
		}

		// Check for underscore in name (detects SCREAMING_SNAKE_CASE and snake_case)
		if hasUnderscore003(varName) {
			msg, ok := messages.Get(ruleCodeVar003)
			// Defensive: avoid panic if message is missing
			if !ok {
				pass.Reportf(name.Pos(), "%s: utiliser camelCase pour %q", ruleCodeVar003, varName)
				continue
			}
			pass.Reportf(
				name.Pos(),
				"%s: %s",
				ruleCodeVar003,
				msg.Format(config.Get().Verbose, varName),
			)
		}
	}
}

// hasUnderscore003 vérifie si un nom contient un underscore.
//
// Params:
//   - name: nom à vérifier
//
// Returns:
//   - bool: true si le nom contient un underscore
func hasUnderscore003(name string) bool {
	// Blank identifier is allowed
	if name == "_" {
		// Retour de la fonction
		return false
	}

	// Check for underscore in name
	return strings.Contains(name, "_")
}
