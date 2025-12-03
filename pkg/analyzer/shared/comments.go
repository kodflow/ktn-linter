// Shared utilities for comments handling.
package shared

import "go/ast"

// Constantes pour la vérification des directives "want"
const (
	// COMMENT_PREFIX_LENGTH est la longueur de "//" ou "/*"
	COMMENT_PREFIX_LENGTH int = 2
	// WANT_MIN_LENGTH est la longueur minimale pour "// want" ou "/*want"
	WANT_MIN_LENGTH int = 8 // "//" + " want" ou "/*" + " want"
	// WANT_KEYWORD_LENGTH longueur du mot "want"
	WANT_KEYWORD_LENGTH int = 4
	// WANT_WITH_SPACE_LENGTH longueur de " want"
	WANT_WITH_SPACE_LENGTH int = 5
)

// HasValidComment vérifie si un groupe de commentaires contient des commentaires valides et ignore les directives "want" utilisées par analysistest.
//
// Params:
//   - cg: groupe de commentaires à vérifier
//
// Returns:
//   - bool: true si commentaire valide trouvé (pas juste directives "want")
func HasValidComment(cg *ast.CommentGroup) bool {
	// Vérification de la condition
	if cg == nil || len(cg.List) == 0 {
		// Pas de commentaires
		return false
	}

	// Check if any comment is NOT a "want" directive
	// Itération sur les commentaires pour trouver un commentaire valide
	for _, comment := range cg.List {
		text := comment.Text

		// Skip "want" directives used by analysistest
		// Format: "// want ..." ou "/* want ..."
		if len(text) >= WANT_MIN_LENGTH {
			// Extraire le contenu après "//" ou "/*"
			content := text[COMMENT_PREFIX_LENGTH:]
			// Vérifier si c'est une directive "want"
			if len(content) >= WANT_WITH_SPACE_LENGTH && (content[:WANT_KEYWORD_LENGTH] == "want" || content[:WANT_WITH_SPACE_LENGTH] == " want") {
				continue
			}
		}

		// Found a valid comment
		return true
	}

	// Aucun commentaire valide trouvé
	return false
}
