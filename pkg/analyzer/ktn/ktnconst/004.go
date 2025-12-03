// Analyzer 004 for the ktnconst package.
package ktnconst

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer004 checks that every constant has an associated comment
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnconst004",
	Doc:      "KTN-CONST-004: Vérifie que chaque constante a un commentaire associé",
	Run:      runConst004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runConst004 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runConst004(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

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
						"KTN-CONST-004: la constante '%s' doit avoir un commentaire associé",
						name.Name,
					)
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
