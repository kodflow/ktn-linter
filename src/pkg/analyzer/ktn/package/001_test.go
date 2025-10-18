package ktn_package_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_package "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/package"
)

func TestRule001_NoDotImports(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_package.Rule001, "package001")
}
