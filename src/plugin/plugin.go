package plugin

import (
	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin implements the golangci-lint plugin interface
type AnalyzerPlugin struct{}

// GetAnalyzers returns all custom analyzers for golangci-lint
func (*AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.ConstAnalyzer,
		analyzer.VarAnalyzer,
		analyzer.FuncAnalyzer,
	}
}

// New creates a new instance of the plugin.
//
// Params:
//   - conf: la configuration du plugin (non utilisée actuellement)
//
// Returns:
//   - []*analysis.Analyzer: la liste des analyseurs fournis par le plugin
//   - error: toujours nil dans l'implémentation actuelle
func New(conf interface{}) ([]*analysis.Analyzer, error) {
	return (&AnalyzerPlugin{}).GetAnalyzers(), nil
}
