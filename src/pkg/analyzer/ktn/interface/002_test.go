package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

func TestRule002_NoPublicStructs(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule002, "interface002")
}
