package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct002(t *testing.T) {
	// good.go: 0 errors (interface complète), bad.go: 2 errors (1 sans interface + 1 interface incomplète)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer002, "struct002", 2)
}
