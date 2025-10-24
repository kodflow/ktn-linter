package ktnvar

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Constantes pour la détection des directives "want" dans les commentaires
const (
	WANT_MIN_LENGTH      int = 6 // longueur minimale pour "// want"
	WANT_SPACE_MIN_LENGTH int = 7 // longueur minimale pour "// want"
	COMMENT_PREFIX_END   int = 2 // fin du préfixe "//" dans un commentaire
	WANT_DIRECTIVE_END   int = 6 // position de fin de "want" dans "// want"
	WANT_SPACE_END       int = 7 // position de fin de " want" dans "//  want"
)

// Analyzer004 checks that every package-level variable has an associated comment
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar004",
	Doc:      "KTN-VAR-004: Vérifie que chaque variable de package a un commentaire associé",
	Run:      runVar004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar004 exécute l'analyse KTN-VAR-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar004(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations only
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

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
			hasGenDeclDoc := hasValidComment(genDecl.Doc)

			// Itération sur les éléments
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)

				// Check if this specific ValueSpec has a doc comment or line comment
				// Filter out "want" directives used by analysistest
				hasValueSpecDoc := hasValidComment(valueSpec.Doc)
				hasValueSpecComment := hasValidComment(valueSpec.Comment)

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
							"KTN-VAR-004: la variable '%s' doit avoir un commentaire associé",
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

// hasValidComment vérifie si un groupe de commentaires contient des commentaires valides.
//
// Params:
//   - cg: groupe de commentaires à vérifier
//
// Returns:
//   - bool: true si commentaire valide (pas juste directives "want")
func hasValidComment(cg *ast.CommentGroup) bool {
	// Vérification de la condition
	if cg == nil || len(cg.List) == 0 {
		// Retour de la fonction
		return false
	}

	// Check if any comment is NOT a "want" directive
	// Itération sur les commentaires pour trouver un commentaire valide
	for _, comment := range cg.List {
		text := comment.Text
		// Skip "want" directives used by analysistest
		// Vérification si c'est un commentaire de ligne "want"
		if len(text) >= WANT_MIN_LENGTH && text[COMMENT_PREFIX_END:WANT_DIRECTIVE_END] == "want" {
			continue
		}
		// Vérification si c'est un commentaire de bloc "want"
		if len(text) >= WANT_SPACE_MIN_LENGTH && text[COMMENT_PREFIX_END:WANT_SPACE_END] == " want" {
			continue
		}
		// Found a valid comment
		return true
	}

	// Retour de la fonction
	return false
}
