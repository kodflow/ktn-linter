package filesystem_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// TestInterfaces vérifie que le fichier interfaces.go existe et compile.
//
// Params:
//   - t: instance de test
func TestInterfaces(t *testing.T) {
	t.Log("Test de validation pour interfaces.go")
}

// TestOSFileSystemStat teste OSFileSystem.Stat.
//
// Params:
//   - t: instance de test
func TestOSFileSystemStat(t *testing.T) {
	fs := filesystem.NewOSFileSystem()

	// Tester avec un fichier qui devrait exister (ce fichier de test lui-même)
	_, err := fs.Stat("interfaces_test.go")
	if err != nil {
		t.Errorf("Expected file to exist, got error: %v", err)
	}

	// Tester avec un fichier qui n'existe pas
	_, err = fs.Stat("nonexistent_file_xyz_123.go")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
