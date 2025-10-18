package test001 // want `\[KTN_TEST_001\] Fichier de test 'bad_test.go' a le package 'test001' au lieu de 'test001_test'`

import "testing"

// Mauvais : package sans suffixe _test dans fichier _test.go
func TestSomething(t *testing.T) {
	t.Log("This test has the wrong package name")
}
