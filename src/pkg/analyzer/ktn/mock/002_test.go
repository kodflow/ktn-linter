package ktn_mock_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_mock "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/mock"
)

// TestRule002_InterfaceHasMock tests Rule002 when interface has mock.
//
// Params:
//   - t: testing instance
func TestRule002_InterfaceHasMock(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_mock.Rule002, "mock002")
}
