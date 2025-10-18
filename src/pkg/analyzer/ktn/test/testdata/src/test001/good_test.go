package test001_test

import "testing"

// Bon : package avec suffixe _test
func TestSomething(t *testing.T) {
	t.Log("This test has the correct package name")
}
