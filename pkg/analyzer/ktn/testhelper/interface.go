// Test helper interfaces for the ktn-linter.
package testhelper

// TestingT est l'interface pour les méthodes de testing.T utilisées par testhelper.
type TestingT interface {
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
	Logf(format string, args ...any)
}
