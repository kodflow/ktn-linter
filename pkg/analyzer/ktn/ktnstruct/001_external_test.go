package ktnstruct_test

import (
	"testing"

	ktnstruct "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestStruct001(t *testing.T) {
	// good.go: 0 errors (interface complète), bad.go: 1 error (struct sans interface complète)
	testhelper.TestGoodBad(t, ktnstruct.Analyzer001, "struct001", 1)
}
