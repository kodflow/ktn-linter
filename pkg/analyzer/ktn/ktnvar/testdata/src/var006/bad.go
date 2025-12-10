// Bad examples for the var007 test case.
package var006

import (
	"bytes"
	"strings"
)

// Constantes pour les tests
const (
	BadLoopCountLarge int = 100
	BadLoopCountSmall int = 50
)

// badStringsBuilderNoGrow creates a strings.Builder without Grow.
//
// Returns:
//   - string: concatenated result
func badStringsBuilderNoGrow() string {
	// Bad: composite literal without Grow() call
	sb := strings.Builder{}

	// Iteration over data to append
	for i := 0; i < BadLoopCountLarge; i++ {
		sb.WriteString("item")
	}

	// Return the result
	return sb.String()
}

// badBytesBufferNoGrow creates a bytes.Buffer without Grow.
//
// Returns:
//   - []byte: concatenated result
func badBytesBufferNoGrow() []byte {
	// Bad: composite literal without Grow() call
	buf := bytes.Buffer{}

	// Iteration over data to append
	for i := 0; i < BadLoopCountLarge; i++ {
		buf.WriteString("item")
	}

	// Return the result
	return buf.Bytes()
}

// badShortFormBuilder uses short form without Grow.
//
// Returns:
//   - string: concatenated result
func badShortFormBuilder() string {
	// Bad: short declaration without Grow
	sb := strings.Builder{}

	// Iteration over data to append
	for i := 0; i < BadLoopCountSmall; i++ {
		sb.WriteString("x")
	}

	// Return the result
	return sb.String()
}

// badShortFormBuffer uses short form bytes.Buffer without Grow.
//
// Returns:
//   - []byte: concatenated result
func badShortFormBuffer() []byte {
	// Bad: short declaration without Grow
	buf := bytes.Buffer{}

	// Iteration over data to append
	for i := 0; i < BadLoopCountSmall; i++ {
		buf.Write([]byte("x"))
	}

	// Return the result
	return buf.Bytes()
}

// init utilise les fonctions privées
func init() {
	// Appel de badStringsBuilderNoGrow
	badStringsBuilderNoGrow()
	// Appel de badBytesBufferNoGrow
	badBytesBufferNoGrow()
	// Appel de badShortFormBuilder
	badShortFormBuilder()
	// Appel de badShortFormBuffer
	badShortFormBuffer()
}
