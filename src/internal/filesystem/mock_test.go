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
		// Retourne les informations du fichier simulé
		return &mockFileInfo{name: name}, nil
	}
	// Fichier non trouvé dans le mock
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
	// Retourne le nom stocké
	return m.name
}

// Size retourne la taille (toujours 0 pour le mock).
//
// Returns:
//   - int64: la taille du fichier
func (m *mockFileInfo) Size() int64 {
	// Taille simulée fixe
	return 0
}

// Mode retourne le mode (toujours 0644 pour le mock).
//
// Returns:
//   - os.FileMode: le mode du fichier
func (m *mockFileInfo) Mode() os.FileMode {
	// Mode simulé fixe (rw-r--r--)
	return 0644
}

// ModTime retourne le temps de modification (epoch pour le mock).
//
// Returns:
//   - time.Time: le temps de modification
func (m *mockFileInfo) ModTime() time.Time {
	// Retourne l'epoch Unix comme temps simulé
	return time.Unix(0, 0)
}

// IsDir retourne false (pas un répertoire pour le mock).
//
// Returns:
//   - bool: false car c'est un fichier
func (m *mockFileInfo) IsDir() bool {
	// Simule toujours un fichier, jamais un répertoire
	return false
}

// Sys retourne nil (pas d'infos système pour le mock).
//
// Returns:
//   - interface{}: nil
func (m *mockFileInfo) Sys() interface{} {
	// Pas d'informations système pour le mock
	return nil
}
