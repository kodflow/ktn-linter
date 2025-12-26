package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest004(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "table-driven pattern detection",
			analyzer: "test004",
		},
		{
			name:     "validate test structure compliance",
			analyzer: "test004",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 8 erreurs: tests sans table-driven pattern (>= 2 assertions)
			// - TestCalculatorTwoAssertions (2 assertions)
			// - TestStringLengthMultipleCases
			// - TestIsEmptyRepeatedAssertions
			// - TestToUpperManyScenarios
			// - TestContainsManyChecks
			// - TestCountWordsMultipleInputs
			// - TestWithAssertManyScenarios (testify/assert)
			// - TestWithRequireManyScenarios (testify/require)
			testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer004, tt.analyzer, "good_test.go", "bad_test.go", 8)
		})
	}
}
