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
