// Diagnostic types for the cmd package.
package cmd

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// diagWithFset associe un diagnostic avec son FileSet et son analyseur.
type diagWithFset struct {
	diag         analysis.Diagnostic
	fset         *token.FileSet
	analyzerName string
}
