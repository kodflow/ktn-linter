package ktn_control_flow_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
)

func TestForRule001_UnderscoreUnused(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_control_flow.RuleFor001, "for001")
}
