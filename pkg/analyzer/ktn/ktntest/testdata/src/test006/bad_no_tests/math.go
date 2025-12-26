package test008 // want "KTN-TEST-006: le fichier 'math.go' contient des fonctions publiques. Il doit avoir un fichier 'math_external_test.go' \\(black-box\\)"

// Multiply multiplie deux entiers
func Multiply(a int, b int) int {
	return a * b
}
