package test004_test

import "testing"

// Bon : fonctions de test dans un fichier _test.go
func TestProcess(t *testing.T) {
	t.Log("Testing Process")
}

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 1 + 2
	}
}
