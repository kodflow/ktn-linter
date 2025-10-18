package ktn_pool_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_pool "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/pool"
)

func TestRule001_PoolGetWithDeferPut(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_pool.Rule001, "pool001")
}
