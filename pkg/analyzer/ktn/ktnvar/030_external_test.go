package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar030 tests the detection of slice cloning patterns.
// Expected errors in bad.go:
// - badMakeCopyInt: make+copy pattern (1)
// - badMakeCopyString: make+copy pattern (1)
// - badAppendNilInt: append([]T(nil), s...) pattern (1)
// - badAppendNilString: append([]T(nil), s...) pattern (1)
// - badMakeCopyBytes: make+copy pattern (1)
// - badAppendNilBytes: append([]T(nil), s...) pattern (1)
// Total: 6 errors
func TestVar030(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Slice cloning patterns that should use slices.Clone",
			analyzer:       ktnvar.Analyzer030,
			testdataDir:    "var030",
			expectedErrors: 6,
		},
		{
			name:           "Valid slices.Clone usage and non-clone patterns",
			analyzer:       ktnvar.Analyzer030,
			testdataDir:    "var030",
			expectedErrors: 6,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
