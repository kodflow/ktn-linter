package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestTest002(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "package naming convention",
			analyzer: "test002",
		},
		{
			name:     "validate analyzer consistency",
			analyzer: "test002",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1 erreur: mauvais nom de package
			testhelper.TestGoodBadWithFiles(t, ktntest.Analyzer002, tt.analyzer, "good_test.go", "bad_test.go", 1)
		})
	}
}
