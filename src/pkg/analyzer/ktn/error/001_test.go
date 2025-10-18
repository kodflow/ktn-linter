package ktn_error_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_error "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/error"
)

func TestRule001_ErrorWrapping(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_error.Rule001, "error001")
}
