// Shared utilities for comments handling.
package shared

import "go/ast"

// Constantes pour la vérification des directives "want"
const (
	// commentPrefixLength est la longueur de "//" ou "/*"
	commentPrefixLength int = 2
	// wantMinLength est la longueur minimale pour "// want" ou "/*want"
	wantMinLength int = 8 // "//" + " want" ou "/*" + " want"
	// wantKeywordLength longueur du mot "want"
	wantKeywordLength int = 4
	// wantWithSpaceLength longueur de " want"
	wantWithSpaceLength int = 5
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
		if len(text) >= wantMinLength {
			// Extraire le contenu après "//" ou "/*"
			content := text[commentPrefixLength:]
			// Vérifier si c'est une directive "want"
			if len(content) >= wantWithSpaceLength && (content[:wantKeywordLength] == "want" || content[:wantWithSpaceLength] == " want") {
				continue
			}
		}

		// Found a valid comment
		// Retour commentaire valide trouvé
		return true
	}

	// Aucun commentaire valide trouvé
	return false
}
