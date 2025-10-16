package filesystem_test

import (
	"os"
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
)

// TestNewFileSystem vérifie que NewFileSystem crée une instance valide.
//
// Params:
//   - t: contexte de test
func TestNewFileSystem(t *testing.T) {
	fs := filesystem.NewFileSystem()
	if fs == nil {
		t.Fatal("NewFileSystem() returned nil")
	}

	// Tester avec un fichier qui existe (ce fichier de test lui-même)
	_, err := fs.Stat("impl_test.go")
	if err != nil {
		t.Errorf("Stat() failed for existing file: %v", err)
	}

	// Tester avec un fichier qui n'existe pas
	_, err = fs.Stat("nonexistent_file_xyz.go")
	if err == nil {
		t.Error("Stat() should fail for nonexistent file")
	}
}

// TestNewOSFileSystem vérifie que NewOSFileSystem crée une instance valide.
//
// Params:
//   - t: contexte de test
func TestNewOSFileSystem(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	if fs == nil {
		t.Fatal("NewOSFileSystem() returned nil")
	}

	// Vérifier que Stat fonctionne correctement
	info, err := fs.Stat("impl_test.go")
	if err != nil {
		t.Fatalf("Stat() failed: %v", err)
	}

	if info.IsDir() {
		t.Error("impl_test.go should be a file, not a directory")
	}

	if info.Size() == 0 {
		t.Error("impl_test.go should have non-zero size")
	}
}

// TestOSFileSystemStatError vérifie le comportement en cas d'erreur.
//
// Params:
//   - t: contexte de test
func TestOSFileSystemStatError(t *testing.T) {
	fs := filesystem.NewOSFileSystem()

	_, err := fs.Stat("this_file_absolutely_does_not_exist_xyz123.go")
	if err == nil {
		t.Error("Stat() should return error for nonexistent file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected os.IsNotExist error, got: %v", err)
	}
}
