package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest004(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "public function test coverage",
			analyzer: "test004",
		},
		{
			name:     "validate test completeness",
			analyzer: "test004",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 2 erreurs: fonctions Multiply et Divide sans test
			testhelper.TestGoodBadPackage(t, ktntest.Analyzer004, tt.analyzer, 2)
		})
	}
}
