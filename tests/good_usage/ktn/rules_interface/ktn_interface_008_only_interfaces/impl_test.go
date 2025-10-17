package ktninterface008_test

import (
	"testing"

	ktninterface008 "github.com/kodflow/ktn-linter/tests/good_usage/rules_interface/ktn_interface_008_only_interfaces"
)

// TestNewService vérifie que NewService crée une instance valide.
//
// Params:
//   - t: contexte de test
func TestNewService(t *testing.T) {
	service := ktninterface008.NewService()
	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	// Vérifier que GetStatus retourne la valeur initiale
	status := service.GetStatus()
	if status != "ready" {
		t.Errorf("Expected status 'ready', got '%s'", status)
	}

	// Vérifier que Process ne retourne pas d'erreur
	err := service.Process("test data")
	if err != nil {
		t.Errorf("Process() returned unexpected error: %v", err)
	}
}

// TestNewRepository vérifie que NewRepository crée une instance valide.
//
// Params:
//   - t: contexte de test
func TestNewRepository(t *testing.T) {
	repo := ktninterface008.NewRepository()
	if repo == nil {
		t.Fatal("NewRepository() returned nil")
	}

	// Tester Save
	err := repo.Save("test data")
	if err != nil {
		t.Errorf("Save() returned unexpected error: %v", err)
	}

	// Tester Load après Save
	data, err := repo.Load("default")
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}
	if data != "test data" {
		t.Errorf("Expected 'test data', got '%s'", data)
	}
}
