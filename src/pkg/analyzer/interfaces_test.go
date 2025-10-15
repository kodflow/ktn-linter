package analyzer_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestAnalyzerPluginInterface vérifie que l'interface AnalyzerPlugin est bien définie.
//
// Params:
//   - t: instance de test
func TestAnalyzerPluginInterface(t *testing.T) {
	plugin, err := analyzer.New(nil)
	if err != nil {
		t.Fatalf("Erreur lors de la création du plugin: %v", err)
	}

	// Vérifier que le plugin implémente l'interface attendue
	if plugin == nil {
		t.Fatal("Le plugin créé est nil")
	}

	// Vérifier BuildAnalyzers
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers a retourné une erreur: %v", err)
	}
	if len(analyzers) == 0 {
		t.Fatal("BuildAnalyzers n'a retourné aucun analyseur")
	}

	// Vérifier GetLoadMode
	loadMode := plugin.GetLoadMode()
	if loadMode == "" {
		t.Fatal("GetLoadMode a retourné une chaîne vide")
	}
}
