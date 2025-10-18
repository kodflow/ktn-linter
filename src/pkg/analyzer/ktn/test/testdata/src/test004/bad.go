package test004

// Mauvais : fonction de test dans un fichier non-test
func TestSomething() { // want `\[KTN_TEST_004\] Fonction de test 'TestSomething' dans un fichier non-test 'bad.go'`
	// This test function should be in bad_test.go
}

func BenchmarkProcess() { // want `\[KTN_TEST_004\] Fonction de test 'BenchmarkProcess' dans un fichier non-test 'bad.go'`
	// This benchmark function should be in bad_test.go
}

// Bon : fonction normale
func DoSomething() string {
	return "something"
}
