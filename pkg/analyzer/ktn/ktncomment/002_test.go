package ktncomment_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestComment002(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktncomment.Analyzer002, "comment002")
}
