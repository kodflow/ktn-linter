package test011 // want "KTN-TEST-011: le fichier 'service_external_test.go' doit utiliser 'package service_test'"

import "testing"

// TestPublicService - ERREUR: package test011 dans _external_test.go
func TestPublicService(t *testing.T) {
	result := PublicService()
	if result != "service" {
		t.Errorf("expected 'service', got '%s'", result)
	}
}
