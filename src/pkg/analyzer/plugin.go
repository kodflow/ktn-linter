package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// plugin implements the golangci-lint module plugin interface
type plugin struct{}

// New creates a new instance of the plugin.
//
// Params:
//   - settings: les paramètres de configuration du plugin
//
// Returns:
//   - register.LinterPlugin: l'instance du plugin créée
//   - error: toujours nil dans l'implémentation actuelle
func New(settings any) (register.LinterPlugin, error) {
	// Retourne une nouvelle instance du plugin
	return &plugin{}, nil
}

// BuildAnalyzers returns all analyzers provided by this plugin.
//
// Returns:
//   - []*analysis.Analyzer: la liste des analyseurs fournis par ce plugin
//   - error: toujours nil dans l'implémentation actuelle
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		ConstAnalyzer,
		VarAnalyzer,
		FuncAnalyzer,
		StructAnalyzer,
		InterfaceAnalyzer,
	}, nil
}

// GetLoadMode returns the load mode for the analyzers.
//
// Returns:
//   - string: le mode de chargement (LoadModeSyntax)
func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
