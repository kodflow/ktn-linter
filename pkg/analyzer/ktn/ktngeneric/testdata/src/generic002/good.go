package generic002

import "io"

// readSome lit des donnees depuis un reader.
// Correct: utilise l'interface directement.
//
// Params:
//   - r: reader source
//
// Returns:
//   - []byte: donnees lues
//   - error: erreur eventuelle
func readSome(r io.Reader) ([]byte, error) {
	// Creer un buffer
	buf := make([]byte, 1024)
	// Lire les donnees
	n, err := r.Read(buf)
	// Verifier l'erreur
	if err != nil {
		// Retourner l'erreur
		return nil, err
	}
	// Retourner les donnees
	return buf[:n], nil
}

// writeAll ecrit des donnees dans un writer.
// Correct: utilise l'interface directement.
//
// Params:
//   - w: writer destination
//   - data: donnees a ecrire
//
// Returns:
//   - error: erreur eventuelle
func writeAll(w io.Writer, data []byte) error {
	// Ecrire les donnees
	_, err := w.Write(data)
	// Retourner le resultat
	return err
}

// closeResource ferme une ressource.
// Correct: utilise l'interface directement.
//
// Params:
//   - c: closer a fermer
//
// Returns:
//   - error: erreur eventuelle
func closeResource(c io.Closer) error {
	// Fermer la ressource
	return c.Close()
}

// passThrough retourne le meme reader.
// Correct: generique justifie pour la preservation du type.
//
// Params:
//   - r: reader source
//
// Returns:
//   - T: le meme reader avec son type preserve
func passThrough[T io.Reader](r T) T {
	// Retourner le reader tel quel
	return r
}

// wrapWriter encapsule un writer et retourne le meme type.
// Correct: generique justifie pour la preservation du type.
//
// Params:
//   - w: writer a encapsuler
//
// Returns:
//   - T: le meme writer avec son type preserve
func wrapWriter[T io.Writer](w T) T {
	// Retourner le writer tel quel
	return w
}

// processAndReturn traite un closer et le retourne.
// Correct: generique justifie pour la preservation du type.
//
// Params:
//   - c: closer a traiter
//
// Returns:
//   - T: le meme closer avec son type preserve
//   - error: erreur eventuelle
func processAndReturn[T io.Closer](c T) (T, error) {
	// Retourner le closer et nil
	return c, nil
}

// mapSlice applique une fonction a chaque element.
// Correct: utilise any sans contrainte interface specifique.
//
// Params:
//   - s: slice source
//   - f: fonction de transformation
//
// Returns:
//   - []U: slice transforme
func mapSlice[T, U any](s []T, f func(T) U) []U {
	// Creer le slice resultat
	result := make([]U, len(s))
	// Parcourir et transformer
	for i, x := range s {
		// Appliquer la fonction
		result[i] = f(x)
	}
	// Retourner le resultat
	return result
}

// filterSlice filtre un slice.
// Correct: utilise any sans contrainte interface specifique.
//
// Params:
//   - s: slice source
//   - pred: predicat de filtre
//
// Returns:
//   - []T: slice filtre
func filterSlice[T any](s []T, pred func(T) bool) []T {
	// Creer le slice resultat
	var result []T
	// Parcourir et filtrer
	for _, x := range s {
		// Tester le predicat
		if pred(x) {
			// Ajouter au resultat
			result = append(result, x)
		}
	}
	// Retourner le resultat
	return result
}
