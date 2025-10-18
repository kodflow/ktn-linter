package ktn_data_structures_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_data_structures "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/data_structures"
)

func TestArrayRule001_InconsistentSize(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_data_structures.RuleArray001, "array001")
}
