// Package ktntest provides analyzers for test file lint rules.
package ktntest

// testFilesStatus contient l'état des fichiers de test.
type testFilesStatus struct {
	baseName    string
	fileBase    string
	hasInternal bool
	hasExternal bool
}
