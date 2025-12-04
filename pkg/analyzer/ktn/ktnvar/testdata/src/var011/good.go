// Good examples for the var011 test case.
package var011

import (
	"bytes"
	"strings"
)

const (
	// GROW_SIZE_LARGE is large grow size
	GROW_SIZE_LARGE int = 400
	// LOOP_COUNT_LARGE is large loop count
	LOOP_COUNT_LARGE int = 100
	// LOOP_COUNT_SMALL is small loop count
	LOOP_COUNT_SMALL int = 10
	// GROW_SIZE_SMALL is small grow size
	GROW_SIZE_SMALL int = 50
)

// goodStringsBuilderTypeDecl uses type declaration (not composite literal).
//
// Returns:
//   - string: concatenated result
func goodStringsBuilderTypeDecl() string {
	// Good: var declaration without composite literal is allowed
	var sb strings.Builder
	sb.Grow(GROW_SIZE_LARGE)

	// Iteration over data to append
	for i := 0; i < LOOP_COUNT_LARGE; i++ {
		sb.WriteString("item")
	}

	// Return the result
	return sb.String()
}

// goodBytesBufferTypeDecl uses type declaration (not composite literal).
//
// Returns:
//   - []byte: concatenated result
func goodBytesBufferTypeDecl() []byte {
	// Good: var declaration without composite literal is allowed
	var buf bytes.Buffer
	buf.Grow(GROW_SIZE_LARGE)

	// Iteration over data to append
	for i := 0; i < LOOP_COUNT_LARGE; i++ {
		buf.WriteString("item")
	}

	// Return the result
	return buf.Bytes()
}

// goodBuilderPointer uses a pointer to strings.Builder (allowed).
//
// Returns:
//   - string: concatenated result
func goodBuilderPointer() string {
	// Good: pointer type is allowed (different use case)
	sb := &strings.Builder{}

	// Iteration over data to append
	for i := 0; i < LOOP_COUNT_SMALL; i++ {
		sb.WriteString("x")
	}

	// Return the result
	return sb.String()
}

// goodNoLoopTypeDecl uses type declaration without loop (allowed).
//
// Returns:
//   - string: concatenated result
func goodNoLoopTypeDecl() string {
	// Good: var declaration without composite literal
	var sb strings.Builder
	sb.WriteString("single")

	// Return the result
	return sb.String()
}

// goodShortFormNew uses new() to create Builder (allowed).
//
// Returns:
//   - string: concatenated result
func goodShortFormNew() string {
	// Good: using new() instead of composite literal
	sb := new(strings.Builder)
	sb.Grow(GROW_SIZE_SMALL)

	// Iteration over data to append
	for i := 0; i < GROW_SIZE_SMALL; i++ {
		sb.WriteString("x")
	}

	// Return the result
	return sb.String()
}

// init utilise les fonctions privées
func init() {
	// Appel de goodStringsBuilderTypeDecl
	goodStringsBuilderTypeDecl()
	// Appel de goodBytesBufferTypeDecl
	goodBytesBufferTypeDecl()
	// Appel de goodBuilderPointer
	goodBuilderPointer()
	// Appel de goodNoLoopTypeDecl
	goodNoLoopTypeDecl()
	// Appel de goodShortFormNew
	goodShortFormNew()
}
