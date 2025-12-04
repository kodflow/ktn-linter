package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct004(t *testing.T) {
	// good.go: 0 errors (documentation complète), bad.go: 1 error (documentation insuffisante)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer004, "struct004", 1)
}
