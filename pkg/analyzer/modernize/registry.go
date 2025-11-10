package modernize

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/modernize"
)

// Analyzers retourne tous les analyseurs modernize recommandés.
// Ces analyseurs suggèrent des simplifications du code Go en utilisant
// des fonctionnalités modernes du langage et de la bibliothèque standard.
//
// La suite modernize contient ~20 analyseurs qui détectent:
//   - Patterns obsolètes remplaçables par nouvelles features (Go 1.18+)
//   - Optimisations avec slices, maps packages (Go 1.21+)
//   - Simplifications avec min/max, range int (Go 1.21-1.22)
//   - Modernisation testing (t.Context, b.Loop, etc.)
//
// Chaque diagnostic inclut un fix automatique qui préserve le comportement.
//
// Note: Certains analyseurs peuvent être désactivés s'ils causent des problèmes.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs modernize
// Params: TODO
func Analyzers() []*analysis.Analyzer {
	// Liste des analyseurs à désactiver (bugs connus ou instables)
	disabled := map[string]bool{
		"newexpr": true, // Désactivé: panic dans certains cas (nil pointer dereference)
	}

	// Filtrer les analyseurs désactivés
	var filtered []*analysis.Analyzer
 // Verification de la condition
	for _, a := range modernize.Suite {
		// Vérification si l'analyseur est désactivé
		if !disabled[a.Name] {
			filtered = append(filtered, a)
		}
	}

	// Retour des analyseurs filtrés
	return filtered
}
