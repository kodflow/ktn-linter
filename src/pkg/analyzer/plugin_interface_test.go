package analyzer_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestNewAnalyzerPlugin vérifie le constructeur du plugin.
//
// Params:
//   - t: instance de test
func TestNewAnalyzerPlugin(t *testing.T) {
	plugin, err := analyzer.NewAnalyzerPlugin(nil)
	if err != nil {
		t.Fatalf("NewAnalyzerPlugin a retourné une erreur: %v", err)
	}

	if plugin == nil {
		t.Fatal("NewAnalyzerPlugin a retourné nil")
	}

	// Vérifier que le plugin retourné est utilisable
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers a retourné une erreur: %v", err)
	}

	if len(analyzers) == 0 {
		t.Fatal("BuildAnalyzers n'a retourné aucun analyseur")
	}
}
