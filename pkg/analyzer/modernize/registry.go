// Package modernize wraps golang.org/x/tools modernize analyzers.
package modernize

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/modernize"
)

// Analyzers retourne tous les analyseurs modernize recommandés.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs modernize
func Analyzers() []*analysis.Analyzer {
	// Liste des analyseurs à désactiver (bugs connus ou instables avec Go 1.25.5)
	disabled := map[string]bool{
		"newexpr":        true, // Désactivé: panic dans certains cas (nil pointer dereference)
		"omitzero":       true, // Désactivé: panic avec Go 1.25.5 (nil pointer dans checkOmitEmptyField)
		"slicescontains": true, // Désactivé: panic avec Go 1.25.5 (nil pointer dans CoreType)
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
