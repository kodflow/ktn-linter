package test008 // want "KTN-TEST-008: le fichier 'number.go' n'a pas de fichier 'number_internal_test.go'"

// Double double un nombre
func Double(x int) int {
	return helper(x)
}

// helper est une fonction privée
func helper(x int) int {
	return x * 2
}
