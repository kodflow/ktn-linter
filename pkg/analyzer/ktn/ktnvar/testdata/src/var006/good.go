// Package var006 provides good test cases.
package var006

import (
	"bytes"
	"strings"
)

const (
	// GrowSizeLarge is large grow size
	GrowSizeLarge int = 400
	// LoopCountLarge is large loop count
	LoopCountLarge int = 100
	// LoopCountSmall is small loop count
	LoopCountSmall int = 10
	// GrowSizeSmall is small grow size
	GrowSizeSmall int = 50
)

// init demonstrates proper Builder usage
func init() {
	// Good: var declaration without composite literal is allowed
	var sb strings.Builder
	sb.Grow(GrowSizeLarge)
	// Iteration over data to append
	for i := range LoopCountLarge {
		sb.WriteString("item")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	_ = sb.String()

	// Good: var declaration without composite literal is allowed
	var buf bytes.Buffer
	buf.Grow(GrowSizeLarge)
	// Iteration over data to append
	for i := range LoopCountLarge {
		buf.WriteString("item")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	_ = buf.Bytes()

	// Good: pointer type is allowed (different use case)
	sb2 := &strings.Builder{}
	// Iteration over data to append
	for i := range LoopCountSmall {
		sb2.WriteString("x")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	_ = sb2.String()

	// Good: var declaration without composite literal
	var sb3 strings.Builder
	sb3.WriteString("single")
	_ = sb3.String()

	// Good: using new() instead of composite literal
	sb4 := new(strings.Builder)
	sb4.Grow(GrowSizeSmall)
	// Iteration over data to append
	for i := range GrowSizeSmall {
		sb4.WriteString("x")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	_ = sb4.String()
}
