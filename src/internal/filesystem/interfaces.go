package filesystem

import "os"

// FileSystem définit l'interface pour les opérations sur le système de fichiers.
// Cette interface permet le mocking dans les tests.
//
// Constructeur:
//   - NewOSFileSystem(): crée une implémentation utilisant le système de fichiers réel
type FileSystem interface {
	// Stat retourne les informations sur un fichier.
	//
	// Params:
	//   - name: le chemin du fichier
	//
	// Returns:
	//   - os.FileInfo: les informations du fichier
	//   - error: erreur si le fichier n'existe pas
	Stat(name string) (os.FileInfo, error)
}

// OSFileSystem implémente FileSystem en utilisant le vrai système de fichiers.
type osFileSystem struct{}

// NewFileSystem crée une nouvelle instance de FileSystem utilisant l'implémentation par défaut.
//
// Returns:
//   - FileSystem: l'implémentation du système de fichiers
func NewFileSystem() FileSystem {
	// Retourne une instance de osFileSystem qui utilise os.Stat
	return &osFileSystem{}
}

// NewOSFileSystem crée une nouvelle instance de FileSystem utilisant os.Stat.
//
// Returns:
//   - FileSystem: l'implémentation réelle du système de fichiers
func NewOSFileSystem() FileSystem {
	// Retourne une instance de osFileSystem qui utilise os.Stat
	return &osFileSystem{}
}

// Stat appelle os.Stat pour obtenir les informations d'un fichier.
//
// Params:
//   - name: le chemin du fichier
//
// Returns:
//   - os.FileInfo: les informations du fichier
//   - error: erreur si le fichier n'existe pas
func (fs *osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
