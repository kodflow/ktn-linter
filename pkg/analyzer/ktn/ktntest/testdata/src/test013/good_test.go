package test013_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test004"
)

// TestProcessData teste ProcessData AVEC cas d'erreur.
// ProcessData retourne error → test avec cas d'erreur = OK.
//
// Params:
//   - t: contexte de test
func TestProcessData(t *testing.T) {
	// Test cas valide
	result, err := test004.ProcessData("hello")
	// Vérification pas d'erreur
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Vérification résultat
	if result != "processed:hello" {
		t.Logf("got %s", result)
	}

	// Test cas d'erreur - données vides
	_, err = test004.ProcessData("")
	// Vérification erreur attendue
	if err == nil {
		t.Log("expected error for empty data")
	}
}

// TestGetName teste GetName.
// GetName NE retourne PAS error → test simple = OK.
//
// Params:
//   - t: contexte de test
func TestGetName(t *testing.T) {
	// Test simple
	name := test004.GetName()
	// Vérification résultat
	if name != "test" {
		t.Logf("got %s", name)
	}
}

// TestGetCount teste GetCount.
// GetCount NE retourne PAS error → pas de cas d'erreur = OK.
//
// Params:
//   - t: contexte de test
func TestGetCount(t *testing.T) {
	// Test simple
	got := test004.GetCount()
	// Vérification résultat
	if got != 42 {
		t.Logf("got %d, want 42", got)
	}
}
