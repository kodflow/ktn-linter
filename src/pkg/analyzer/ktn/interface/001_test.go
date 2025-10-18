package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestRule001InterfacesFileRequired teste la règle 001 avec un fichier interfaces.go requis.
func TestRule001InterfacesFileRequired(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule001, "interface001")
}

// TestRule001AllowedTypes teste la règle 001 avec des types autorisés.
func TestRule001AllowedTypes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule001, "interface001_allowed_types")
}

// TestRule001WithInterface teste la règle 001 avec une interface.
func TestRule001WithInterface(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule001, "interface001_with_interface")
}
