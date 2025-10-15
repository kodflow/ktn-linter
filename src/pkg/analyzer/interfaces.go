package analyzer

import "golang.org/x/tools/go/analysis"

// AnalyzerPlugin définit l'interface pour les plugins d'analyseurs.
//
// Constructeur:
//   - NewAnalyzerPlugin(): crée une nouvelle instance du plugin
type AnalyzerPlugin interface {
	// BuildAnalyzers construit et retourne la liste des analyseurs.
	//
	// Returns:
	//   - []*analysis.Analyzer: la liste des analyseurs
	//   - error: erreur éventuelle lors de la construction
	BuildAnalyzers() ([]*analysis.Analyzer, error)

	// GetLoadMode retourne le mode de chargement requis.
	//
	// Returns:
	//   - string: le mode de chargement
	GetLoadMode() string
}
