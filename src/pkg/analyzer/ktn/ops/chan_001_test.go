package ktn_ops_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
)

// TestChanRule001_CloseByReceiver tests the functionality of the corresponding implementation.
func TestChanRule001_CloseByReceiver(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_ops.RuleChan001, "chan001")
}
