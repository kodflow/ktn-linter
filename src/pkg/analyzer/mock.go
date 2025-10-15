//go:build test
// +build test

package analyzer

import "golang.org/x/tools/go/analysis"

// MockAnalyzerPlugin est un mock de AnalyzerPlugin pour les tests.
type MockAnalyzerPlugin struct {
	BuildAnalyzersFunc func() ([]*analysis.Analyzer, error)
	GetLoadModeFunc    func() string
}

// BuildAnalyzers implémente AnalyzerPlugin.BuildAnalyzers.
func (m *MockAnalyzerPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if m.BuildAnalyzersFunc != nil {
		return m.BuildAnalyzersFunc()
	}
	return nil, nil
}

// GetLoadMode implémente AnalyzerPlugin.GetLoadMode.
func (m *MockAnalyzerPlugin) GetLoadMode() string {
	if m.GetLoadModeFunc != nil {
		return m.GetLoadModeFunc()
	}
	return ""
}

// Vérification à la compilation que MockAnalyzerPlugin implémente AnalyzerPlugin
var _ AnalyzerPlugin = (*MockAnalyzerPlugin)(nil)
