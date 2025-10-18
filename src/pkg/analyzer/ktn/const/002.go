package ktn_const

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule002 analyzer for KTN linter.
var Rule002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_CONST_002",
	Doc:  "Vérifie que les groupes de constantes ont un commentaire de groupe",
	Run:  runRule002,
}

func runRule002(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			// Vérifier uniquement les groupes (avec parenthèses)
			if genDecl.Lparen == token.NoPos {
				continue
			}

			// Vérifier la présence d'un commentaire de groupe
			hasGroupComment := genDecl.Doc != nil && len(genDecl.Doc.List) > 0

			if !hasGroupComment {
				pass.Reportf(genDecl.Pos(),
					"[KTN_CONST_002] Groupe de constantes sans commentaire de groupe.\nAjoutez un commentaire avant le bloc const () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de constantes\n  const (...)")
			}
		}
	}
	// Analysis completed successfully.
	return nil, nil
}
