package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestInterface001(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktninterface.Analyzer001, "interface001")
}
