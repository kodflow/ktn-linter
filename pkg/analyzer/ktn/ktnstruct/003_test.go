package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct003(t *testing.T) {
	// good.go: 0 errors (champs exportés avant privés), bad.go: 5 errors (champs mélangés)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer003, "struct003", 5)
}
