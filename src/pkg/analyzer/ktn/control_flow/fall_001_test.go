package ktn_control_flow_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
)

// TestFallRule001_FallthroughOutsideSwitch tests the functionality of the corresponding implementation.
func TestFallRule001_FallthroughOutsideSwitch(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_control_flow.RuleFall001, "fall001")
}
