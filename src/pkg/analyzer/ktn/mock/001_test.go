package ktn_mock_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_mock "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/mock"
)

func TestRule001_MockFileExists(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock001")
}

func TestRule001_NoInterfaces(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock003")
}

func TestRule001_OnlyConstants(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule001, "mock004")
}
