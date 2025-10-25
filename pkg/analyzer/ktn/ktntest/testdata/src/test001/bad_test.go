package test001

import "testing"

// TestBadPackageName utilise le mauvais package (test001 au lieu de test001_test).
// Ceci viole KTN-TEST-001.
//
// Params:
//   - t: contexte de test
func TestBadPackageName(t *testing.T) {
	// Test avec package interne (PAS BIEN)
	err := CheckValue(10)
	// Vérification pas d'erreur
	if err != nil {
		t.Errorf("CheckValue devrait réussir: %v", err)
	}

	// Test cas erreur
	err = CheckValue(0)
	// Vérification erreur
	if err == nil {
		t.Error("CheckValue(0) devrait échouer")
	}
}

// TestInternalPackage utilise le package interne (PAS BIEN).
// Devrait utiliser test001_test.
//
// Params:
//   - t: contexte de test
func TestInternalPackage(t *testing.T) {
	const EXPECTED string = ">hello<"
	// Test avec package interne
	result, err := Format("hello")
	// Vérification pas d'erreur
	if err != nil {
		t.Fatalf("Format devrait réussir: %v", err)
	}
	// Vérification résultat
	if result != EXPECTED {
		t.Errorf("got %q, want %q", result, EXPECTED)
	}

	// Test cas erreur
	_, err = Format("")
	// Vérification erreur
	if err == nil {
		t.Error("Format should fail on empty")
	}
}
