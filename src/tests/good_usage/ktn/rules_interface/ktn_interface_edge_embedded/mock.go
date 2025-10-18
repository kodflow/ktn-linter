//go:build test
// +build test

package goodembedded

// MockReader est le mock de Reader.
type MockReader struct {
	ReadFunc func(p []byte) (n int, err error)
}

// Read implémente l'interface Reader.
//
// Params:
//   - p: buffer de destination
//
// Returns:
//   - n: nombre de bytes lus
//   - err: erreur éventuelle
func (m *MockReader) Read(p []byte) (n int, err error) {
	if m.ReadFunc != nil {
		// Early return from function.
		return m.ReadFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// MockWriter est le mock de Writer.
type MockWriter struct {
	WriteFunc func(p []byte) (n int, err error)
}

// Write implémente l'interface Writer.
//
// Params:
//   - p: buffer source
//
// Returns:
//   - n: nombre de bytes écrits
//   - err: erreur éventuelle
func (m *MockWriter) Write(p []byte) (n int, err error) {
	if m.WriteFunc != nil {
		// Early return from function.
		return m.WriteFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// MockReadWriter est le mock de ReadWriter.
type MockReadWriter struct {
	ReadFunc  func(p []byte) (n int, err error)
	WriteFunc func(p []byte) (n int, err error)
}

// Read implémente l'interface ReadWriter.
//
// Params:
//   - p: buffer de destination
//
// Returns:
//   - n: nombre de bytes lus
//   - err: erreur éventuelle
func (m *MockReadWriter) Read(p []byte) (n int, err error) {
	if m.ReadFunc != nil {
		// Early return from function.
		return m.ReadFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// Write implémente l'interface ReadWriter.
//
// Params:
//   - p: buffer source
//
// Returns:
//   - n: nombre de bytes écrits
//   - err: erreur éventuelle
func (m *MockReadWriter) Write(p []byte) (n int, err error) {
	if m.WriteFunc != nil {
		// Early return from function.
		return m.WriteFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// MockComplexEmbedded est le mock de ComplexEmbedded.
type MockComplexEmbedded struct {
	ReadFunc         func(p []byte) (n int, err error)
	WriteFunc        func(p []byte) (n int, err error)
	CloseFunc        func() error
	CustomMethodFunc func() error
}

// Read implémente l'interface ComplexEmbedded.
//
// Params:
//   - p: buffer de destination
//
// Returns:
//   - n: nombre de bytes lus
//   - err: erreur éventuelle
func (m *MockComplexEmbedded) Read(p []byte) (n int, err error) {
	if m.ReadFunc != nil {
		// Early return from function.
		return m.ReadFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// Write implémente l'interface ComplexEmbedded.
//
// Params:
//   - p: buffer source
//
// Returns:
//   - n: nombre de bytes écrits
//   - err: erreur éventuelle
func (m *MockComplexEmbedded) Write(p []byte) (n int, err error) {
	if m.WriteFunc != nil {
		// Early return from function.
		return m.WriteFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// Close implémente l'interface ComplexEmbedded.
//
// Returns:
//   - error: erreur si l'opération échoue
func (m *MockComplexEmbedded) Close() error {
	if m.CloseFunc != nil {
		// Early return from function.
		return m.CloseFunc()
	}
	// Early return from function.
	return nil
}

// CustomMethod implémente l'interface ComplexEmbedded.
//
// Returns:
//   - error: erreur si l'opération échoue
func (m *MockComplexEmbedded) CustomMethod() error {
	if m.CustomMethodFunc != nil {
		// Early return from function.
		return m.CustomMethodFunc()
	}
	// Early return from function.
	return nil
}

// MockProcessor est le mock de Processor.
type MockProcessor struct {
	ReadFunc    func(p []byte) (n int, err error)
	ProcessFunc func()
}

// Read implémente l'interface Processor.
//
// Params:
//   - p: buffer de destination
//
// Returns:
//   - n: nombre de bytes lus
//   - err: erreur éventuelle
func (m *MockProcessor) Read(p []byte) (n int, err error) {
	if m.ReadFunc != nil {
		// Early return from function.
		return m.ReadFunc(p)
	}
	// Early return from function.
	return 0, nil
}

// Process implémente l'interface Processor.
func (m *MockProcessor) Process() {
	if m.ProcessFunc != nil {
		m.ProcessFunc()
	}
}
