package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest004(t *testing.T) {
	// 2 erreurs: tests de fonctions retournant error sans couverture d'erreur
	// - TestParseConfig: ParseConfig retourne error mais test sans cas d'erreur
	// - TestValidateInput: ValidateInput retourne error mais test sans cas d'erreur
	testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer004, "test004", "good_test.go", "bad_test.go", 2)
}
