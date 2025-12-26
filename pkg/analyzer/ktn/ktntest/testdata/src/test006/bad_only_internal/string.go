package test008 // want "KTN-TEST-006: le fichier 'string.go' contient des fonctions publiques. Il doit avoir un fichier 'string_external_test.go' \\(black-box\\)"

// Concat concatène deux strings
func Concat(a string, b string) string {
	return a + b
}
