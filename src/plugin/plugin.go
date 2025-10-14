package plugin

import (
	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// analyzerPlugin implements the golangci-lint plugin interface
type analyzerPlugin struct{}

// GetAnalyzers returns all custom analyzers for golangci-lint
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.ConstAnalyzer,
		analyzer.VarAnalyzer,
		analyzer.FuncAnalyzer,
		analyzer.InterfaceAnalyzer,
		analyzer.TestAnalyzer,
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
	return (&analyzerPlugin{}).GetAnalyzers(), nil
}
