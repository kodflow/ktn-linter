package ktntest_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestTest005 teste la règle KTN-TEST-007
func TestTest005(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
	}{
		{
			name:     "Skip statement validation",
			analyzer: "test005",
		},
		{
			name:     "verify Skip usage compliance",
			analyzer: "test005",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdata := analysistest.TestData()
			// Test bad_test.go contient les cas d'erreur Skip/Skipf/SkipNow invalid
			analysistest.Run(t, testdata, ktntest.Analyzer005, tt.analyzer)
		})
	}
}
