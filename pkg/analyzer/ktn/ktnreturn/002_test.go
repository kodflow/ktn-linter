package ktnreturn_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestReturn002(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnreturn.Analyzer002, "return002")
}
