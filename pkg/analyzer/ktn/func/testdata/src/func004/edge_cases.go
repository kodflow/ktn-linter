package func004

import "unsafe"

// Prevent "unsafe imported but not used" error
var _ = unsafe.Pointer(nil)

// Good: External function linked via go:linkname (no body to analyze)
// This tests the funcDecl.Body == nil branch in runFunc004
//
//go:linkname externalLinkedFunc runtime.convT64
func externalLinkedFunc(v int) (result unsafe.Pointer)

// Good: Another external function with named return
//
//go:linkname anotherExternal runtime.convTstring
func anotherExternal(v string) (ptr unsafe.Pointer)
