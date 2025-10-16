package goodembedded

// Implementation implémente ComplexEmbedded avec toutes les méthodes requises.
type implementation struct {
	data []byte
}

// Read implémente io.Reader.
//
// Params:
//   - p: buffer de destination
//
// Returns:
//   - n: nombre de bytes lus
//   - err: erreur éventuelle
func (i *implementation) Read(p []byte) (n int, err error) {
	// Retourne 0 bytes lus et nil pour l'erreur
	return 0, nil
}

// Write implémente io.Writer.
//
// Params:
//   - p: buffer source
//
// Returns:
//   - n: nombre de bytes écrits
//   - err: erreur éventuelle
func (i *implementation) Write(p []byte) (n int, err error) {
	// Retourne le nombre de bytes écrits et nil pour l'erreur
	return len(p), nil
}

// Close implémente io.Closer.
//
// Returns:
//   - error: erreur si la fermeture échoue
func (i *implementation) Close() error {
	// Retourne nil car la fermeture est réussie
	return nil
}

// CustomMethod implémente la méthode spécifique de ComplexEmbedded.
//
// Returns:
//   - error: erreur si l'opération échoue
func (i *implementation) CustomMethod() error {
	// Retourne nil car l'opération est terminée avec succès
	return nil
}

// NewImplementation crée une nouvelle instance de Implementation.
//
// Returns:
//   - ComplexEmbedded: l'interface implémentée
func NewImplementation() ComplexEmbedded {
	// Retourne une nouvelle instance de l'implémentation
	return &implementation{}
}
