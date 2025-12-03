package test004_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test004"
)

// TestParseConfig teste ParseConfig SANS cas d'erreur.
// ParseConfig retourne error → devrait avoir des cas d'erreur.
//
// Params:
//   - t: contexte de test
func TestParseConfig(t *testing.T) { // want "KTN-TEST-004: le test 'TestParseConfig' teste une fonction qui retourne error, il devrait couvrir les cas d'erreur"
	// Test uniquement le cas valide
	result, _ := test004.ParseConfig("config.yaml")
	// Vérification résultat
	if result == "" {
		t.Log("empty result")
	}
}

// TestValidateInput teste ValidateInput SANS cas d'erreur.
// ValidateInput retourne error → devrait avoir des cas d'erreur.
//
// Params:
//   - t: contexte de test
func TestValidateInput(t *testing.T) { // want "KTN-TEST-004: le test 'TestValidateInput' teste une fonction qui retourne error, il devrait couvrir les cas d'erreur"
	// Test uniquement le cas valide
	_ = test004.ValidateInput(10)
}

// TestGetVersion teste GetVersion.
// GetVersion NE retourne PAS error → pas d'erreur attendue.
//
// Params:
//   - t: contexte de test
func TestGetVersion(t *testing.T) {
	// Test simple - pas d'erreur attendue car GetVersion ne retourne pas error
	got := test004.GetVersion()
	// Vérification résultat
	if got != "1.0.0" {
		t.Logf("got %s, want 1.0.0", got)
	}
}
