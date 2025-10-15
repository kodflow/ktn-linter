package filesystem_test

import (
	"errors"
	"os"
	"time"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// mockFileSystem est un mock du FileSystem pour les tests.
type mockFileSystem struct {
	files map[string]bool
}

// NewMockFileSystem crée un nouveau mock FileSystem.
//
// Params:
//   - files: map des chemins de fichiers (true = existe, false = n'existe pas)
//
// Returns:
//   - filesystem.FileSystem: le mock FileSystem
func NewMockFileSystem(files map[string]bool) filesystem.FileSystem {
	// Retourne une instance de mockFileSystem avec la map de fichiers fournie
	return &mockFileSystem{files: files}
}

// Stat simule os.Stat en utilisant la map de fichiers.
//
// Params:
//   - name: le chemin du fichier
//
// Returns:
//   - os.FileInfo: informations simulées si le fichier existe
//   - error: erreur si le fichier n'existe pas
func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	if exists, ok := m.files[name]; ok && exists {
		return &mockFileInfo{name: name}, nil
	}
	return nil, errors.New("file does not exist")
}

// mockFileInfo implémente os.FileInfo pour le mock.
type mockFileInfo struct {
	name string
}

// Name retourne le nom du fichier.
//
// Returns:
//   - string: le nom du fichier
func (m *mockFileInfo) Name() string {
	return m.name
}

// Size retourne la taille (toujours 0 pour le mock).
//
// Returns:
//   - int64: la taille du fichier
func (m *mockFileInfo) Size() int64 {
	return 0
}

// Mode retourne le mode (toujours 0644 pour le mock).
//
// Returns:
//   - os.FileMode: le mode du fichier
func (m *mockFileInfo) Mode() os.FileMode {
	return 0644
}

// ModTime retourne le temps de modification (epoch pour le mock).
//
// Returns:
//   - time.Time: le temps de modification
func (m *mockFileInfo) ModTime() time.Time {
	return time.Unix(0, 0)
}

// IsDir retourne false (pas un répertoire pour le mock).
//
// Returns:
//   - bool: false car c'est un fichier
func (m *mockFileInfo) IsDir() bool {
	return false
}

// Sys retourne nil (pas d'infos système pour le mock).
//
// Returns:
//   - interface{}: nil
func (m *mockFileInfo) Sys() interface{} {
	return nil
}
