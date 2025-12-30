// External tests for ktnapi Analyzer001.
package ktnapi_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestAPI001 tests the Analyzer001 for external concrete type dependencies.
func TestAPI001(t *testing.T) {
	tests := []struct {
		name             string
		analyzer         *analysis.Analyzer
		testDataDir      string
		expectedBadCount int
	}{
		{
			name:             "external concrete type with method calls",
			analyzer:         ktnapi.Analyzer001,
			testDataDir:      "api001",
			expectedBadCount: 6,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// good.go: 0 errors, bad.go: 6 errors for external concrete types
			// - badHTTPClientWithOneMethod: 1 error
			// - badHTTPClientWithMultipleMethods: 1 error
			// - badFileWithMethods: 1 error
			// - badMultipleParams: 2 errors (client + file)
			// - badAliasExternalType: 1 error (alias to http.Client)
			// Note: bytes.Buffer and strings.Builder are now in allowed list
			testhelper.TestGoodBad(t, tt.analyzer, tt.testDataDir, tt.expectedBadCount)
		})
	}
}
