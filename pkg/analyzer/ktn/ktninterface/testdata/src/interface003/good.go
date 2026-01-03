// Package interface003 contains test cases for KTN-INTERFACE-003.
package interface003

// Reader follows the -er convention for single-method interface.
type Reader interface {
	Read()
}

// Writer follows the -er convention for single-method interface.
type Writer interface {
	Write(data []byte)
}

// Stringer follows the -er convention for single-method interface.
type Stringer interface {
	String() string
}

// Handler follows the -or convention for single-method interface.
type Handler interface {
	Handle()
}

// ReadWriter is a multi-method interface (no -er check needed).
type ReadWriter interface {
	Read()
	Write()
}

// Empty interface (no -er check needed).
type Empty interface{}
