//go:build test
// +build test

package filesystem

import "os"

// MockFileSystem est un mock de FileSystem pour les tests.
type MockFileSystem struct {
	StatFunc func(name string) (os.FileInfo, error)
}

// Stat implémente FileSystem.Stat.
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, nil
}

// Vérification à la compilation que MockFileSystem implémente FileSystem
var _ FileSystem = (*MockFileSystem)(nil)
