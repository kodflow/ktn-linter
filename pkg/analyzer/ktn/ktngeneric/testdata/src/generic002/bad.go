// Package generic002 contient les tests pour KTN-GENERIC-002.
// Ce fichier contient des exemples de code non-conforme.
package generic002

import (
	"fmt"
	"io"
)

// BadWriteAll ecrit des donnees dans un writer.
// Erreur: utilise un generique inutile sur une interface.
//
// Params:
//   - w: writer destination
//   - data: donnees a ecrire
//
// Returns:
//   - error: erreur eventuelle
func BadWriteAll[T io.Writer](w T, data []byte) error { // want "KTN-GENERIC-002"
	// Ecrire les donnees
	_, err := w.Write(data)
	// Retourner le resultat
	return err
}

// BadCloser ferme une ressource.
// Erreur: utilise un generique inutile sur une interface.
//
// Params:
//   - c: closer a fermer
//
// Returns:
//   - error: erreur eventuelle
func BadCloser[T io.Closer](c T) error { // want "KTN-GENERIC-002"
	// Fermer la ressource
	return c.Close()
}

// BadSeeker repositionne un flux.
// Erreur: utilise un generique inutile sur une interface.
//
// Params:
//   - s: seeker a repositionner
//   - offset: decalage
//   - whence: origine
//
// Returns:
//   - int64: nouvelle position
//   - error: erreur eventuelle
func BadSeeker[T io.Seeker](s T, offset int64, whence int) (int64, error) { // want "KTN-GENERIC-002"
	// Repositionner le flux
	return s.Seek(offset, whence)
}

// BadStringer convertit en string.
// Erreur: utilise un generique inutile sur fmt.Stringer.
//
// Params:
//   - s: valeur a convertir
//
// Returns:
//   - string: representation textuelle
func BadStringer[T fmt.Stringer](s T) string { // want "KTN-GENERIC-002"
	// Convertir en string
	return s.String()
}
