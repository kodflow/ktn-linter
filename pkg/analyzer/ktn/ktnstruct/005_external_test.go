package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct005(t *testing.T) {
	// good.go: 0 errors (constructeur NewX présent), bad.go: 1 error (constructeur manquant)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer005, "struct005", 1)
}
