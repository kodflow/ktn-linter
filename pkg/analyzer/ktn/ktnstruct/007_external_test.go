package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestStruct007 vérifie la détection des getters manquants pour structs non-DTO.
//
// Params:
//   - t: contexte de test
func TestStruct007(t *testing.T) {
	// 4 champs privés sans getters dans des structs non-DTO
	testhelper.TestGoodBad(t, ktnstruct.Analyzer007, "struct007", 4)
}
