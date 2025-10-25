package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest002(t *testing.T) {
	// 1 erreur: bad_test.go sans fichier bad.go correspondant
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer002, "test002", "good_test.go", "bad_test.go", 1)
}
