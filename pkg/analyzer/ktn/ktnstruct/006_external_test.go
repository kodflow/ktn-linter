package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct006 vérifie la détection des champs privés tagués dans les DTOs.
//
// Params:
//   - t: contexte de test
func TestStruct006(t *testing.T) {
	// 4 champs privés avec tags dans des DTOs
	testhelper.TestGoodBad(t, ktnstruct.Analyzer006, "struct006", 4)
}
