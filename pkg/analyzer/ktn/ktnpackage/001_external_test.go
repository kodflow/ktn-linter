package ktnpackage_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnpackage"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAnalyzer001 teste KTN-PACKAGE-001.
//
// Params:
//   - t: contexte de test
func TestAnalyzer001(t *testing.T) {
	// Expected 2 errors in bad package:
	// - no_comment.go (pas de commentaire de fichier)
	// - empty_comment.go (commentaire vide)
	testhelper.TestGoodBadPackage(t, ktnpackage.Analyzer001, "package001", 2)
}
