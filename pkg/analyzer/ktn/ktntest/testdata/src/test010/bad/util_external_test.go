package test010_test

import "testing"

// Test_doWork teste doWork (fonction privée) dans external - ERREUR! // want "KTN-TEST-010: le test 'Test_doWork' dans 'util_external_test.go' teste une fonction privée 'doWork'"
func Test_doWork(t *testing.T) {
	// Ne peut pas compiler car doWork est privée, mais on veut tester la détection
	t.Log("This should not be here")
}
