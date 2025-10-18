package ktn_interface

import (
	"golang.org/x/tools/go/analysis"
)

var Rule006 = &analysis.Analyzer{
	Name: "KTN_INTERFACE_006",
	Doc:  "Règle INTERFACE-006 (réservée pour usage futur)",
	Run:  runRule006,
}

func runRule006(pass *analysis.Pass) (any, error) {
	// Cette règle est réservée pour une utilisation future
	// Elle n'était pas implémentée dans l'ancien code
	return nil, nil
}
