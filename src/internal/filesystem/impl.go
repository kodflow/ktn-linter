package filesystem

import "os"

// osFileSystem implémente FileSystem en utilisant le vrai système de fichiers.
type osFileSystem struct{}

// NewFileSystem crée une nouvelle instance de FileSystem utilisant l'implémentation par défaut.
//
// Returns:
//   - FileSystem: l'implémentation du système de fichiers
func NewFileSystem() FileSystem {
	// Retourne une nouvelle instance du système de fichiers OS
	return &osFileSystem{}
}

// NewOSFileSystem crée une nouvelle instance de FileSystem utilisant os.Stat.
//
// Returns:
//   - FileSystem: l'implémentation réelle du système de fichiers
func NewOSFileSystem() FileSystem {
	// Retourne une nouvelle instance du système de fichiers OS
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
	// Délègue à la fonction os.Stat standard
	return os.Stat(name)
}
