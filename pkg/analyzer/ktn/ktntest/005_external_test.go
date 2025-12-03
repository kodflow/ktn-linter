package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest005(t *testing.T) {
	// 8 erreurs: tests sans table-driven pattern (>= 2 assertions)
	// - TestCalculatorTwoAssertions (2 assertions)
	// - TestStringLengthMultipleCases
	// - TestIsEmptyRepeatedAssertions
	// - TestToUpperManyScenarios
	// - TestContainsManyChecks
	// - TestCountWordsMultipleInputs
	// - TestWithAssertManyScenarios (testify/assert)
	// - TestWithRequireManyScenarios (testify/require)
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer005, "test005", "good_test.go", "bad_test.go", 8)
}
