// Analyzer 006 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer006 checks that package-level variables are declared after constants
var Analyzer006 = &analysis.Analyzer{
	Name:     "ktnvar006",
	Doc:      "KTN-VAR-006: Vérifie que les variables de package sont déclarées après les constantes",
	Run:      runVar006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar006 exécute l'analyse KTN-VAR-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar006(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	// Process each file
	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)
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
				pass.Reportf(
					genDecl.Pos(),
					"KTN-VAR-006: les constantes doivent être déclarées avant les variables",
				)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
