package ktn_mock_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_mock "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/mock"
)

// TestRule001_MockFileExists tests Rule001 when mock file exists.
//
// Params:
//   - t: testing instance
func TestRule001_MockFileExists(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock001")
}

// TestRule001_NoInterfaces tests Rule001 with no interfaces.
//
// Params:
//   - t: testing instance
func TestRule001_NoInterfaces(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock003")
}

// TestRule001_OnlyConstants tests Rule001 with only constants.
//
// Params:
//   - t: testing instance
func TestRule001_OnlyConstants(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock004")
}
