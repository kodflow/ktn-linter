package main

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
	}
}

// New creates a new instance of the plugin
func New(conf interface{}) ([]*analysis.Analyzer, error) {
	return (&AnalyzerPlugin{}).GetAnalyzers(), nil
}
