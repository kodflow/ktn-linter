package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest006(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "orphan test file detection",
			analyzer: "test006",
		},
		{
			name:     "verify test-source pairing",
			analyzer: "test006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1 erreur: fichier orphan_test.go sans fichier orphan.go
			testhelper.TestGoodBadPackage(t, ktntest.Analyzer006, tt.analyzer, 1)
		})
	}
}
