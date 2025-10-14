package plugin_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/plugin"

	"golang.org/x/tools/go/analysis"
)

func TestNew(t *testing.T) {
	analyzers, err := plugin.New(nil)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if analyzers == nil {
		t.Fatal("New() returned nil")
	}

	if len(analyzers) == 0 {
		t.Fatal("New() returned empty analyzers list")
	}

	// Verify we have at least ConstAnalyzer and VarAnalyzer
	var foundConst, foundVar bool
	for _, a := range analyzers {
		if a.Name == "ktnconst" {
			foundConst = true
		}
		if a.Name == "ktnvar" {
			foundVar = true
		}
	}

	if !foundConst {
		t.Error("New() did not include ConstAnalyzer")
	}
	if !foundVar {
		t.Error("New() did not include VarAnalyzer")
	}
}

func TestGetAnalyzers(t *testing.T) {
	analyzers, err := plugin.New(nil)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if analyzers == nil {
		t.Fatal("New() returned nil")
	}

	if len(analyzers) == 0 {
		t.Fatal("New() returned empty list")
	}

	// Verify analyzers are valid
	for i, a := range analyzers {
		if a == nil {
			t.Errorf("analyzers[%d] is nil", i)
			continue
		}
		if a.Name == "" {
			t.Errorf("analyzers[%d] has empty name", i)
		}
		if a.Run == nil {
			t.Errorf("analyzers[%d] has nil Run function", i)
		}
	}
}

func TestIntegration(t *testing.T) {
	analyzers, err := plugin.New(nil)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	// Test that analyzers can be used
	for _, a := range analyzers {
		// Verify it's a valid analyzer
		if a == nil {
			t.Fatal("Got nil analyzer")
		}

		// Just checking that it conforms to the interface
		var _ *analysis.Analyzer = a
	}
}
