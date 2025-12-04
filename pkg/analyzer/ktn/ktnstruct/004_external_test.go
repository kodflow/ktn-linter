package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct004(t *testing.T) {
	// good.go: 0 errors (1 struct), bad.go: 2 errors (3 structs - les 2 dernières sont en violation)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer004, "struct004", 2)
}
