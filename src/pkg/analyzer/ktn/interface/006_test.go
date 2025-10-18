package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestRule006_Reserved tests the functionality of the corresponding implementation.
// TestRule006Reserved tests the functionality of the corresponding implementation.
func TestRule006Reserved(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule006, "interface006")
}
