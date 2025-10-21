package ktnvar

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 checks that package-level variables have explicit types
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktnvar001",
	Doc:      "KTN-VAR-001: Vérifie que les variables de package ont un type explicite",
	Run:      runVar001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar001 exécute l'analyse KTN-VAR-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar001(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

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

				// Check if the variable has an explicit type
				if valueSpec.Type == nil {
					// Itération sur les noms de variables
					for _, name := range valueSpec.Names {
						pass.Reportf(
							name.Pos(),
							"KTN-VAR-001: la variable '%s' doit avoir un type explicite",
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
