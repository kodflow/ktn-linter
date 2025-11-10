package test008 // want "KTN-TEST-008: le fichier 'number.go' contient des fonctions privées. Il doit avoir un fichier 'number_internal_test.go' \\(white-box\\)"

// Double double un nombre
func Double(x int) int {
	return helper(x)
}

// helper est une fonction privée
func helper(x int) int {
	return x * 2
}
