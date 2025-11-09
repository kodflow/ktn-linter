package test010_test

import "testing"

// Testhelper teste helper (fonction privée) dans external - ERREUR! // want "KTN-TEST-010: le test 'Testhelper' dans 'util_external_test.go' teste une fonction privée 'helper'"
func Testhelper(t *testing.T) {
	// Ne peut pas compiler car helper est privée, mais on veut tester la détection
	t.Log("This should not be here")
}
