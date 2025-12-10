package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest013(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "error case coverage validation",
			analyzer: "test013",
		},
		{
			name:     "verify error path testing",
			analyzer: "test013",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 2 erreurs: tests de fonctions retournant error sans couverture d'erreur
			// - TestParseConfig: ParseConfig retourne error mais test sans cas d'erreur
			// - TestValidateInput: ValidateInput retourne error mais test sans cas d'erreur
			testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer013, tt.analyzer, "good_test.go", "bad_test.go", 2)
		})
	}
}
