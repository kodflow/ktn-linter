// Analyzer 004 for the ktncomment package.
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
	// ruleCodeComment004 is the rule code for this analyzer
	ruleCodeComment004 string = "KTN-COMMENT-004"
)

// Analyzer004 checks that every package-level variable has an associated comment
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment004",
	Doc:      "KTN-COMMENT-004: Vérifie que chaque variable de package a un commentaire associé",
	Run:      runComment004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment004 exécute l'analyse KTN-COMMENT-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runComment004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment004) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations only
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		if cfg.IsFileExcluded(ruleCodeComment004, filename) {
			// Fichier exclu
			return
		}

		// Check package-level declarations only
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip if not a GenDecl
			if !ok {
				continue
			}

			// Only check var declarations
			if genDecl.Tok != token.VAR {
				continue
			}

			// Check if the GenDecl has a doc comment (applies to all variables in the group)
			// Filter out "want" directives used by analysistest
			hasGenDeclDoc := shared.HasValidComment(genDecl.Doc)

			// Itération sur les éléments
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)

				// Check if this specific ValueSpec has a doc comment or line comment
				// Filter out "want" directives used by analysistest
				hasValueSpecDoc := shared.HasValidComment(valueSpec.Doc)
				hasValueSpecComment := shared.HasValidComment(valueSpec.Comment)

				// A variable is considered documented if:
				// 1. The GenDecl has a doc comment (group documentation), OR
				// 2. The ValueSpec has a doc comment (above the variable), OR
				// 3. The ValueSpec has a line comment (on the same line)
				hasComment := hasGenDeclDoc || hasValueSpecDoc || hasValueSpecComment

				// Vérification de la condition
				if !hasComment {
					// Itération sur les éléments
					for _, name := range valueSpec.Names {
						pass.Reportf(
							name.Pos(),
							"KTN-COMMENT-004: la variable '%s' doit avoir un commentaire associé",
							name.Name,
						)
					}
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
