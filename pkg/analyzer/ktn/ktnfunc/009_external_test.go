package ktnfunc_test

import (
	"testing"

	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc009(t *testing.T) {
	testhelper.TestGoodBad(t, ktnfunc.Analyzer009, "func009", 9)
}
