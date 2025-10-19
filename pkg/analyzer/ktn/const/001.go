package ktnconst

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer001 checks that constants have explicit types
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktnconst001",
	Doc:      "KTN-CONST-001: Vérifie que les constantes ont un type explicite",
	Run:      runConst001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runConst001(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
   // Retour de la fonction
			return
		}

  // Itération sur les éléments
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Check if the constant has an explicit type
			if valueSpec.Type == nil {
				// If there are values, it's an error (not inheriting from iota pattern)
				// If there are no values, it's OK (inheriting type and value from previous line - iota pattern)
				if len(valueSpec.Values) > 0 {
     // Itération sur les éléments
					for _, name := range valueSpec.Names {
						pass.Reportf(
							name.Pos(),
							"KTN-CONST-001: la constante '%s' doit avoir un type explicite",
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
