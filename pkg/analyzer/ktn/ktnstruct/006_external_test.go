package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct006(t *testing.T) {
	// good.go: 0 errors (getters idiomatiques sans Get), bad.go: 3 errors (getters avec préfixe Get)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer006, "struct006", 3)
}
