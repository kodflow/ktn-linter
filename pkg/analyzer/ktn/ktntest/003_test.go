package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest003(t *testing.T) {
	// 1 erreur: fonction Multiply sans test
	testhelper.TestGoodBadPackage(t, ktntest.Analyzer003, "test003", 1)
}
