// Analyzer 003 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeComment003 is the rule code for this analyzer
	ruleCodeComment003 string = "KTN-COMMENT-003"
)

// Analyzer003 checks that every constant has an associated comment
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment003",
	Doc:      "KTN-COMMENT-003: Vérifie que chaque constante a un commentaire associé",
	Run:      runComment003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment003 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runComment003(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment003) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeComment003, filename) {
			// File excluded by configuration
			return
		}

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			// Retour de la fonction
			return
		}

		// Check if the GenDecl has a doc comment (applies to all constants in the group)
		// Filter out "want" directives used by analysistest
		hasGenDeclDoc := shared.HasValidComment(genDecl.Doc)

		// Itération sur les éléments
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Check if this specific ValueSpec has a doc comment or line comment
			// Filter out "want" directives used by analysistest
			hasValueSpecDoc := shared.HasValidComment(valueSpec.Doc)
			hasValueSpecComment := shared.HasValidComment(valueSpec.Comment)

			// A constant is considered documented if:
			// 1. The GenDecl has a doc comment (group documentation), OR
			// 2. The ValueSpec has a doc comment (above the constant), OR
			// 3. The ValueSpec has a line comment (on the same line)
			hasComment := hasGenDeclDoc || hasValueSpecDoc || hasValueSpecComment

			// Vérification de la condition
			if !hasComment {
				// Itération sur les éléments
				for _, name := range valueSpec.Names {
					pass.Reportf(
						name.Pos(),
						"KTN-COMMENT-003: la constante '%s' doit avoir un commentaire associé",
						name.Name,
					)
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
