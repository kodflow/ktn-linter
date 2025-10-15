package goodembedded

// Implementation implémente ComplexEmbedded avec toutes les méthodes requises.
type Implementation struct {
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
func (i *Implementation) Read(p []byte) (n int, err error) {
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
func (i *Implementation) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Close implémente io.Closer.
//
// Returns:
//   - error: erreur si la fermeture échoue
func (i *Implementation) Close() error {
	return nil
}

// CustomMethod implémente la méthode spécifique de ComplexEmbedded.
//
// Returns:
//   - error: erreur si l'opération échoue
func (i *Implementation) CustomMethod() error {
	return nil
}

// NewImplementation crée une nouvelle instance de Implementation.
//
// Returns:
//   - ComplexEmbedded: l'interface implémentée
func NewImplementation() ComplexEmbedded {
	return &Implementation{}
}
