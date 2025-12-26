package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest003(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "public function test coverage",
			analyzer: "test003",
		},
		{
			name:     "validate test completeness",
			analyzer: "test003",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 2 erreurs: fonctions Multiply et Divide sans test
			testhelper.TestGoodBadPackage(t, ktntest.Analyzer003, tt.analyzer, 2)
		})
	}
}
