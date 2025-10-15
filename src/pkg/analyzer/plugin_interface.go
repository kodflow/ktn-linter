package analyzer

import "github.com/golangci/plugin-module-register/register"

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
