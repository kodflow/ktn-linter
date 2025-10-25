package test004_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test004"
)

// TestValidateInputWithErrors teste avec cas d'erreur (BIEN).
//
// Params:
//   - t: contexte de test
func TestValidateInputWithErrors(t *testing.T) {
	// Test cas valide
	err := test004.ValidateInput(10)
	// Vérification pas d'erreur
	if err != nil {
		t.Errorf("ValidateInput(10) devrait réussir, got error: %v", err)
	}

	// Test cas invalide - cas d'erreur
	err = test004.ValidateInput(0)
	// Vérification erreur attendue
	if err == nil {
		t.Error("ValidateInput(0) devrait échouer")
	}

	// Test cas invalide négatif - cas d'erreur
	err = test004.ValidateInput(-5)
	// Vérification erreur attendue
	if err == nil {
		t.Error("ValidateInput(-5) devrait échouer")
	}
}

// TestProcessDataWithErrors teste avec cas d'erreur (BIEN).
//
// Params:
//   - t: contexte de test
func TestProcessDataWithErrors(t *testing.T) {
	// Test cas valide
	result, err := test004.ProcessData("test")
	// Vérification pas d'erreur
	if err != nil {
		t.Fatalf("ProcessData should not fail on valid data: %v", err)
	}
	// Vérification résultat
	if result != "processed: test" {
		t.Errorf("got %q, want %q", result, "processed: test")
	}

	// Test cas d'erreur - données vides
	_, err = test004.ProcessData("")
	// Vérification erreur attendue
	if err == nil {
		t.Error("ProcessData should fail on empty data")
	}
}
