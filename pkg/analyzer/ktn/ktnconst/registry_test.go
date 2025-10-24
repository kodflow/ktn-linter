package ktnconst

import (
	"testing"
)

// TestGetAnalyzers tests GetAnalyzers returns all analyzers
func TestGetAnalyzers(t *testing.T) {
	analyzers := GetAnalyzers()

	// Check that we have exactly 4 analyzers
	expectedCount := 4
	if len(analyzers) != expectedCount {
		t.Errorf("GetAnalyzers() returned %d analyzers, expected %d", len(analyzers), expectedCount)
	}

	// Check that all analyzers are non-nil
	for i, analyzer := range analyzers {
		if analyzer == nil {
			t.Errorf("Analyzer at index %d is nil", i)
		}
	}

	// Check that the analyzers have the expected names
	expectedNames := []string{
		"ktnconst001",
		"ktnconst002",
		"ktnconst003",
		"ktnconst004",
	}

	for i, analyzer := range analyzers {
		if analyzer.Name != expectedNames[i] {
			t.Errorf("Analyzer at index %d has name %q, expected %q", i, analyzer.Name, expectedNames[i])
		}
	}
}
