package ktn_interface_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestRule005InterfacesFileNotEmpty teste la règle 005 avec un fichier interfaces.go non vide.
func TestRule005InterfacesFileNotEmpty(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule005, "interface005")
}

// TestRule005PrivateOnly teste la règle 005 avec uniquement des interfaces privées.
func TestRule005PrivateOnly(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_interface.Rule005, "interface005_private_only")
}
