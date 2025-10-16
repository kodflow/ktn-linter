package analyzer_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestNewAnalyzerPlugin teste la création d'un plugin d'analyseur.
//
// Params:
//   - t: l'instance de test
func TestNewAnalyzerPlugin(t *testing.T) {
	plugin, err := analyzer.NewAnalyzerPlugin(nil)
	if err != nil {
		t.Fatalf("NewAnalyzerPlugin() retourné une erreur: %v", err)
	}

	if plugin == nil {
		t.Fatal("NewAnalyzerPlugin() retourné nil")
	}

	// Vérifier que BuildAnalyzers retourne des analyseurs
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers() retourné une erreur: %v", err)
	}

	if len(analyzers) == 0 {
		t.Error("BuildAnalyzers() retourné aucun analyseur")
	}

	// Vérifier que GetLoadMode retourne une valeur
	loadMode := plugin.GetLoadMode()
	if loadMode == "" {
		t.Error("GetLoadMode() retourné une chaîne vide")
	}
}
