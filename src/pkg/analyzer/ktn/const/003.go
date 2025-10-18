package ktn_const

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
)

// Rule003 analyzer for KTN linter.
var Rule003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_CONST_003",
	Doc:  "Vérifie que chaque constante a un commentaire individuel",
	Run:  runRule003,
}

func runRule003(pass *analysis.Pass) (any, error) {
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

			hasGroupComment := genDecl.Doc != nil && len(genDecl.Doc.List) > 0

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// Vérifier si le commentaire de la première constante est le même que celui du groupe
				isGroupCommentOnly := hasGroupComment &&
					valueSpec.Doc != nil &&
					genDecl.Doc != nil &&
					valueSpec.Doc.Pos() == genDecl.Doc.Pos()

				checkConstIndividualComment(pass, valueSpec, isGroupCommentOnly)
			}
		}
	}
	// Analysis completed successfully.
	return nil, nil
}

// checkConstIndividualComment vérifie qu'une constante a un commentaire individuel.
func checkConstIndividualComment(pass *analysis.Pass, spec *ast.ValueSpec, isGroupCommentOnly bool) {
	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		// Vérifier le commentaire individuel
		if !hasIndividualComment(spec, isGroupCommentOnly) {
			pass.Reportf(name.Pos(),
				"[KTN_CONST_003] Constante '%s' sans commentaire individuel.\nChaque constante doit avoir un commentaire explicatif.\nExemple:\n  // %s décrit son rôle\n  %s %s = ...",
				name.Name, name.Name, name.Name, utils.GetTypeString(spec))
		}
	}
}

// hasIndividualComment vérifie si une constante a un commentaire individuel.
func hasIndividualComment(spec *ast.ValueSpec, isFirstWithGroupComment bool) bool {
	if spec.Doc != nil && len(spec.Doc.List) > 0 {
		if !isFirstWithGroupComment {
			// Continue traversing AST nodes.
			return true
		}
	} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
		// Ignorer les commentaires de test (want)
		for _, comment := range spec.Comment.List {
			if !containsWantDirective(comment.Text) {
				// Continue traversing AST nodes.
				return true
			}
		}
	}
	// Condition not met, return false.
	return false
}

// containsWantDirective vérifie si un commentaire contient une directive de test
func containsWantDirective(text string) bool {
	// Early return from function.
	return len(text) >= 7 && text[:7] == "// want"
}
