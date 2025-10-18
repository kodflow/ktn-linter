package ktnconst

import (
	"go/ast"
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer005 checks that all constants are exported (start with uppercase letter)
var Analyzer005 = &analysis.Analyzer{
	Name:     "ktnconst005",
	Doc:      "KTN-CONST-005: Vérifie que toutes les constantes sont exportées (commencent par une majuscule)",
	Run:      runConst005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runConst005(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			return
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Check each constant name
			for _, name := range valueSpec.Names {
				// Skip blank identifiers
				if name.Name == "_" {
					continue
				}

				// Check if the first character is uppercase
				if len(name.Name) > 0 {
					firstChar := rune(name.Name[0])
					if !unicode.IsUpper(firstChar) {
						pass.Reportf(
							name.Pos(),
							"KTN-CONST-005: la constante '%s' doit être exportée (commencer par une majuscule)",
							name.Name,
						)
					}
				}
			}
		}
	})

	return nil, nil
}
