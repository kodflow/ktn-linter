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

// Analyzer002 checks that package-level variables have explicit type AND value
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar002",
	Doc:      "KTN-VAR-002: Les variables de package doivent avoir le format 'var name type = value'",
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

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar002, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Check package-level declarations only
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip if not a GenDecl
			if !ok {
				// Not a general declaration
				continue
			}

			// Only check var declarations
			if genDecl.Tok != token.VAR {
				// Continue traversing AST nodes.
				continue
			}

			// Itération sur les spécifications
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				// Vérifier si le type est explicite ou visible dans la valeur
				checkVarSpec(pass, valueSpec)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkVarSpec vérifie une spécification de variable.
// Style requis: var name type (= value optionnel, zéro-value accepté)
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVarSpec(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	hasExplicitType := valueSpec.Type != nil

	// Seule exigence: type explicite obligatoire
	// L'initialisation est optionnelle (zéro-value idiomatique en Go)
	if !hasExplicitType {
		// Parcourir les noms
		for _, name := range valueSpec.Names {
			// Ignorer les blank identifiers
			if name.Name == "_" {
				continue
			}

			msg, _ := messages.Get(ruleCodeVar002)
			pass.Reportf(
				name.Pos(),
				"%s: %s",
				ruleCodeVar002,
				msg.Format(config.Get().Verbose, name.Name),
			)
		}
	}
	// Type explicite présent = OK (avec ou sans valeur)
}
