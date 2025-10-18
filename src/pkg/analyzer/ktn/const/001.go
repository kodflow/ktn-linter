package ktn_const

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
)

var Rule001 = &analysis.Analyzer{
	Name: "KTN_CONST_001",
	Doc:  "Vérifie que les constantes sont regroupées dans un bloc const()",
	Run:  runRule001,
}

func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			// Vérifier si c'est une déclaration non groupée (pas de parenthèses)
			if genDecl.Lparen == token.NoPos {
				reportUngroupedConst(pass, genDecl)
			}
		}
	}
	return nil, nil
}

// reportUngroupedConst signale une constante non groupée.
func reportUngroupedConst(pass *analysis.Pass, genDecl *ast.GenDecl) {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range valueSpec.Names {
			pass.Reportf(name.Pos(),
				"[KTN_CONST_001] Constante '%s' déclarée individuellement. Regroupez les constantes dans un bloc const ().\nExemple:\n  const (\n      %s %s = ...\n  )",
				name.Name, name.Name, utils.GetTypeString(valueSpec))
		}
	}
}
