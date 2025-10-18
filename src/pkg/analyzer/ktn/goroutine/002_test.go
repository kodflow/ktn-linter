package ktn_goroutine_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_goroutine "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/goroutine"
)

// TestRule002_GoroutineSynchronization tests the functionality of the corresponding implementation.
func TestRule002_GoroutineSynchronization(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_goroutine.Rule002, "goroutine002")
}
