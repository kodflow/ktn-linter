package ktn_alloc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_alloc "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/alloc"
)

func TestRule002_MakeSlicePrealloc(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_alloc.Rule002, "alloc002")
}
