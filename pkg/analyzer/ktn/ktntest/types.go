// Types for the ktntest package.
package ktntest

// testFilesStatus contient l'état des fichiers de test.
type testFilesStatus struct {
	baseName    string
	fileBase    string
	hasInternal bool
	hasExternal bool
}
