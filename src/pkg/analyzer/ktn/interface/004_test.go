package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestRule004ConstructorsRequired teste la règle 004 avec des constructeurs requis.
func TestRule004ConstructorsRequired(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule004, "interface004")
}

// TestRule004EmptyInterface teste la règle 004 avec une interface vide.
func TestRule004EmptyInterface(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule004, "interface004_empty_interface")
}
