package ktn_pool_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_pool "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/pool"
)

// TestRule001_PoolGetWithDeferPut tests the functionality of the corresponding implementation.
func TestRule001_PoolGetWithDeferPut(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_pool.Rule001, "pool001")
}
