package goodembedded

import "io"

// Interfaces embarquées correctement documentées.

// Reader définit l'interface de lecture.
type Reader interface {
	// Read lit les données dans p.
	//
	// Params:
	//   - p: buffer de destination
	//
	// Returns:
	//   - n: nombre de bytes lus
	//   - err: erreur éventuelle
	Read(p []byte) (n int, err error)
}

// Writer définit l'interface d'écriture.
type Writer interface {
	// Write écrit les données de p.
	//
	// Params:
	//   - p: buffer source
	//
	// Returns:
	//   - n: nombre de bytes écrits
	//   - err: erreur éventuelle
	Write(p []byte) (n int, err error)
}

// ReadWriter combine Reader et Writer via embarquement d'interfaces.
// Cette interface hérite de toutes les méthodes de Reader et Writer.
type ReadWriter interface {
	Reader
	Writer
}

// ComplexEmbedded combine plusieurs interfaces standard io avec une méthode custom.
// Embarque io.Reader, io.Writer, et io.Closer pour réutiliser leurs contrats.
type ComplexEmbedded interface {
	io.Reader
	io.Writer
	io.Closer
	// CustomMethod méthode spécifique à cette interface.
	//
	// Returns:
	//   - error: erreur si l'opération échoue
	CustomMethod() error
}

// Processor interface qui combine io.Reader avec une méthode de traitement.
type Processor interface {
	io.Reader
	// Process traite les données lues.
	Process()
}
