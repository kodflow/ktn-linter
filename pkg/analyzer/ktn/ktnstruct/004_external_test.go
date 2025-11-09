package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct004(t *testing.T) {
	// good.go: 0 errors (documentation complète), bad.go: 4 errors (documentation manquante/insuffisante)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer004, "struct004", 4)
}
