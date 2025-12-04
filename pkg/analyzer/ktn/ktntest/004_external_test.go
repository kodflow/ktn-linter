package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest004(t *testing.T) {
	// 2 erreurs: fonctions Multiply et Divide sans test
	testhelper.TestGoodBadPackage(t, ktntest.Analyzer004, "test004", 2)
}
