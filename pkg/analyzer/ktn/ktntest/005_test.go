package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest005(t *testing.T) {
	// 1 erreur: test sans table-driven pattern
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer005, "test005", "good_test.go", "bad_test.go", 1)
}
