package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar023 vérifie la détection de math/rand dans un contexte sécurité.
// Erreurs attendues dans bad.go:
// - badGenerateKey: appel rand.Intn dans fonction avec "Key" (1)
// - badCreateToken: appel rand.Int63 dans fonction avec "Token" (1)
// - badSecretKey: variable globale avec "Secret" et appel rand.Uint64 (1)
// - badGenerateSecret: appel rand.Intn dans fonction avec "Secret" (1)
// - secretValue: variable locale avec "secret" et rand.Intn (1)
// - badNormalFunction: tokenValue variable locale (déclenche checkLocalVarDecl) (1)
// Total: 6 erreurs
func TestVar023(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "math/rand in security context",
			analyzer:       ktnvar.Analyzer023,
			testdataDir:    "var023",
			expectedErrors: 7,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}
