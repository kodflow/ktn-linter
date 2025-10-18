//go:build test
// +build test

// Package ktninterface004 fournit les mocks pour les tests.
package ktninterface004

// MockService est un mock réutilisable de Service.
type MockService struct {
	ProcessFunc   func(data string) error
	GetStatusFunc func() string
}

// Process implémente Service.Process.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func (m *MockService) Process(data string) error {
	if m.ProcessFunc != nil {
		// Early return from function.
		return m.ProcessFunc(data)
	}
	// Early return from function.
	return nil
}

// GetStatus implémente Service.GetStatus.
//
// Returns:
//   - string: le statut actuel
func (m *MockService) GetStatus() string {
	if m.GetStatusFunc != nil {
		// Early return from function.
		return m.GetStatusFunc()
	}
	// Early return from function.
	return "mock"
}

// Vérification à la compilation que MockService implémente Service
var _ Service = (*MockService)(nil)

// MockRepository est un mock réutilisable de Repository.
type MockRepository struct {
	SaveFunc func(data string) error
	LoadFunc func(id string) (string, error)
}

// Save implémente Repository.Save.
//
// Params:
//   - data: les données à sauvegarder
//
// Returns:
//   - error: erreur si la sauvegarde échoue
func (m *MockRepository) Save(data string) error {
	if m.SaveFunc != nil {
		// Early return from function.
		return m.SaveFunc(data)
	}
	// Early return from function.
	return nil
}

// Load implémente Repository.Load.
//
// Params:
//   - id: l'identifiant des données
//
// Returns:
//   - string: les données chargées
//   - error: erreur si le chargement échoue
func (m *MockRepository) Load(id string) (string, error) {
	if m.LoadFunc != nil {
		// Early return from function.
		return m.LoadFunc(id)
	}
	// Early return from function.
	return "", nil
}

// Vérification à la compilation que MockRepository implémente Repository
var _ Repository = (*MockRepository)(nil)
