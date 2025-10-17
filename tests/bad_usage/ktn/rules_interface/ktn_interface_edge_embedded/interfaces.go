package badembedded

import "io"

// Violations avec interfaces embarquées

// reader interface privée (devrait être publique)
type reader interface {
	Read(p []byte) (n int, err error)
}

// writer sans commentaire
type writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriter interface qui embarque sans documentation appropriée
type ReadWriter interface {
	reader
	writer
}

// ComplexEmbedded interface sans doc claire sur l'embarquement
type ComplexEmbedded interface {
	io.Reader
	io.Writer
	io.Closer
	CustomMethod() error
}

// badNaming mauvais nommage pour interface
type badNaming interface {
	io.Reader
	Process()
}
