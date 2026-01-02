package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar022 vérifie la détection des pointeurs vers interfaces.
// Erreurs attendues dans bad.go:
// - badProcess: param *io.Reader (1)
// - badHandle: param *io.Writer (1)
// - badHandler: var *interface{} (1)
// - badAny: var *any (1)
// - BadService.reader: field *io.Reader (1)
// Total: 5 erreurs
func TestVar022(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Pointer to interface detection",
			analyzer:       ktnvar.Analyzer022,
			testdataDir:    "var022",
			expectedErrors: 5,
		},
		{
			name:           "Valid interface usage without pointer",
			analyzer:       ktnvar.Analyzer022,
			testdataDir:    "var022",
			expectedErrors: 5,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
