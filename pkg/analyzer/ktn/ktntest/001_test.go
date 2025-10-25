package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest001(t *testing.T) {
	// 1 erreur: mauvais nom de package
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer001, "test001", "good_test.go", "bad_test.go", 1)
}
