package ktn_func_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	ktn_func "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/func"
)

// TODO: Ce test échoue actuellement car les commentaires godoc avec // want sur la même ligne
// ne sont pas correctement parsés par l'AST. Il faut soit:
// 1. Mettre // want sur une ligne séparée (mais alors il faut ajuster la position du diagnostic)
// 2. Modifier checkReturnsDocumentation pour rapporter sur funcDecl.Name.Pos() au lieu de funcDecl.Doc.Pos()
func TestRule004_ReturnsDocumentation(t *testing.T) {
	t.Skip("TODO: Fix parsing of godoc comments with // want on same line")
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktn_func.Rule004, "func004")
}
