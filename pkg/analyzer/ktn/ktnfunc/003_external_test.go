package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestFunc003 teste KTN-FUNC-003.
// Erreurs attendues dans bad.go:
// - badCheckPositive: else après return
// - badProcessValue: else après return
// - badFindMax: else après return
// - badLoopExample: else après continue
// - badSwitchExample: else après break
// - badValidateInput: else après return
// - badPanicExample: else après panic
// - badElseIfExample: else if après return + else après return (2 erreurs)
// Total: 9 erreurs
//
// Params:
//   - t: contexte de test
func TestFunc003(t *testing.T) {
	testhelper.TestGoodBad(t, ktnfunc.Analyzer003, "func003", 9)
}
