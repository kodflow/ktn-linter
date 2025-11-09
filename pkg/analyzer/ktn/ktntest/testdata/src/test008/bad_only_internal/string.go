package test008 // want "KTN-TEST-008: le fichier 'string.go' n'a pas de fichier 'string_external_test.go'"

// Concat concatène deux strings
func Concat(a string, b string) string {
	return a + b
}
