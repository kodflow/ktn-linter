// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

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
	// ruleCodeVar002 is the rule code for this analyzer
	ruleCodeVar002 string = "KTN-VAR-002"
)

// Analyzer002 checks that package-level variables are declared after constants
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar002",
	Doc:      "KTN-VAR-002: Vérifie que les variables de package sont déclarées après les constantes",
	Run:      runVar002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar002 exécute l'analyse KTN-VAR-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar002) {
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
		(*ast.File)(nil),
	}

	// Process each file
	insp.Preorder(nodeFilter, func(n ast.Node) {
		file, ok := n.(*ast.File)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar002, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}
		varSeen := false

		// Check declarations in order
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip non-GenDecl (functions, etc.)
			if !ok {
				continue
			}

			// Track variable declarations
			if genDecl.Tok == token.VAR {
				varSeen = true
			}

			// Error: const after var
			if genDecl.Tok == token.CONST && varSeen {
				msg, ok := messages.Get(ruleCodeVar002)
				// Defensive: avoid panic if message is missing
				if !ok {
					pass.Reportf(genDecl.Pos(), "%s: const doit être déclaré avant var", ruleCodeVar002)
					continue
				}
				pass.Reportf(
					genDecl.Pos(),
					"%s: %s",
					ruleCodeVar002,
					msg.Format(config.Get().Verbose),
				)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
