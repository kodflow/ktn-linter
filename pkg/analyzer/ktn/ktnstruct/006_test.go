package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct006(t *testing.T) {
	// good.go: 0 errors (champs privés + getters), bad.go: 5 errors (champs publics/getters manquants)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer006, "struct006", 5)
}
