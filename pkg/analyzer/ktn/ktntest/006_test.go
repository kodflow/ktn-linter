package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest006(t *testing.T) {
	// 1 erreur: fichier orphan_test.go sans fichier orphan.go
	testhelper.TestGoodBadPackage(t, ktntest.Analyzer006, "test006", 1)
}
