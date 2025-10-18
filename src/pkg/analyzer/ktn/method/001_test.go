package ktn_method_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_method "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/method"
)

// TestRule001_PointerReceiverForMutation tests the functionality of the corresponding implementation.
func TestRule001_PointerReceiverForMutation(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_method.Rule001, "method001")
}
