package ktn_goroutine_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_goroutine "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/goroutine"
)

// TestRule001_GoroutineInLoop tests the functionality of the corresponding implementation.
func TestRule001_GoroutineInLoop(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_goroutine.Rule001, "goroutine001")
}
