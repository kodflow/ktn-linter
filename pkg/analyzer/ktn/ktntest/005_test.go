package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest005(t *testing.T) {
	// 5 erreurs: tests sans table-driven pattern
	// - TestStringLengthMultipleCases
	// - TestIsEmptyRepeatedAssertions
	// - TestToUpperManyScenarios
	// - TestContainsManyChecks
	// - TestCountWordsMultipleInputs
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer005, "test005", "good_test.go", "bad_test.go", 5)
}
