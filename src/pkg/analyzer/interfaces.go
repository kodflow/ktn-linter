package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

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

// NewAnalyzerPlugin crée une nouvelle instance du plugin d'analyseur.
//
// Params:
//   - settings: les paramètres de configuration du plugin
//
// Returns:
//   - register.LinterPlugin: l'instance du plugin créée
//   - error: erreur éventuelle lors de la création
func NewAnalyzerPlugin(settings any) (register.LinterPlugin, error) {
	// Retourne une nouvelle instance via New()
	return New(settings)
}
