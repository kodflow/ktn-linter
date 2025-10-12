package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// Plugin implements the golangci-lint module plugin interface
type Plugin struct{}

// New creates a new instance of the plugin
func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

// BuildAnalyzers returns all analyzers provided by this plugin
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		ConstAnalyzer,
	}, nil
}

// GetLoadMode returns the load mode for the analyzers
func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
