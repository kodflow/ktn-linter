package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestRule002NoPublicStructs teste la règle 002 sans structures publiques.
func TestRule002NoPublicStructs(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule002, "interface002")
}

// TestRule002AllowedTypes teste la règle 002 avec des types autorisés.
func TestRule002AllowedTypes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule002, "interface002_allowed")
}
