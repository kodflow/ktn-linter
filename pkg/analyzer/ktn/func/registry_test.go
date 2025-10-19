package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
)

func TestAnalyzers(t *testing.T) {
	analyzers := ktnfunc.GetAnalyzers()

	// Check that we have at least 11 analyzers
	minExpectedCount := 11
	if len(analyzers) < minExpectedCount {
		t.Errorf("Analyzers() returned %d analyzers, expected at least %d", len(analyzers), minExpectedCount)
	}

	// Check that all analyzers are non-nil
	for i, analyzer := range analyzers {
		if analyzer == nil {
			t.Errorf("Analyzer at index %d is nil", i)
		}
	}

	// Check that the analyzers have the expected names
	expectedNames := map[string]bool{
		"ktnfunc001": true, // Max 35 lines
		"ktnfunc002": true, // Max 5 parameters
		"ktnfunc003": true, // Function names must start with a verb
		"ktnfunc004": true, // No naked returns
		"ktnfunc005": true, // Max cyclomatic complexity 10
		"ktnfunc006": true, // Error last
		"ktnfunc007": true, // Documentation stricte
		"ktnfunc008": true, // Context must be first parameter
		"ktnfunc009": true, // No side effects in getters
		"ktnfunc010": true, // Named returns for >3 return values
		"ktnfunc011": true, // Comments on branches/returns
	}

	for _, analyzer := range analyzers {
		if !expectedNames[analyzer.Name] {
			t.Errorf("Unexpected analyzer name: %s", analyzer.Name)
		}
	}
}
