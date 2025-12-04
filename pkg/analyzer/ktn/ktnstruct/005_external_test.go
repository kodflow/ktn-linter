package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct005(t *testing.T) {
	// good.go: 0 errors (champs exportés avant privés), bad.go: 5 errors (champs mélangés)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer005, "struct005", 5)
}
